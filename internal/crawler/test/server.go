package test

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
)

type Server struct {
	Ts *httptest.Server
}

func NewServer() Server {
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal("failed to get pwd")
	}
	return Server{
		Ts: httptest.NewServer(http.FileServer(http.Dir(pwd + "/test/html"))),
	}
}

func (s *Server) Stop() {
	s.Ts.Close()
}
