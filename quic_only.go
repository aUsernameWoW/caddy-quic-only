package quiconly

import (
	"errors"
	"net"
	
	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
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
// Note: The actual protocol filtering is done in the Provision method.
// This wrapper is kept for consistency with the ListenerWrapper interface.
func (qo QuicOnly) WrapListener(ln net.Listener) net.Listener {
	// Log that the listener wrapper is being applied
	qo.logger.Info("QuicOnly listener wrapper applied", zap.String("mode", qo.Mode))
	
	// In this implementation, the actual protocol filtering is done in the Provision method
	// by modifying the server's Protocols field. This wrapper simply returns the original
	// listener without modification, as the protocol filtering is handled at a higher level.
	return ln
}

// Provision implements caddy.Provisioner.
// This method is called when the module is being set up.
func (qo *QuicOnly) Provision(ctx caddy.Context) error {
	qo.logger = ctx.Logger()
	
	// Log the configured mode
	qo.logger.Info("Provisioning QuicOnly module", zap.String("mode", qo.Mode))
	
	// If we're in a server context, modify the server's protocols
	serverValue := ctx.Value(caddyhttp.ServerCtxKey)
	if serverValue == nil {
		qo.logger.Warn("Server context not available, skipping protocol modification")
		return nil
	}
	
	server, ok := serverValue.(*caddyhttp.Server)
	if !ok {
		qo.logger.Warn("Server context is not of expected type, skipping protocol modification")
		return nil
	}
	
	if server != nil {
		qo.logger.Debug("Modifying server protocols", zap.String("mode", qo.Mode))
		
		// Modify the server's protocols based on the mode
		switch qo.Mode {
		case "quic_only":
			// For QUIC-only mode, we want to only enable HTTP/3
			// Filter out h1 and h2 protocols, keep only h3
			protocols := []string{}
			for _, p := range server.Protocols {
				if p == "h3" {
					protocols = append(protocols, p)
				}
			}
			// If h3 is not in the original protocols, we don't add it automatically
			// This prevents enabling HTTP/3 when it wasn't intended
			if len(protocols) == 0 {
				qo.logger.Info("QUIC-only mode enabled but h3 not in original protocols, no protocols configured")
			} else {
				server.Protocols = protocols
				qo.logger.Info("Configured server for QUIC-only mode", zap.Strings("protocols", server.Protocols))
			}
		case "tcp_only":
			// For TCP-only mode, we want to only enable HTTP/1.1 and HTTP/2
			// Filter out h3 protocol, keep h1 and h2
			protocols := []string{}
			for _, p := range server.Protocols {
				if p == "h1" || p == "h2" {
					protocols = append(protocols, p)
				}
			}
			// If neither h1 nor h2 is in the original protocols, we don't add them automatically
			if len(protocols) == 0 {
				qo.logger.Info("TCP-only mode enabled but neither h1 nor h2 in original protocols, no protocols configured")
			} else {
				server.Protocols = protocols
				qo.logger.Info("Configured server for TCP-only mode", zap.Strings("protocols", server.Protocols))
			}
		default:
			qo.logger.Info("Default mode enabled - using server's configured protocols", zap.Strings("protocols", server.Protocols))
		}
	}
	
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
		return errors.New("invalid mode: " + qo.Mode + " (must be one of: quic_only, tcp_only, default)")
	}
}

// UnmarshalCaddyfile implements caddyfile.Unmarshaler.
func (qo *QuicOnly) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	for d.Next() {
		args := d.RemainingArgs()
		if len(args) > 0 {
			return d.ArgErr()
		}
		for nesting := d.Nesting(); d.NextBlock(nesting); {
			switch d.Val() {
			case "mode":
				if !d.NextArg() {
					return d.ArgErr()
				}
				qo.Mode = d.Val()
			default:
				return d.Errf("unrecognized subdirective '%s'", d.Val())
			}
		}
	}
	return nil
}



// Interface guards
var (
	_ caddy.ListenerWrapper     = (*QuicOnly)(nil)
	_ caddy.Provisioner         = (*QuicOnly)(nil)
	_ caddy.Validator           = (*QuicOnly)(nil)
	_ caddyfile.Unmarshaler     = (*QuicOnly)(nil)
)