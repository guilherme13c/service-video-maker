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
	log.Println("Starting service...")

	log.Println("Loading config...")
	cfg, errCfg := config.New()
	if errCfg != nil {
		log.Fatalf("Error loading config: %v", errCfg)
	}

	log.Println("Creating rest client...")
	restClient := resty.New()

	log.Println("Creating repository...")
	repository := repository.NewRepository()

	log.Println("Creating service...")
	service := service.NewService(restClient, repository)

	log.Println("Creating rest handler...")
	restHandler := restHandler.NewRestHandler(service)

	log.Println("Creating server...")
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

	log.Println("Running server...")
	s.Run()

	serverShutdown.Wait()
	fmt.Println("Running cleanup tasks...")
	fmt.Println("End of cleanup tasks...")
}
