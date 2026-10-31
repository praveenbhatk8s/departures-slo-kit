package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Metrics: these are what Prometheus scrapes from /metrics.
var (
	requests = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "HTTP requests by path and status code.",
	}, []string{"path", "code"})

	duration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "Request latency in seconds.",
		Buckets: prometheus.DefBuckets,
	}, []string{"path"})
)

type Departure struct {
	Line        string `json:"line"`
	Destination string `json:"destination"`
	InMinutes   int    `json:"in_minutes"`
}

// errorRate reads ERROR_RATE (0.0 to 1.0). Missing or invalid means 0: no failures.
func errorRate() float64 {
	v, err := strconv.ParseFloat(os.Getenv("ERROR_RATE"), 64)
	if err != nil {
		return 0
	}
	return v
}

func departures(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	code := http.StatusOK

	// Runs when the function returns, so it records the final status code.
	defer func() {
		requests.WithLabelValues("/departures", strconv.Itoa(code)).Inc()
		duration.WithLabelValues("/departures").Observe(time.Since(start).Seconds())
	}()

	// Fault injection: fail this share of requests on purpose.
	if rand.Float64() < errorRate() {
		code = http.StatusServiceUnavailable
		http.Error(w, "departure data unavailable", code)
		return
	}

	board := []Departure{
		{Line: "4", Destination: "Hauptbahnhof", InMinutes: 3},
		{Line: "12", Destination: "Universität", InMinutes: 7},
		{Line: "N1", Destination: "Messe", InMinutes: 15},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(board)
}

func healthz(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))
}

func main() {
	http.HandleFunc("/departures", departures)
	http.HandleFunc("/healthz", healthz)
	http.Handle("/metrics", promhttp.Handler())

	log.Printf("departures-api listening on :8080 (ERROR_RATE=%v)", errorRate())
	log.Fatal(http.ListenAndServe(":8080", nil))
}
