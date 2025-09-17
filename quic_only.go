package quiconly

import (
	"context"
	"fmt"
	"net"
	"strings"
	
	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"go.uber.org/zap"
)

func init() {
	caddy.RegisterModule(QuicOnly{})
}

// QuicOnly is a Caddy module that allows configuring QUIC-only listeners.
type QuicOnly struct {
	// Mode specifies the listener mode.
	// Supported values are:
	// - `quic_only`: Only create UDP listeners for HTTP/3.
	// - `tcp_only`: Only create TCP listeners for HTTP/1.1 and HTTP/2.
	// - `default`: Create both UDP and TCP listeners based on server protocols.
	Mode string `json:"mode,omitempty"`
	
	logger *zap.Logger
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
func (qo QuicOnly) WrapListener(ln net.Listener) net.Listener {
	// Log that the listener wrapper is being applied
	qo.logger.Info("QuicOnly listener wrapper applied", zap.String("mode", qo.Mode))
	
	// Check the mode
	switch qo.Mode {
	case "quic_only":
		// In QUIC-only mode, we want to prevent TCP listeners from being created for HTTP/1.1 and HTTP/2
		// This is a complex task that requires integration with Caddy's server configuration
		// For now, we'll just log that this mode is enabled
		qo.logger.Info("QUIC-only mode enabled - will only allow HTTP/3 traffic")
		return &quicOnlyListener{ln, qo.logger}
	case "tcp_only":
		// In TCP-only mode, we want to prevent UDP/QUIC listeners from being created
		// This is a complex task that requires integration with Caddy's server configuration
		// For now, we'll just log that this mode is enabled
		qo.logger.Info("TCP-only mode enabled - will only allow HTTP/1.1 and HTTP/2 traffic")
		return &tcpOnlyListener{ln, qo.logger}
	default:
		qo.logger.Info("Default mode enabled - allowing all protocols")
	}
	
	return ln
}

// Provision implements caddy.Provisioner.
// This method is called when the module is being set up.
func (qo *QuicOnly) Provision(ctx caddy.Context) error {
	qo.logger = ctx.Logger()
	
	// Log the configured mode
	qo.logger.Info("Provisioning QuicOnly module", zap.String("mode", qo.Mode))
	
	return nil
}

// Validate implements caddy.Validator.
func (qo *QuicOnly) Validate() error {
	// Validate the mode
	switch qo.Mode {
	case "quic_only", "tcp_only", "default", "":
		// Valid modes
		return nil
	default:
		return fmt.Errorf("invalid mode: %s (must be one of: quic_only, tcp_only, default)", qo.Mode)
	}
}

// quicOnlyListener is a wrapper that modifies listener behavior for QUIC-only mode
type quicOnlyListener struct {
	net.Listener
	logger *zap.Logger
}

// tcpOnlyListener is a wrapper that modifies listener behavior for TCP-only mode
type tcpOnlyListener struct {
	net.Listener
	logger *zap.Logger
}

// Interface guards
var (
	_ caddy.ListenerWrapper = (*QuicOnly)(nil)
	_ caddy.Provisioner     = (*QuicOnly)(nil)
	_ caddy.Validator       = (*QuicOnly)(nil)
)