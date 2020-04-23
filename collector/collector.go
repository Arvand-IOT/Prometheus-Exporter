package collector

import (
	// "encoding/binary"
    // "fmt"
    // "math"
    // "os/exec"
    // "strings"

	"github.com/prometheus/client_golang/prometheus"
	// log "github.com/sirupsen/logrus"
)

// SensorCollector ...
type SensorCollector struct {
    cmdMetric *prometheus.Desc
}

// NewCollector ...
func NewCollector() *SensorCollector {
    return &SensorCollector{
        cmdMetric: prometheus.NewDesc("cmd_result",
            "Shows the cmd result",
            []string{"name"}, nil,
        ),
    }
}

// Describe ...
func (collector *SensorCollector) Describe(ch chan<- *prometheus.Desc) {
    ch <- collector.cmdMetric
}

// Collect ...
func (collector *SensorCollector) Collect(ch chan<- prometheus.Metric) {
    ch <- prometheus.MustNewConstMetric(collector.cmdMetric, prometheus.GaugeValue, 1, "rack")
    ch <- prometheus.MustNewConstMetric(collector.cmdMetric, prometheus.GaugeValue, 2, "rack2")
}
