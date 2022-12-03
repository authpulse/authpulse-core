package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/arthureichelberger/authpulse/pkg/router"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sign := make(chan os.Signal, 1)
	signal.Notify(sign, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	go func() {
		<-sign
		log.Debug().Msg("Shutting down")
		cancel()
	}()

	r := gin.New()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{})
	})

	if err := router.Run(ctx, ":8080", r); err != nil {
		log.Fatal().Err(err).Msg("Server crashed")
	}
}
