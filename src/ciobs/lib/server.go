// server.go
package ciobs

import (
	"log"
)

type Server struct {
	manager *TaskMan
}

func NewServer() (s *Server, e error) {
	tm, e := NewTaskMan()
	if e != nil {
		return nil, e
	}

	return &Server{tm}, nil
}

func StartServer() {
	_, e := NewServer()
	if e != nil {
		log.Fatalf("Failed to load server; %s", e)
	}
}
