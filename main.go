package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"arvand-exporter/collector"
	"arvand-exporter/config"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

var (
	configFile  = flag.String("config-file", "", "The config file to load")
	metricsPath = flag.String("path", "/metrics", "Path to answer requests on ( Default /metrics )")
	port        = flag.String("port", ":9437", "Port number to listen on ( Default: 9437 )")
	logLevel    = flag.String("log-level", "info", "Log level ( Default: info )")
	logFormat   = flag.String("log-format", "text", "Log format. text / json ( Default: text )")

	cfg *config.Config
)

func main() {
	flag.Parse()

	configureLog()

	log.Info("Welcome to Arvand Prometheus exporter")

	log.Info("Version: 1.2.0")

	c, err := loadConfig()
	if err != nil {
		log.Errorf("Could not load config : %v", err)
		os.Exit(3)
	}
	cfg = c

	startServer()
}

func configureLog() {
	ll, err := log.ParseLevel(*logLevel)
	if err != nil {
		panic(err)
	}

	log.SetLevel(ll)

	if *logFormat == "text" {
		log.SetFormatter(&log.TextFormatter{
			FullTimestamp: true,
		})
	} else {
		log.SetFormatter(&log.JSONFormatter{})
	}
}

func loadConfig() (*config.Config, error) {
	if *configFile != "" {
		return loadConfigFromFile()
	}

	return nil, fmt.Errorf("Missing config file")
}

func loadConfigFromFile() (*config.Config, error) {
	b, err := ioutil.ReadFile(*configFile)
	if err != nil {
		return nil, err
	}

	return config.Load(bytes.NewReader(b))
}

func startServer() {
	collector := collector.NewCollector(cfg.Clients)
	registry := prometheus.NewRegistry()
	registry.MustRegister(collector)
	handler := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})

	http.Handle("/metrics", handler)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`<html>
			<head><title>Arvand Exporter</title></head>
			<body>
			<h1>Arvand Exporter</h1>
			<p><a href="` + *metricsPath + `">Metrics</a></p>
			</body>
			</html>`))
	})

	log.Info("Listening on ", *port)
	log.Fatal(http.ListenAndServe(*port, nil))
}
