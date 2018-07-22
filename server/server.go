package server

import (
	"net/http"

	"github.com/hlts2/gokvs/config"
	"github.com/pkg/errors"
)

// ErrNotExists --
var ErrNotExists = errors.New("server dose not exist")

// Server -
type Server interface {
	Run() error
}

type server struct {
	sname string
	conf  *config.Config
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
	return http.ListenAndServe(sv.Host+":"+sv.Port, sm)
}
