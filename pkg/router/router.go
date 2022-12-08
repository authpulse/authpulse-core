package router

import (
	"context"
	"net"
	"net/http"

	"github.com/rs/zerolog/log"
)

func Run(ctx context.Context, listener net.Listener, handler http.Handler) error {
	srv := &http.Server{Handler: handler}

	go func() {
		<-ctx.Done()
		_ = srv.Shutdown(context.Background())
	}()

	log.Debug().Int("port", listener.Addr().(*net.TCPAddr).Port).Msg("starting server")
	if err := srv.Serve(listener); err != http.ErrServerClosed {
		return err
	}

	return nil
}
