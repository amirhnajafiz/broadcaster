package http

import (
	"log"
	"math/rand"
	"strconv"

	"github.com/sweetie-pie/line-recommendation/internal/model"
	"github.com/sweetie-pie/line-recommendation/internal/port/mysql"

	"github.com/gofiber/fiber/v2"
)

const (
	XUp   = 1000
	XDown = -1000
	YUp   = 1000
	yDown = -1000
)

type Handler struct {
	Repository *mysql.MySQL
}

func (h *Handler) CreateRoute(c *fiber.Ctx) error {
	tmp := c.Query("count", "4")
	count, _ := strconv.Atoi(tmp)

	nodes, err := h.Repository.GetNodes()
	if err != nil {
		log.Println(err)

		return fiber.ErrInternalServerError
	}

	for i := 0; i < count; i++ {
		// index of src and dest nodes
		srcIndex := rand.Intn(len(nodes))
		destIndex := srcIndex

		// make sure they are not the same
		for {
			destIndex = rand.Intn(len(nodes))
			if destIndex != srcIndex {
				break
			}
		}

		route := model.Route{
			StartID: nodes[srcIndex],
			StopID:  nodes[destIndex],
		}

		if err := h.Repository.InsertRoute(&route); err != nil {
			log.Println(err)

			return fiber.ErrInternalServerError
		}
	}

	return c.Status(fiber.StatusOK).SendString("Created!")
}

func (h *Handler) CreateNode(c *fiber.Ctx) error {
	return c.Next()
}

func (h *Handler) Search(c *fiber.Ctx) error {
	return c.Next()
}

func (h *Handler) Data(c *fiber.Ctx) error {
	routes, err := h.Repository.GetRoutes()
	if err != nil {
		log.Panicln(err)

		return fiber.ErrInternalServerError
	}

	routesResponse := make([]RouteResponse, 0)

	for _, route := range routes {
		src, _ := h.Repository.GetNode(route.StartID)
		dest, _ := h.Repository.GetNode(route.StopID)

		routesResponse = append(routesResponse, RouteResponse{
			ID:    route.ID,
			Start: *src,
			Stop:  *dest,
		})
	}

	searches, err := h.Repository.GetSearches()
	if err != nil {
		log.Panicln(err)

		return fiber.ErrInternalServerError
	}

	searchResponse := make([]RouteResponse, 0)

	for _, route := range searches {
		src, _ := h.Repository.GetNode(route.StartID)
		dest, _ := h.Repository.GetNode(route.StopID)

		searchResponse = append(searchResponse, RouteResponse{
			ID:    route.ID,
			Start: *src,
			Stop:  *dest,
		})
	}

	return c.JSON(&Response{
		Routes:   routesResponse,
		Searches: searchResponse,
	})
}
