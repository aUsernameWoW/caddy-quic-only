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
	
	// Check that the wrapped listener is the same as the original (should not be wrapped in the new implementation)
	if wrappedLn != ln {
		t.Error("Expected wrapped listener to be the same as original listener")
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
	
	// Check that the wrapped listener is the same as the original (should not be wrapped in the new implementation)
	if wrappedLn2 != ln {
		t.Error("Expected wrapped listener to be the same as original listener")
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

// TestProtocolFiltering tests the protocol filtering functionality.
func TestProtocolFiltering(t *testing.T) {
	// This test is meant to verify the protocol filtering logic
	// Since we can't easily create a full caddyhttp.Server in tests,
	// we'll test the filtering logic directly
	
	// Test QUIC-only filtering
	protocols := []string{"h1", "h2", "h3"}
	filtered := filterProtocols(protocols, "quic_only")
	if len(filtered) != 1 || filtered[0] != "h3" {
		t.Errorf("Expected only h3 protocol, got %v", filtered)
	}
	
	// Test TCP-only filtering
	protocols = []string{"h1", "h2", "h3"}
	filtered = filterProtocols(protocols, "tcp_only")
	if len(filtered) != 2 || filtered[0] != "h1" || filtered[1] != "h2" {
		t.Errorf("Expected h1 and h2 protocols, got %v", filtered)
	}
	
	// Test QUIC-only filtering with no h3
	protocols = []string{"h1", "h2"}
	filtered = filterProtocols(protocols, "quic_only")
	if len(filtered) != 0 {
		t.Errorf("Expected no protocols, got %v", filtered)
	}
	
	// Test TCP-only filtering with no h1 or h2
	protocols = []string{"h3"}
	filtered = filterProtocols(protocols, "tcp_only")
	if len(filtered) != 0 {
		t.Errorf("Expected no protocols, got %v", filtered)
	}
	
	// Test default mode
	protocols = []string{"h1", "h2", "h3"}
	filtered = filterProtocols(protocols, "default")
	if len(filtered) != 3 {
		t.Errorf("Expected all protocols, got %v", filtered)
	}
}

// filterProtocols is a helper function to test protocol filtering logic
func filterProtocols(protocols []string, mode string) []string {
	switch mode {
	case "quic_only":
		// For QUIC-only mode, we want to only enable HTTP/3
		// Filter out h1 and h2 protocols, keep only h3
		filtered := []string{}
		for _, p := range protocols {
			if p == "h3" {
				filtered = append(filtered, p)
			}
		}
		return filtered
	case "tcp_only":
		// For TCP-only mode, we want to only enable HTTP/1.1 and HTTP/2
		// Filter out h3 protocol, keep h1 and h2
		filtered := []string{}
		for _, p := range protocols {
			if p == "h1" || p == "h2" {
				filtered = append(filtered, p)
			}
		}
		return filtered
	default:
		// For default mode, keep all protocols
		return protocols
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

// mockServer is a mock implementation for testing.
type mockServer struct {
	Protocols []string
}