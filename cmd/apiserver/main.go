package main

import (
	"github.com/VSKrivoshein/test/internal/app/broker"
	"github.com/VSKrivoshein/test/internal/app/cache"
	"github.com/VSKrivoshein/test/internal/app/configs"
	"github.com/VSKrivoshein/test/internal/app/repository"
	"github.com/VSKrivoshein/test/internal/app/service"
	"github.com/VSKrivoshein/test/internal/app/utils"
	log "github.com/sirupsen/logrus"
)

func main() {
	utils.UseJSONLogFormat()

	c := cache.New()
	r := repository.New(c)
	b := broker.New()
	s := service.New(r, b)

	errors := make(chan error)

	go func() {
		errors <- configs.StartGrpc(s)
	}()

	go func() {
		errors <- configs.GracefulShutdown()
	}()

	log.Fatalf("Terminated %s", <-errors)
}
