package router_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/arthureichelberger/authpulse/pkg/router"
	"github.com/phayes/freeport"
	"github.com/stretchr/testify/require"
)

func TestRun(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })

	port, _ := freeport.GetFreePort()
	go func() {
		err := router.Run(ctx, fmt.Sprintf(":%d", port), handler)
		require.NoError(t, err)
	}()

	time.Sleep(time.Second)
	cancel()
}
