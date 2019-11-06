package main

import (
	"log"
	"net/http"
	"net/http/httptest"
)

const html_folder_path = "/Users/sterlingdeng/Software/go/webcrawler/cmd/test/html"

type Server struct {
	ts *httptest.Server
}

func NewServer() Server {
	return Server{
		ts: httptest.NewServer(http.FileServer(http.Dir(html_folder_path))),
	}
}

func (s *Server) Start() {
	s.ts.Start()
}

func (s *Server) Stop() {
	s.ts.Close()
}

func main() {
	h := http.FileServer(http.Dir(html_folder_path))
	err := http.ListenAndServe(":8080", h)
	if err != nil {
		log.Fatal(err)
	}
}
