package geario

import (
	"net"
	"context"
	"time"
)

// Dialer is the interface that wraps the basic DialContext method.
type Dialer interface {
	DialContext(ctx context.Context, network, address string) (net.Conn, error)
}

// Dialer is the interface that wraps the basic DialContext method.
func (g *Gear) Dialer(d Dialer) Dialer {
	return DialGear(d, g, g)
}

// GearDialer returns a new Dialer that wraps the original Dialer with a Gear.
func GearDialer(d Dialer, duration time.Duration, limit B) Dialer {
	g := NewGear(duration, limit)
	return g.Dialer(d)
}

// DialGear returns a new Dialer that wraps the original Dialer with a Gear.
func DialGear(d Dialer, r, w *Gear) Dialer {
	return &gearDialer{
		dialer: d,
		r:      r,
		w:      w,
	}
}

type gearDialer struct {
	r, w   *Gear
	dialer Dialer
}

func (g *gearDialer) DialContext(ctx context.Context, network, address string) (net.Conn, error) {
	conn, err := g.dialer.DialContext(ctx, network, address)
	if err != nil {
		return nil, err
	}
	return ConnGear(conn, g.r, g.w), nil
}
