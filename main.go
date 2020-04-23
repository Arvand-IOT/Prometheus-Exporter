package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"encoding/json"
	"strconv"

	"arvand-exporter/collector"
	"arvand-exporter/config"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
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
	// collectData()

	// handler, err := createMetricsHandler()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// http.Handle(*metricsPath, handler)

	cmd := collector.NewCollector(cfg.Clients)
    prometheus.MustRegister(cmd)

	http.Handle("/metrics", promhttp.Handler())

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

func createMetricsHandler() (http.Handler, error) {
	// opts := collectorOptions()
	// nc, err := collector.NewCollector(cfg, opts...)
	// nc, err := collector.NewCollector()
	// if err != nil {
	// 	return nil, err
	// }

	registry := prometheus.NewRegistry()
	registry.MustRegister(TemperatureGauge)
	registry.MustRegister(HumidityGauge)

	// err = registry.Register(nc)
	// if err != nil {
	// 	return nil, err
	// }

	return promhttp.HandlerFor(registry,
		promhttp.HandlerOpts{
			ErrorLog:      log.New(),
			ErrorHandling: promhttp.ContinueOnError,
		}), nil
}

func collectData() {
	url := "http://192.168.1.29/data"

	res, err := http.Get(url)

	if err != nil {
		panic(err.Error())
	}

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		panic(err.Error())
	}

	var data *config.Sensor
	json.Unmarshal(body, &data)

	t, _ := strconv.ParseFloat(data.Temperature, 64)
	h, _ := strconv.ParseFloat(data.Humidity, 64)

	TemperatureGauge.Set(t)
	HumidityGauge.Set(h)
}

// func collectorOptions() []collector.Option {
// 	opts := []collector.Option{}

// 	opts = append(opts, collector.WithTemp())

// 	return opts
// }

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
