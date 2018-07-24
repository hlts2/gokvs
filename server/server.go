package server

import (
	"net/http"
	"time"

	"github.com/hlts2/gokvs/config"
	"github.com/hlts2/gokvs/icmp"
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
}

// New -
func New(sname string, conf *config.Config) Server {
	return &server{
		sname: sname,
		conf:  conf,
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

	go s.start(s.conf.Servers.GetHostAndPorts())

	err := http.ListenAndServe(sv.Host+":"+sv.Port, sm)
	if err != nil {
		s.finish <- true
		return err
	}

	return nil
}

func (s server) start(ips []string) {
	t := time.NewTicker(1 * time.Second)

	icmp, _ := icmp.New()

END_LOOP:
	for {
		select {
		case _ = <-s.finish:
			t.Stop()
			break END_LOOP
		case _ = <-t.C:

			// Confirm the survival of servers into cluster
			for i, ip := range ips {
				go icmp.Send(ip, i)
			}
		}
	}
}
