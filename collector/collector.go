package collector

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"arvand-exporter/config"
	"github.com/prometheus/client_golang/prometheus"
	// log "github.com/sirupsen/logrus"
)

var (
	clients []config.Client
)

// SensorCollector ...
type SensorCollector struct {
	temperatureMetric *prometheus.Desc
	humidityMetric    *prometheus.Desc
}

// NewCollector ...
func NewCollector(c []config.Client) *SensorCollector {
	clients = c

	return &SensorCollector{
		temperatureMetric: prometheus.NewDesc("temperature",
			"temperature value in celsius",
			[]string{"name"}, nil,
		),
		humidityMetric: prometheus.NewDesc("humidity",
			"humidity value from 100",
			[]string{"name"}, nil,
		),
	}
}

// Describe ...
func (collector *SensorCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.temperatureMetric
	ch <- collector.humidityMetric
}

// Collect ...
func (collector *SensorCollector) Collect(ch chan<- prometheus.Metric) {
	for _, client := range clients {
		url := "http://" + client.IP + "/data"

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

		ch <- prometheus.MustNewConstMetric(collector.temperatureMetric, prometheus.GaugeValue, t, client.Name)
		ch <- prometheus.MustNewConstMetric(collector.humidityMetric, prometheus.GaugeValue, h, client.Name)
	}
}
