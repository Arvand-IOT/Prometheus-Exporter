package main

import (
	"os"
	"fmt"
	"flag"
	"net/http"
	"io/ioutil"
	"bytes"
	// "encoding/json"
	
	"arvand-exporter/config"
	// "arvand-exporter/collector"

	log "github.com/sirupsen/logrus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	configFile  = flag.String("config-file", "", "config file to load")
	metricsPath = flag.String("path", "/metrics", "path to answer requests on")
	port        = flag.String("port", ":9437", "port number to listen on")
	logLevel    = flag.String("log-level", "info", "log level")
	logFormat   = flag.String("log-format", "json", "logformat text or json (default json)")

	cfg *config.Config
)

var (
	// TemperatureGauge is the temperature of sensor in celsius
	TemperatureGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "temperature",
		Help: "Current temperature of the sensor.",
	})

	// HumidityGauge is the humidity of sensor ( 0 - 100 )
	HumidityGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "humidity",
		Help: "Current humidity of the sensor.",
	})
)

func main() {
	flag.Parse()

	configureLog()

	log.Info("Welcome to Arvand Prometheus Exporter")

	log.Info("Version : 1.0")

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
		log.SetFormatter(&log.TextFormatter{})
	} else {
		log.SetFormatter(&log.JSONFormatter{})
	}
}

func loadConfig() (*config.Config, error) {
	if *configFile != "" {
		return loadConfigFromFile()
	}

	return nil, fmt.Errorf("missing config file")
}

func loadConfigFromFile() (*config.Config, error) {
	b, err := ioutil.ReadFile(*configFile)
	if err != nil {
		return nil, err
	}

	return config.Load(bytes.NewReader(b))
}

func startServer() {
	TemperatureGauge.Set(65.3)
	HumidityGauge.Set(64.3)

	handler, err := createMetricsHandler()
	if err != nil {
		log.Fatal(err)
	}

	http.Handle(*metricsPath, handler)

	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("ok"))
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`<html>
			<head><title>Mikrotik Exporter</title></head>
			<body>
			<h1>Mikrotik Exporter</h1>
			<p><a href="` + *metricsPath + `">Metrics</a></p>
			</body>
			</html>`))
	})

	log.Info("Listening on ", *port)
	log.Fatal(http.ListenAndServe(*port, nil))
}

func createMetricsHandler() (http.Handler, error) {
	registry := prometheus.NewRegistry()
	registry.MustRegister(TemperatureGauge)
	registry.MustRegister(HumidityGauge)

	return promhttp.HandlerFor(registry,
		promhttp.HandlerOpts{
			ErrorLog:      log.New(),
			ErrorHandling: promhttp.ContinueOnError,
		}), nil
}

// func main() {
// 	var conf Conf
// 	var data Sensor
//     conf.getConf()

// 	for _, client := range conf.Clients {
// 		data = get(client.IP)
// 		fmt.Println(client.IP , " : " , data.Temperature)
// 	}

// 	os.Exit(1)
// }
