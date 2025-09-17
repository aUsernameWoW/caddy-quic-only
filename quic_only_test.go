package quiconly

import (
	"context"
	"net"
	"testing"
	
	"github.com/caddyserver/caddy/v2"
	"go.uber.org/zap"
)

// TestQuicOnly tests the QuicOnly module.
func TestQuicOnly(t *testing.T) {
	// Create a new QuicOnly instance
	qo := &QuicOnly{
		Protocols: []string{"h3"},
	}
	
	// Create a mock listener
	ln := &mockListener{}
	
	// Wrap the listener
	wrappedLn := qo.WrapListener(ln)
	
	// Check that the wrapped listener is the same as the original
	if wrappedLn != ln {
		t.Error("Expected wrapped listener to be the same as original listener")
	}
	
	// Test Provision method
	ctx, cancel := caddy.NewContext(caddy.Context{Context: context.Background()})
	defer cancel()
	
	// Set a logger in the context
	ctx.Context = context.WithValue(ctx.Context, caddy.LoggerCtxKey, zap.NewNop())
	
	err := qo.Provision(ctx)
	if err != nil {
		t.Errorf("Provision failed: %v", err)
	}
}

// mockListener is a mock implementation of net.Listener for testing.
type mockListener struct{}

func (ml *mockListener) Accept() (net.Conn, error) {
	return nil, nil
}

func (ml *mockListener) Close() error {
	return nil
}

func (ml *mockListener) Addr() net.Addr {
	return &net.TCPAddr{}
}