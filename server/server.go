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

	alive := make(chan func() (string, bool))

END_LOOP:
	for {
		select {
		case _ = <-s.finish:
			t.Stop()
			break END_LOOP
		case _ = <-t.C:

			// Confirm the survival of servers into cluster
			for i, ip := range ips {
				go im.Send(ip, i, alive)
			}
		case f := <-alive:
			ip, starting := f()
			s.lf.Wait()
			s.conf.Servers.SetStartingByIP(ip, starting)
			s.lf.Signal()
		}
	}
}
