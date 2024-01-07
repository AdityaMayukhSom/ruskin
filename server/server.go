package server

import (
	"net/http"

	store "github.com/AdityaMayukhSom/ruskin/store"
)

type ServerConfig struct {
	ListenAddr   string
	StoreFactory store.StoreFactory
}

type Server struct {
	*ServerConfig
	topics map[string]store.Store
}

func NewServer(scfg *ServerConfig) (*Server, error) {
	return &Server{ServerConfig: scfg}, nil
}

func (s *Server) Start() {
	http.ListenAndServe(s.ListenAddr, nil)
}

func (s *Server) CreateTopic(name string) {
	_, ok := s.topics[name]
	if !ok {
		s.topics[name] = s.StoreFactory.Produce()
	}

}
