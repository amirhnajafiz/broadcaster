# Broadcaster

### What is this project?
Building a WebRTC video and audio Broadcaster in Golang using ION-SFU, and media devices.

### What is WebRTC ?
WebRTC or **Web Real-Time Communication** is a free open-source project providing 
web browsers and mobile applications with **real-time communication** (RTC) with APIs.
It allows audio and video communication to work inside web pages by allowing direct
**peer-to-peer** communication, eliminating the need to install plugins or download native
apps.

The technology behind WebRTC are implemented as an open web standard and available
as regular Javascript APIs in all major browsers.

As said before, major components of WebRTC includes several JS APIs:
- **getUserMedia** acquires the audio and video media
- **RTCPeerConnection** enables audio and video communication between peers. It performs signal processing, codec handling, peer-to-peer communication and ...
- **RTCDataChannel** allows bidirectional communication of arbitrary data between peers.

<p align="center">
    <img src="https://www.researchgate.net/profile/Martin-Meszaros-3/publication/328334940/figure/fig18/AS:682651635707904@1539768241002/WebRTC-triangle-with-SDES-and-DTLS-key-exchange-As-discussed-in-section-35-DTLS-SRTP.png" width="400" />
</p>

### What is SFU ?
SFU stands for **Selective Forwarding Unit**. Also known in the specifications as SFM (Selective Forwarding Middlebox).
At times, the term is used to describe a type of video routing device, while at other times
it will be used to indicate the support of routing technology and not a specific device.

An SFU is capable of receiving multiple media streams and then decide which of these media
streams should be sent to which participants.

SFU is a video routing service which allows webrtc sessions to scale more efficiently.

<p align="center">
    <img src="https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcSTiIJSpXcglL4a1_Z8pKNVx7kH7B8nh72xkXPC2xv-5tR6bmNkq2lMjM4PR2URk2HADws&usqp=CAU" width="392" />
</p>

### Why using ION-SFU ?
This package provides a simple, flexible, high performance Go implementation of a WebRTC SFU. 
It can be called directly or through a gPRC or json-rpc interface.

ION-SFU supports real-time processing on media streamed through the sfu using ion-avp.

Ion-avp is an extensible audio/video processing service designed for use with ios-sfu.

### What is media devices ?
The MediaDevices interface provides access to connected media input devices like cameras
and microphones, as well as screen sharing. In essence, it lets you obtain access to any 
hardware source of media data.

### What is ICE ?
ICE stands for **Interactive Connectivity Establishment**. It is a framework used by
WebRTC for connecting two peers, regardless of network topology.

This protocol allows two peers find and establish a connection with one another even
though they may both be using Network Address Translator to share a global IP address with
other devices on their respective local networks.

The framework algorithm looks for the lowest-latency path for connecting the two peers.

### What is SDP ?
The **Session Description Protocol** is a format for describing multimedia communication sessions
for the purposes of announcement and invitation. Its predominant use is in support of
streaming media applications, such as voice over IP and video conferencing. 

SDP does not deliver any media streams itself but is used between endpoints for negotiation of network metrics, media types
and other associated properties. The set of properties and parameters is called a session profile.

### How does this project work ?
WebRTC serves multiple purposes; together with the Media Capture and Streams API, they provide
powerful multimedia capabilities to the web, including support for audio and video conferencing,
file exchange, screen sharing, identity management. Connection between peers can be made
without requiring any special drivers or plugins.

Connections between two peers are represented by the **RTCPeerConnection** interface.
Once a connection has been established and opened using RTCPeerConnection, media streams (**MediaStream**)
and data channels (**RTCDataChannel**) can be added to the connection.

Media streams can consist of any number of tracks of media information; **tracks**, which are
represented by objects based on the **MediaStreamTrack** interface, may contain one of a number
of types of media data, including audio, video, and text.

Most streams consist of at least one audio track and likely also a video track, and
can be used to send and receive both live media or stored media information (such as streamed movie).

First we create a WebRTC connection between our local computer and a remote peer.
then we use bidirectional channels to transfer data between peers.

Once a user runs the application, we use Media devices to get the input data from our user
and send it to our ION-SFU server.

ION-SFU server on the other hand, gets the data and returns the response to all other peers, including our local machine.

### How to use this project ?
First you need the following requirements:
- go 17.1
- uuid 1.3
- websocket 1.4.2
- pion/mediadevices 0.3.2
- webritc 3.1.22
- jsonrpc2 1.0

Now you need to install the following packages in your system:
- npm
- pkg-config
- libx264 or x264
- libvpx or libvpx-dev
- vaapi 
- libopus or libopus-dev 

If you are using **Manjaro**, just run the following command:
```shell
sudo pacman -S a52dec faac faad2 flac jasper lame libdca libdv libmad libmpeg2 libtheora libvorbis libxv opus wavpack x264 xvidcore
```

Now you have to do the following steps to run the project.

Clone the project:
```shell
git clone https://github.com/amirhnajafiz/Broadcaster.git
```

Now enter the root directory:
```shell
cd Broadcaster
```

Clone the ION-SFU and execute it:
```shell
git clone https://github.com/pion/ion-sfu.git
```

Enter the root directory of INO-SFU:
```shell
ce ion-sfu
```

Now build and run the ION server:
```shell
go build ./cmd/signal/json-rpc/main.go && ./main -c config.toml
```

Once ION-SFU server is running you should see something like this:
```shell
[2022-02-08 02:30:40.584] [INFO] [main.go:94] => Config file loaded file=config.toml v=0
[2022-02-08 02:30:40.588] [INFO] [main.go:151] => --- Starting SFU Node --- v=0
[2022-02-08 02:30:40.588] [INFO] [main.go:188] => Started listening addr=http://:7000 v=0
[2022-02-08 02:30:40.589] [INFO] [main.go:130] => Metrics Listening addr=:8100 v=0
```

Now we need to run our WebRTC server:
```shell
cd ..
make serve
```

If you got the following response you are good to go:
```shell
2022/02/08 02:35:02 connecting to ws://localhost:7000/ws
[{79584686-6ab6-4e72-b4e2-bcf3cff5bb24 1 pci-0000:00:14.0-usb-0:8:1.0-video-index0;video0 camera} {7e917d25-5db4-451d-82c2-819b2b28e3cf 1 pci-0000:00:14.0-usb-0:8:1.0-video-index1;video1 camera} {082aeb80-85ee-438e-9212-cf9778555fc1 2 616c73615f6f75747075742e7063692d303030305f30305f31662e332e616e616c6f672d73746572656f2e6d6f6e69746f72 microphone} {65ec460f-5188-4895-9e87-69f908bfbfa1 2 616c73615f696e7075742e7063692d303030305f30305f31662e332e616e616c6f672d73746572656f microphone}]
recv: {"id":2473120800,"result":{"type":"answer","sdp":"v=0\r\no=- 7188600485215558167 1644275102 IN IP4 0.0.0.0\r\ns=-\r\nt=0 0\r\na=fingerprint:sha-256 CC:7F:6A:CD:1E:DD:E3:22:4C:BA:5E:71:70:09:62:1C:2F:84:E6:69:FF:44:41:48:98:6A:AC:7B:A6:5C:15:03\r\na=group:BUNDLE 0\r\nm=video 9 UDP/TLS/RTP/SAVPF 96\r\nc=IN IP4 0.0.0.0\r\na=setup:active\r\na=mid:0\r\na=ice-ufrag:tLebVdfImSyBFMxE\r\na=ice-pwd:MTeiJhQRVWWCMVRhmfNmjhpwjbekFeFN\r\na=rtcp-mux\r\na=rtcp-rsize\r\na=rtpmap:96 VP8/90000\r\na=recvonly\r\n"},"jsonrpc":"2.0"}
recv: {"method":"trickle","params":{"target":0,"candidate":{"candidate":"candidate:3388917061 1 udp 2130706431 192.168.1.104 5086 typ host","sdpMid":"","sdpMLineIndex":0,"usernameFragment":null}},"jsonrpc":"2.0"}
Connection State has changed to checking 
recv: {"method":"offer","params":{"type":"offer","sdp":"v=0\r\no=- 4704832186149651630 1644275102 IN IP4 0.0.0.0\r\ns=-\r\nt=0 0\r\na=fingerprint:sha-256 E6:7E:E2:5B:BA:CC:29:52:33:6A:F1:E8:CD:15:18:8C:D3:55:6F:3E:10:21:8C:5C:D5:D7:39:03:A1:2C:82:E1\r\na=group:BUNDLE 0\r\nm=application 9 UDP/DTLS/SCTP webrtc-datachannel\r\nc=IN IP4 0.0.0.0\r\na=setup:actpass\r\na=mid:0\r\na=sendrecv\r\na=sctp-port:5000\r\na=ice-ufrag:nbjXQnCkVUSRyCmR\r\na=ice-pwd:rHqvfuEMBvmXgtnMlYLISGcWFJQzfDyv\r\n"},"jsonrpc":"2.0"}
recv: {"method":"trickle","params":{"target":1,"candidate":{"candidate":"candidate:3388917061 1 udp 2130706431 192.168.1.104 5093 typ host","sdpMid":"","sdpMLineIndex":0,"usernameFragment":null}},"jsonrpc":"2.0"}
Connection State has changed to connected 
```

Now we need to run the client:
```shell
cd client
node server.js
```

You should get the following result:
```shell
Server is running on port 3030
``` 

Now you can check the application on **localhost:3030** to see something like this:
<p align="center">
    <img src="assets/demo.png" />
</p>