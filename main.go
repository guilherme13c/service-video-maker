package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
	restHandler "github.com/simpleAI/service-video-maker/app/handler/rest"
	"github.com/simpleAI/service-video-maker/app/repository"
	"github.com/simpleAI/service-video-maker/app/resource/config"
	"github.com/simpleAI/service-video-maker/app/server"
	"github.com/simpleAI/service-video-maker/app/service"
)

func main() {
	cfg, errCfg := config.New()
	if errCfg != nil {
		log.Fatalf("Error loading config: %v", errCfg)
	}

	restClient := resty.New()

	repository := repository.NewRepository()

	service := service.NewService(restClient, repository)

	restHandler := restHandler.NewRestHandler(service)

	s := server.NewHttpServer(cfg, restHandler)

	//Set shutdown server
	app := s.GetFiberApp()
	var serverShutdown sync.WaitGroup
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		_ = <-c
		fmt.Println("Gracefully shutting down...")
		serverShutdown.Add(1)
		defer serverShutdown.Done()
		_ = app.ShutdownWithTimeout(30 * time.Second)
	}()

	s.Run()

	serverShutdown.Wait()
	fmt.Println("Running cleanup tasks...")
	fmt.Println("End of cleanup tasks...")
}
