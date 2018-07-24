package server

import (
	"net/http"
	"time"

	"github.com/hlts2/gokvs/config"
	"github.com/hlts2/gokvs/icmp"
	lockfree "github.com/hlts2/lock-free"
	"github.com/pkg/errors"
)

// ErrNotExists --
var ErrNotExists = errors.New("server dose not exist")

// Server -
type Server interface {
	Run() error
}

type server struct {
	sname  string
	conf   *config.Config
	finish chan bool
	lf     lockfree.LockFree
}

// New -
func New(sname string, conf *config.Config) Server {
	return &server{
		sname: sname,
		conf:  conf,
		lf:    lockfree.New(),
	}
}

// Run -
func (s *server) Run() error {
	sv := s.conf.Servers.GetServer(s.sname)
	if sv == nil {
		return errors.WithMessage(ErrNotExists, sv.Name)
	}

	sm := http.NewServeMux()
	sm.HandleFunc("/", s.RootHandler)

	go s.start(s.conf.Servers.GetIPs())

	err := http.ListenAndServe(sv.Host+":"+sv.Port, sm)
	if err != nil {
		s.finish <- true
		return err
	}

	return nil
}

func (s server) start(ips []string) {
	t := time.NewTicker(1 * time.Second)

	im, _ := icmp.New()

	deadNode := make(chan string)

END_LOOP:
	for {
		select {
		case _ = <-s.finish:
			t.Stop()
			break END_LOOP
		case _ = <-t.C:

			// Confirm the survival of servers into cluster
			for i, ip := range ips {
				s.lf.Wait()
				s.conf.Servers.SetStartingByIP(ip, true)
				s.lf.Signal()
				go im.Send(ip, i, deadNode)
			}
		case ip := <-deadNode:
			s.lf.Wait()
			s.conf.Servers.SetStartingByIP(ip, false)
			s.lf.Signal()
		}
	}
}
