package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

var (
	addr1   = "http://localhost:8080"
	addr2   = "http://localhost:8081"
	counter = 0
)

type Config struct {
	proxyPort string
	sliceHost []string
	counter   int
}

type Server struct {
	config *Config
}

func NewServer(config *Config) *Server {
	return &Server{config: config}
}

func (s *Server) Run() {
	http.HandleFunc("/", s.handle)
	log.Fatalln(http.ListenAndServe(s.config.proxyPort, nil))
}

func (s *Server) handle(w http.ResponseWriter, r *http.Request) {
	var urlreq *url.URL
	var err error
	if s.config.counter == len(s.config.sliceHost) {
		s.config.counter = 0
	}
	urlreq, err = url.Parse(s.config.sliceHost[s.config.counter])
	if err != nil {
		w.Write([]byte(err.Error()))
	}
	s.config.counter++
	proxy := httputil.NewSingleHostReverseProxy(urlreq)
	proxy.ServeHTTP(w, r)
}

func main() {
	config := &Config{":9000", []string{"http://localhost:8080", "http://localhost:8081"}, 0}
	server := NewServer(config)
	server.Run()
}
