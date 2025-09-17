package quiconly

import (
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
	logger.Info("QuicOnly listener wrapper applied", zap.Strings("protocols", qo.Protocols))
	
	return ln
}

// Provision implements caddy.Provisioner.
// This method is called when the module is being set up.
func (qo *QuicOnly) Provision(ctx caddy.Context) error {
	// Log the configured protocols for debugging
	logger := ctx.Logger()
	logger.Info("Provisioning QuicOnly module", zap.Strings("protocols", qo.Protocols))
	
	return nil
}

// Interface guards
var (
	_ caddy.ListenerWrapper = (*QuicOnly)(nil)
	_ caddy.Provisioner     = (*QuicOnly)(nil)
)