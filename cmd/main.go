package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/pion/mediadevices"
	"github.com/pion/mediadevices/pkg/codec/vpx"
	"github.com/pion/mediadevices/pkg/frame"
	"github.com/pion/mediadevices/pkg/prop"
	"github.com/pion/webrtc/v3"
	"github.com/sourcegraph/jsonrpc2"
	"io"
	"log"
	"net/url"
)

var (
	addr           string
	peerConnection *webrtc.PeerConnection
)

type SendOffer struct {
	SID   string                     `json:"sid"`
	Offer *webrtc.SessionDescription `json:"offer"`
}

type Candidate struct {
	Target    int                  `json:"target"`
	Candidate *webrtc.ICECandidate `json:"candidate"`
}

func main() {
	flag.StringVar(&addr, "a", "localhost:7000", "address to use")
	flag.Parse()

	u := url.URL{
		Scheme: "ws",
		Host:   addr,
		Path:   "/ws",
	}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer func(c *websocket.Conn) {
		err := c.Close()
		if err != nil {
			log.Fatal("close fatal:", err)
		}
	}(c)

	done := make(chan struct{})

	go readMessage(c, done)

	<-done

	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
		SDPSemantics: webrtc.SDPSemanticsUnifiedPlanWithFallback,
	}

	mediaEngine := webrtc.MediaEngine{}

	vpxParams, err := vpx.NewVP8Params()
	if err != nil {
		panic(err)
	}
	vpxParams.BitRate = 500_000 // 500kbps

	codecSelector := mediadevices.NewCodecSelector(
		mediadevices.WithVideoEncoders(&vpxParams),
	)

	codecSelector.Populate(&mediaEngine)
	api := webrtc.NewAPI(webrtc.WithMediaEngine(&mediaEngine))
	peerConnection, err = api.NewPeerConnection(config)
	if err != nil {
		panic(err)
	}

	fmt.Println(mediadevices.EnumerateDevices())

	s, err := mediadevices.GetUserMedia(mediadevices.MediaStreamConstraints{
		Video: func(constraints *mediadevices.MediaTrackConstraints) {
			constraints.FrameFormat = prop.FrameFormat(frame.FormatYUY2)
			constraints.Width = prop.Int(640)
			constraints.Height = prop.Int(480)
		},
		Codec: codecSelector,
	})

	if err != nil {
		panic(err)
	}

	for _, track := range s.GetTracks() {
		track.OnEnded(func(err error) {
			fmt.Printf("Track (ID: %s) ended with error: %v\n", track.ID(), err)
		})
		_, err = peerConnection.AddTransceiverFromTrack(track,
			webrtc.RTPTransceiverInit{
				Direction: webrtc.RTPTransceiverDirectionSendonly,
			},
		)

		if err != nil {
			panic(err)
		}
	}

	// WebRTC offer
	offer, err := peerConnection.CreateOffer(nil)

	// Remote Session description
	err = peerConnection.SetLocalDescription(offer)
	if err != nil {
		panic(err)
	}

	offerJson, err := json.Marshal(&SendOffer{
		Offer: peerConnection.LocalDescription(),
		SID:   "test room",
	})

	params := (*json.RawMessage)(&offerJson)

	connectionUUID := uuid.New()
	connectionID := uint64(connectionUUID.ID())

	offerMessage := &jsonrpc2.Request{
		Method: "join",
		Params: params,
		ID: jsonrpc2.ID{
			IsString: false,
			Str:      "",
			Num:      connectionID,
		},
	}

	reqBodyBytes := new(bytes.Buffer)
	_ = json.NewEncoder(reqBodyBytes).Encode(offerMessage)

	messageBytes := reqBodyBytes.Bytes()
	_ = c.WriteMessage(websocket.TextMessage, messageBytes)

	// Handling OnICECandidate event
	peerConnection.OnICECandidate(func(candidate *webrtc.ICECandidate) {
		if candidate != nil {
			candidateJSON, err := json.Marshal(&Candidate{
				Candidate: candidate,
				Target:    0,
			})

			params := (*json.RawMessage)(&candidateJSON)

			if err != nil {
				log.Fatal(err)
			}

			message := &jsonrpc2.Request{
				Method: "trickle",
				Params: params,
			}

			reqBodyBytes := new(bytes.Buffer)
			_ = json.NewEncoder(reqBodyBytes).Encode(message)

			messageBytes := reqBodyBytes.Bytes()
			_ = c.WriteMessage(websocket.TextMessage, messageBytes)
		}
	})

	peerConnection.OnICEConnectionStateChange(func(state webrtc.ICEConnectionState) {
		fmt.Printf("Connection State has changed to %s \n", state.String())
	})
}

func readMessage(connection *websocket.Conn, done chan struct{}) {
	defer close(done)

	for {
		_, message, err := connection.ReadMessage()
		if err != nil || err == io.EOF {
			log.Fatal("Error reading: ", err)
			return
		}

		fmt.Printf("recv: %s", message)
	}
}
