package router_test

import (
	"context"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/arthureichelberger/authpulse/pkg/router"
	"github.com/stretchr/testify/require"
)

func TestRun(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })

	listener, _ := net.Listen("tcp", ":0")

	go func() {
		err := router.Run(ctx, listener, handler)
		require.NoError(t, err)
	}()

	time.Sleep(time.Second)
	cancel()
}
