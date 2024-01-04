package server

type ServerConfig struct {
	port         string
	storeFactory StoreFactory
}

type Server struct {
	*ServerConfig
	topics map[string]Store
}

func NewServer(scfg *ServerConfig) (*Server, error) {
	return &Server{ServerConfig: scfg}, nil
}

func (s *Server) Start() {

}
