package icmp

import (
	"net"
	"os"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

// ICMP -
type ICMP struct {
	conn *icmp.PacketConn
}

// New -
func New() (*ICMP, error) {
	c, err := icmp.ListenPacket("ipv4:icmp", "0.0.0.0")
	if err != nil {
		return nil, err
	}

	return &ICMP{
		conn: c,
	}, nil
}

// Send -
func (i *ICMP) Send(ip string, seq int) (bool, error) {
	wbyte, err := (&icmp.Message{
		Type: ipv4.ICMPTypeEchoReply,
		Body: &icmp.Echo{
			ID:   os.Getpid() & 0xffff,
			Seq:  seq,
			Data: []byte("HELLO"),
		},
	}).Marshal(nil)
	if err != nil {
		return false, err
	}

	_, err = i.conn.WriteTo(wbyte, &net.IPAddr{IP: net.ParseIP(ip)})
	if err != nil {
		return false, err
	}

	rbyte := make([]byte, 1500)
	n, _, err := i.conn.ReadFrom(rbyte)
	if err != nil {
		return false, err
	}

	h, err := icmp.ParseMessage(ipv4.ICMPTypeEcho.Protocol(), rbyte[:n])
	if err == nil && h.Type == ipv4.ICMPTypeEchoReply {
		return true, nil
	}

	return true, nil
}
