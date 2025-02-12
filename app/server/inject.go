package server

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	resthandler "github.com/simpleAI/service-video-maker/app/handler/rest"
	"github.com/simpleAI/service-video-maker/app/resource/config"
)

type IServer interface {
	Run()
	GetFiberApp() *fiber.App
}

func (hs *httpServer) configure() {
	hs.Use(recover.New())
	hs.Use(logger.New())
	hs.Use(cors.New())
}

func (hs *httpServer) GetFiberApp() *fiber.App {
	return hs.App
}

func (hs *httpServer) Run() {
	err := hs.Listen(":" + hs.cfg.Port)
	if err != nil {
		log.Fatal(err)
	}
}

type httpServer struct {
	*fiber.App
	cfg         *config.Config
	restHandler resthandler.IRestHandler
}

func NewHttpServer(cfg *config.Config, restHandler resthandler.IRestHandler) IServer {
	prefork := false
	if cfg.Environment == "PRODUCTION" {
		prefork = true
	}

	server := new(httpServer)

	server.App = fiber.New(fiber.Config{
		Prefork: prefork,
	})
	server.cfg = cfg
	server.restHandler = restHandler

	server.configure()
	server.setupRoutes()

	return server
}
