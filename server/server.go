package server

import (
	store "github.com/AdityaMayukhSom/ruskin/store"
)

type ServerConfig struct {
	Port         string
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

}
