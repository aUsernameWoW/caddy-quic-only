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
	
	// Initialize the logger
	logger := zap.NewNop()
	qo.logger = logger
	
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
	ctx.Context = context.WithValue(ctx.Context, caddy.CtxKey("logger"), logger)
	
	err := qo.Provision(ctx)
	if err != nil {
		t.Errorf("Provision failed: %v", err)
	}
	
	// Test tcp_only mode
	qo2 := &QuicOnly{
		Mode: "tcp_only",
	}
	qo2.logger = logger
	
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
	qo3.logger = logger
	
	wrappedLn3 := qo3.WrapListener(ln)
	
	// Check that the wrapped listener is the same as the original
	if wrappedLn3 != ln {
		t.Error("Expected wrapped listener to be the same as original listener")
	}
	
	err = qo3.Provision(ctx)
	if err != nil {
		t.Errorf("Provision failed: %v", err)
	}
	
	// Test empty mode (should default to default mode)
	qo4 := &QuicOnly{
		Mode: "",
	}
	qo4.logger = logger
	
	wrappedLn4 := qo4.WrapListener(ln)
	
	// Check that the wrapped listener is the same as the original
	if wrappedLn4 != ln {
		t.Error("Expected wrapped listener to be the same as original listener")
	}
	
	err = qo4.Provision(ctx)
	if err != nil {
		t.Errorf("Provision failed: %v", err)
	}
}

// TestValidate tests the Validate method.
func TestValidate(t *testing.T) {
	// Test valid modes
	validModes := []string{"quic_only", "tcp_only", "default", ""}
	for _, mode := range validModes {
		qo := &QuicOnly{Mode: mode}
		if err := qo.Validate(); err != nil {
			t.Errorf("Expected mode %s to be valid, but got error: %v", mode, err)
		}
	}
	
	// Test invalid mode
	qo := &QuicOnly{Mode: "invalid"}
	if err := qo.Validate(); err == nil {
		t.Error("Expected invalid mode to return an error, but got nil")
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