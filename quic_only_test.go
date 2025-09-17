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
	// Create a new QuicOnly instance with quic_only mode
	qo := &QuicOnly{
		Mode: "quic_only",
	}
	
	// Create a mock listener
	ln := &mockListener{}
	
	// Wrap the listener
	wrappedLn := qo.WrapListener(ln)
	
	// Check that the wrapped listener is not the same as the original (should be wrapped)
	if wrappedLn == ln {
		t.Error("Expected wrapped listener to be different from original listener")
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
	
	// Test tcp_only mode
	qo2 := &QuicOnly{
		Mode: "tcp_only",
	}
	
	wrappedLn2 := qo2.WrapListener(ln)
	
	// Check that the wrapped listener is not the same as the original (should be wrapped)
	if wrappedLn2 == ln {
		t.Error("Expected wrapped listener to be different from original listener")
	}
	
	err = qo2.Provision(ctx)
	if err != nil {
		t.Errorf("Provision failed: %v", err)
	}
	
	// Test default mode
	qo3 := &QuicOnly{
		Mode: "default",
	}
	
	wrappedLn3 := qo3.WrapListener(ln)
	
	// Check that the wrapped listener is the same as the original
	if wrappedLn3 != ln {
		t.Error("Expected wrapped listener to be the same as original listener")
	}
	
	err = qo3.Provision(ctx)
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