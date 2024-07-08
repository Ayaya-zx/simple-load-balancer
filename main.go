package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/Ayaya-zx/simple-load-balancer/slb"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type ServerConf struct {
	URL   *url.URL `json:"url"`
	Ready string   `json:"ready"`
}

func main() {
	var port int
	v := viper.New()

	v.SetDefault("Port", 8765)
	v.SetDefault("Host", "localhost")

	v.AutomaticEnv()
	v.SetEnvPrefix("SLB")
	v.BindEnv("Port", "port")

	pflag.IntVarP(&port, "port", "p", 8765, "Port")
	pflag.Parse()
	if err := v.BindPFlag("Port", pflag.Lookup("port")); err != nil {
		log.Fatalln(err)
	}

	v.SetConfigName("conf")
	v.SetConfigType("json")
	v.AddConfigPath(".")
	err := v.ReadInConfig()
	if err != nil {
		log.Fatalln(err)
	}

	var serverConf []*ServerConf
	err = v.UnmarshalKey("Servers", &serverConf)
	if err != nil {
		log.Fatal(err)
	}

	servers := make([]*slb.Server, 0, len(serverConf))
	for _, conf := range serverConf {
		servers = append(servers, slb.NewServer(
			conf.URL, slb.WithReadyCheck(conf.Ready)))
	}

	b := slb.NewBalancer(servers)
	http.HandleFunc("/", b.Handle)
	err = http.ListenAndServe(fmt.Sprintf(":%d", v.GetInt("Port")), nil)
	if err != nil {
		log.Fatalln(err)
	}
}
