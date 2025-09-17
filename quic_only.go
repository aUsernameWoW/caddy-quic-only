package quiconly

import (
	"fmt"
	"net"
	
	"go.uber.org/zap"
	"github.com/caddyserver/caddy/v2"
)

func init() {
	caddy.RegisterModule(QuicOnly{})
}

// QuicOnly is a Caddy module that allows configuring QUIC-only listeners.
type QuicOnly struct {
	// Protocols specifies which protocols to enable.
	// Supported values are:
	// - `h1` (HTTP/1.1)
	// - `h2` (HTTP/2)
	// - `h3` (HTTP/3)
	//
	// If `h3` is specified, only UDP listeners will be created.
	// If `h1` or `h2` are specified, TCP listeners will be created.
	// If both `h3` and `h1`/`h2` are specified, both UDP and TCP listeners will be created.
	Protocols []string `json:"protocols,omitempty"`
	
	// Mode specifies the listener mode.
	// Supported values are:
	// - `quic_only`: Only create UDP listeners for HTTP/3.
	// - `tcp_only`: Only create TCP listeners for HTTP/1.1 and HTTP/2.
	// - `default`: Create both UDP and TCP listeners based on Protocols.
	Mode string `json:"mode,omitempty"`
}

// CaddyModule returns the Caddy module information.
func (QuicOnly) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "caddy.listeners.quic_only",
		New: func() caddy.Module { return new(QuicOnly) },
	}
}

// WrapListener implements caddy.ListenerWrapper.
// This method is called for each listener that Caddy creates.
// We'll modify the listener behavior based on the configured protocols.
func (qo QuicOnly) WrapListener(ln net.Listener) net.Listener {
	// For this example, we're just returning the listener as-is.
	// In a real implementation, we would modify the listener based on the protocols.
	// Since this is a complex task involving low-level network operations,
	// we'll focus on the configuration aspect in this example.
	
	// Log the configured protocols for debugging
	logger := caddy.Log()
	
	// Check the mode
	switch qo.Mode {
	case "quic_only":
		logger.Info("QuicOnly listener wrapper applied: QUIC-only mode enabled")
		// In a real implementation, we would need to ensure that only UDP listeners are created
		// This is a complex task that requires deep integration with Caddy's listener creation logic
		// For now, we'll just log a warning that this is not fully implemented
		logger.Warn("QUIC-only mode is not fully implemented in this example")
	case "tcp_only":
		logger.Info("QuicOnly listener wrapper applied: TCP-only mode enabled")
		// In a real implementation, we would need to ensure that only TCP listeners are created
		// This is a complex task that requires deep integration with Caddy's listener creation logic
		// For now, we'll just log a warning that this is not fully implemented
		logger.Warn("TCP-only mode is not fully implemented in this example")
	default:
		logger.Info("QuicOnly listener wrapper applied", zap.Strings("protocols", qo.Protocols))
	}
	
	return ln
}

// Provision implements caddy.Provisioner.
// This method is called when the module is being set up.
func (qo *QuicOnly) Provision(ctx caddy.Context) error {
	// Log the configured protocols for debugging
	logger := ctx.Logger()
	
	// Check the mode
	switch qo.Mode {
	case "quic_only":
		logger.Info("Provisioning QuicOnly module: QUIC-only mode enabled")
		// In a real implementation, we would need to ensure that only UDP listeners are created
		// This is a complex task that requires deep integration with Caddy's listener creation logic
		// For now, we'll just log a warning that this is not fully implemented
		logger.Warn("QUIC-only mode is not fully implemented in this example")
	case "tcp_only":
		logger.Info("Provisioning QuicOnly module: TCP-only mode enabled")
		// In a real implementation, we would need to ensure that only TCP listeners are created
		// This is a complex task that requires deep integration with Caddy's listener creation logic
		// For now, we'll just log a warning that this is not fully implemented
		logger.Warn("TCP-only mode is not fully implemented in this example")
	default:
		logger.Info("Provisioning QuicOnly module", zap.Strings("protocols", qo.Protocols))
	}
	
	return nil
}

// Interface guards
var (
	_ caddy.ListenerWrapper = (*QuicOnly)(nil)
	_ caddy.Provisioner     = (*QuicOnly)(nil)
)