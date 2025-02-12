package server

import (
	"github.com/gofiber/fiber/v2"
)

func (hs *httpServer) setupRoutes() {
    v1 := hs.Group("v1")

	hs.healthCheckRoutes()

	hs.videoMakerRoutesV1(v1)
}

func (hs *httpServer) healthCheckRoutes() {
	hs.Get("/health", func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusOK) })
}

func (hs *httpServer) videoMakerRoutesV1(version fiber.Router) {
	group := version.Group("video")

	group.Post("/generate", hs.restHandler.MakeVideo)
}
