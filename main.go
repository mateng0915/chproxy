package main

import (
	"flag"
	"github.com/hagen1778/chproxy/config"
	"github.com/hagen1778/chproxy/log"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

var (
	addr       = flag.String("h", "http://localhost:8123", "ClickHouse web-interface host:port address with scheme")
	port       = flag.String("p", ":8080", "Proxy addr to listen to for incoming requests")
	configFile = flag.String("config", "proxy.yml", "Proxy configuration filename")
)

func main() {
	log.Debugf("Loading config file from %s", *configFile)
	cfg, err := config.LoadFile(*configFile)
	if err != nil {
		log.Fatalf("can't load config %q: %s", *configFile, err)
	}
	log.Debugf("Loading config: %s", "success")

	proxy, err := NewReverseProxy(cfg)
	if err != nil {
		log.Fatalf("error while creating proxy: %s", err)
	}

	http.HandleFunc("/favicon.ico", serveFavicon)
	http.HandleFunc("/metrics", promhttp.Handler().ServeHTTP)
	http.HandleFunc("/", proxy.ServeHTTP)
	log.Fatalf("error while listening at %d: %s", *port, http.ListenAndServe(*port, nil))
}

func serveFavicon(w http.ResponseWriter, r *http.Request) {}