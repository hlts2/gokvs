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
func (i *ICMP) Send(ip string, seq int, deadNode chan string) {
	wbyte, err := (&icmp.Message{
		Type: ipv4.ICMPTypeEchoReply,
		Body: &icmp.Echo{
			ID:   os.Getpid() & 0xffff,
			Seq:  seq,
			Data: []byte("HELLO"),
		},
	}).Marshal(nil)
	if err != nil {
		deadNode <- ip
		return
	}

	_, err = i.conn.WriteTo(wbyte, &net.IPAddr{IP: net.ParseIP(ip)})
	if err != nil {
		deadNode <- ip
		return
	}

	rbyte := make([]byte, 1500)
	n, _, err := i.conn.ReadFrom(rbyte)
	if err != nil {
		deadNode <- ip
		return
	}

	h, err := icmp.ParseMessage(ipv4.ICMPTypeEcho.Protocol(), rbyte[:n])
	if err != nil || h.Type != ipv4.ICMPTypeEchoReply {
		deadNode <- ip
		return
	}
}
