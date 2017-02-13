package main

func NewServer() (Server, error) {
	return Server{}, nil
}

type Server struct {
}

func (server *Server) Run() error {
	return nil
}
