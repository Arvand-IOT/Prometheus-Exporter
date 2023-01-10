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
	airQualityMetric  *prometheus.Desc
	lightMetric       *prometheus.Desc
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
		airQualityMetric: prometheus.NewDesc("air",
			"air quality value",
			[]string{"name"}, nil,
		),
		lightMetric: prometheus.NewDesc("light",
			"light value",
			[]string{"name"}, nil,
		),
	}
}

// Describe for metrics
func (collector *SensorCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.temperatureMetric
	ch <- collector.humidityMetric
	ch <- collector.airQualityMetric
	ch <- collector.lightMetric
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
		var a float64 = 0.0
		var l float64 = 0.0

		if err != nil {
			log.WithFields(log.Fields{
				"IP": client.IP,
			}).Error("Client not reachable:")
		} else {
			body, err := ioutil.ReadAll(res.Body)

			if err != nil {
				log.WithFields(log.Fields{
					"IP": client.IP,
				}).Error("Error fetching data:")
			} else {
				log.WithFields(log.Fields{
					"IP": client.IP,
				}).Info("Done:")

				var data *config.Sensor
				json.Unmarshal(body, &data)

				t, _ = strconv.ParseFloat(data.Temperature, 64)
				h, _ = strconv.ParseFloat(data.Humidity, 64)
				a, _ = strconv.ParseFloat(data.AirQuality, 64)
				l, _ = strconv.ParseFloat(data.Light, 64)

				ch <- prometheus.MustNewConstMetric(collector.temperatureMetric, prometheus.GaugeValue, t, client.Name)
				ch <- prometheus.MustNewConstMetric(collector.humidityMetric, prometheus.GaugeValue, h, client.Name)
				ch <- prometheus.MustNewConstMetric(collector.airQualityMetric, prometheus.GaugeValue, a, client.Name)
				ch <- prometheus.MustNewConstMetric(collector.lightMetric, prometheus.GaugeValue, l, client.Name)
			}
		}
	}

	elapsed := time.Since(start)
	log.Info("All data collected in ", elapsed)
}
