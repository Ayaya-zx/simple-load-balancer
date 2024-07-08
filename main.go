package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/Ayaya-zx/simple-load-balancer/slb"
)

type ServerConf struct {
	URL       *url.URL `json:"url"`
	ReadyPath string   `json:"ready"`
}

func main() {
	port := flag.Int("p", 8765, "Port")

	data, err := os.ReadFile("servers.json")
	if err != nil {
		log.Fatal(err)
	}

	var serversConf []ServerConf
	err = json.Unmarshal(data, &serversConf)
	if err != nil {
		log.Fatal(err)
	}

	servers := make([]*slb.Server, 0, len(serversConf))
	for _, conf := range serversConf {
		if conf.URL == nil {
			log.Fatal("Errors in servers.json")
		}
		servers = append(servers, slb.NewServer(
			conf.URL, slb.WithReadyCheck(conf.ReadyPath),
		))
	}

	b := slb.NewBalancer(servers)
	http.HandleFunc("/", b.Handle)
	err = http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
	if err != nil {
		log.Fatal(err)
	}
}
