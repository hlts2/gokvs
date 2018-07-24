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

// EchoMessage -
func (i *ICMP) EchoMessage(ip string, seq int) (bool, error) {
	msg := icmp.Message{
		Type: ipv4.ICMPTypeEchoReply,
		Code: 0,
		Body: &icmp.Echo{
			ID:   os.Getpid() & 0xffff,
			Seq:  seq,
			Data: []byte("HELLO"),
		},
	}

	wd, err := msg.Marshal(nil)
	if err != nil {
		return false, err
	}

	_, err = i.conn.WriteTo(wd, &net.IPAddr{IP: net.ParseIP(ip)})
	if err != nil {
		return false, err
	}

	rb := make([]byte, 1500)
	n, _, err := i.conn.ReadFrom(rb)
	if err != nil {
		return false, err
	}

	h, err := icmp.ParseMessage(ipv4.ICMPTypeEcho.Protocol(), rb[:n])
	if err == nil && h.Type == ipv4.ICMPTypeEchoReply {
		return true, nil
	}

	return true, nil
}
