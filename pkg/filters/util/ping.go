package util

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/sparrc/go-ping"
)

const (
	DefaultPingIP = "8.8.8.8"
)

// Creates a new Pinger object by providing the IP address to ping to.
// Provides a hard coded values to the number of pings to make(1) and
// sets the timeout to one second.
func NewPinger(testIP string) (*ping.Pinger, error) {
	p, err := ping.NewPinger(testIP)
	if err != nil {
		return nil, err
	}

	p.Count = 1
	p.SetPrivileged(true)
	p.Timeout = (100 * time.Millisecond)

	return p, nil
}

func CheckPingSuccess(ip string) error {
	pingIP := DefaultPingIP

	if ip != "" {
		pingIP = ip
	}

	pinger, err := NewPinger(pingIP)
	if err != nil {
		return fmt.Errorf("creating pinger: %v", err)
	}

	pinger.Run()

	if pinger.PacketsRecv == 0 {
		return fmt.Errorf("imposible to ping %q", ip)
	}

	return nil
}

func CheckPingFail(ip string) error {
	pingIP := DefaultPingIP

	if ip != "" {
		pingIP = ip
	}

	pinger, err := NewPinger(pingIP)
	if err != nil {
		return fmt.Errorf("creating pinger: %v", err)
	}

	pinger.Run()

	if pinger.PacketsRecv >= 1 {
		return fmt.Errorf("ping to %q should have fail", ip)
	}

	return nil
}

func CheckDNSSuccess(ip string) error {
	r := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: time.Second * time.Duration(5),
			}

			return d.DialContext(ctx, "udp", ip+":53")
		},
	}

	_, err := r.LookupHost(context.Background(), "www.google.com")

	return err
}

func CheckDNSFail(ip string) error {
	r := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: time.Second * time.Duration(5),
			}

			return d.DialContext(ctx, "udp", ip+":53")
		},
	}

	_, err := r.LookupHost(context.Background(), "www.google.com")
	if err == nil {
		return fmt.Errorf("DNS request to %s:53 should have failed", ip)
	}

	return nil
}
