package router_test

import (
	"context"
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

	go func() {
		err := router.Run(ctx, ":8080", handler)
		require.NoError(t, err)
	}()

	time.Sleep(time.Second)
	cancel()
}
