package collector

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"arvand-exporter/config"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

var (
	clients []config.Client
)

// SensorCollector is the list of all collectors
type SensorCollector struct {
	temperatureMetric *prometheus.Desc
	humidityMetric    *prometheus.Desc
}

// NewCollector is the main collector function
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

// Describe for metrics
func (collector *SensorCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.temperatureMetric
	ch <- collector.humidityMetric
}

// Collect data for metrics
func (collector *SensorCollector) Collect(ch chan<- prometheus.Metric) {
	start := time.Now()

	for _, client := range clients {
		url := "http://" + client.IP + "/data"

		httpClient := http.Client{
			Timeout: 5 * time.Second,
		}
		res, err := httpClient.Get(url)

		var t float64 = 0.0
		var h float64 = 0.0

		if err != nil {
			log.WithFields(log.Fields{
				"IP": client.IP,
			}).Warn("Client not reachable !")
		} else {

			body, err := ioutil.ReadAll(res.Body)

			if err != nil {
				log.WithFields(log.Fields{
					"IP": client.IP,
				}).Warn("Error fetching data !")
			} else {

				log.WithFields(log.Fields{
					"IP": client.IP,
				}).Info("Done !")

				var data *config.Sensor
				json.Unmarshal(body, &data)

				t, _ = strconv.ParseFloat(data.Temperature, 64)
				h, _ = strconv.ParseFloat(data.Humidity, 64)

				ch <- prometheus.MustNewConstMetric(collector.temperatureMetric, prometheus.GaugeValue, t, client.Name)
				ch <- prometheus.MustNewConstMetric(collector.humidityMetric, prometheus.GaugeValue, h, client.Name)
			}

		}
	}

	elapsed := time.Since(start)
	log.Info("All data collected in ", elapsed)
}
