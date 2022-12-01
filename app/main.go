package main

import (
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	http_response_status = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "sample_external_url_up",
			Help: "http status",
		},
		[]string{"url"},
	)

	http_response_time_milliseconds = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "sample_external_url_response_ms",
			Help: "http response in ms",
		},
		[]string{"url"},
	)
)

func recordMetrics(URL string) error {
	start := time.Now()
	resp, err := http.Get(URL)
	if err != nil {
		return err
	}
	duration := time.Since(start)
	http_response_time_milliseconds.WithLabelValues(URL).Set(float64(duration.Milliseconds()))
	if resp.StatusCode == http.StatusOK {
		http_response_status.WithLabelValues(URL).Set(1)
	} else {
		http_response_status.WithLabelValues(URL).Set(0)
	}
	resp.Body.Close()
	return nil
}

func main() {
	prometheus.MustRegister(http_response_status)
	prometheus.MustRegister(http_response_time_milliseconds)

	URLs := []string{"https://httpstat.us/200", "https://httpstat.us/503"}
	go func() {
		for {
			for _, u := range URLs {
				if err := recordMetrics(u); err != nil {
					break
				}
			}
			time.Sleep(time.Minute * 1)
		}
	}()

	http.Handle("/metrics", promhttp.Handler())

	log.Println("Server Running...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
