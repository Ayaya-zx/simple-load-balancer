package slb

import (
	"net/http"
	"time"
)

type Balancer struct {
	servers []*Server
	count   int
}

func NewBalancer(servers []*Server) *Balancer {
	return &Balancer{servers: servers}
}

func (b *Balancer) Handle(w http.ResponseWriter, r *http.Request) {
	if len(b.servers) == 0 {
		return
	}
	server := b.getNextServer()
	start := server
	for !server.isReady() {
		server = b.getNextServer()
		if server == start {
			time.Sleep(time.Second)
		}
	}
	server.handle(w, r)
}

func (b *Balancer) getNextServer() *Server {
	server := b.servers[b.count]
	b.count++
	if b.count >= len(b.servers) {
		b.count = 0
	}
	return server
}
