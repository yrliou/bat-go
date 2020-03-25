package middleware

import (
	"net/http"
	"reflect"
	"runtime"
	"strings"

	"github.com/go-chi/chi"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	latencyBuckets = []float64{.25, .5, 1, 2.5, 5, 10}
)

// InstrumentRoundTripper instruments an http.RoundTripper to capture metrics like the number
// of active requests, the total number of requests made and latency information
func InstrumentRoundTripper(roundTripper http.RoundTripper, service string) http.RoundTripper {
	inFlightGauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Name:        "client_in_flight_requests",
		Help:        "A gauge of in-flight requests for the wrapped client.",
		ConstLabels: prometheus.Labels{"service": service},
	})

	counter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name:        "client_api_requests_total",
			Help:        "A counter for requests from the wrapped client.",
			ConstLabels: prometheus.Labels{"service": service},
		},
		[]string{"code", "method"},
	)

	// dnsLatencyVec uses custom buckets based on expected dns durations.
	// It has an instance label "event", which is set in the
	// DNSStart and DNSDonehook functions defined in the
	// InstrumentTrace struct below.
	dnsLatencyVec := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:        "client_dns_duration_seconds",
			Help:        "Trace dns latency histogram.",
			Buckets:     []float64{.005, .01, .025, .05},
			ConstLabels: prometheus.Labels{"service": service},
		},
		[]string{"event"},
	)

	// tlsLatencyVec uses custom buckets based on expected tls durations.
	// It has an instance label "event", which is set in the
	// TLSHandshakeStart and TLSHandshakeDone hook functions defined in the
	// InstrumentTrace struct below.
	tlsLatencyVec := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:        "client_tls_duration_seconds",
			Help:        "Trace tls latency histogram.",
			Buckets:     []float64{.05, .1, .25, .5},
			ConstLabels: prometheus.Labels{"service": service},
		},
		[]string{"event"},
	)

	// histVec has no labels, making it a zero-dimensional ObserverVec.
	histVec := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:        "client_request_duration_seconds",
			Help:        "A histogram of request latencies.",
			Buckets:     prometheus.DefBuckets,
			ConstLabels: prometheus.Labels{"service": service},
		},
		[]string{},
	)

	// Register all of the metrics in the standard registry.
	prometheus.MustRegister(counter, tlsLatencyVec, dnsLatencyVec, histVec, inFlightGauge)

	// Define functions for the available httptrace.ClientTrace hook
	// functions that we want to instrument.
	trace := &promhttp.InstrumentTrace{
		DNSStart: func(t float64) {
			dnsLatencyVec.WithLabelValues("dns_start")
		},
		DNSDone: func(t float64) {
			dnsLatencyVec.WithLabelValues("dns_done")
		},
		TLSHandshakeStart: func(t float64) {
			tlsLatencyVec.WithLabelValues("tls_handshake_start")
		},
		TLSHandshakeDone: func(t float64) {
			tlsLatencyVec.WithLabelValues("tls_handshake_done")
		},
	}

	// Wrap the specified RoundTripper with middleware.
	return promhttp.InstrumentRoundTripperInFlight(inFlightGauge,
		promhttp.InstrumentRoundTripperCounter(counter,
			promhttp.InstrumentRoundTripperTrace(trace,
				promhttp.InstrumentRoundTripperDuration(histVec, roundTripper),
			),
		),
	)
}

// this does not work well if there are wrapped middlewares directly around a handler
// such as HTTPSignedOnly and such
func getFunctionName(f interface{}) string {
	// get the function name of what we are going to call.
	var (
		fn    = runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
		parts = strings.Split(fn, ".")
	)
	if len(parts) > 2 {
		if parts[len(parts)-1] == "func1" {
			return parts[len(parts)-2]
		}
		return parts[len(parts)-1]
	}
	return fn
}

// InstrumentHandler instruments an http.Handler to capture metrics like the number
// the total number of requests served and latency information
func InstrumentHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// get "next" function name so we don't have to specify in middleware call
		var (
			name = getFunctionName(h)
			path = chi.RouteContext(r.Context()).RoutePattern()
		)

		if name == "Mount" || strings.Contains(name, "InstrumentHandler") {
			h.ServeHTTP(w, r)
			return
		}

		hRequests := prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "api_requests_total",
				Help: "Number of requests per handler.",
				ConstLabels: prometheus.Labels{
					"handler": name,
					"path":    path,
				},
			},
			[]string{"code", "method"},
		)
		if err := prometheus.Register(hRequests); err != nil {
			if aerr, ok := err.(prometheus.AlreadyRegisteredError); ok {
				hRequests = aerr.ExistingCollector.(*prometheus.CounterVec)
			} else {
				panic(err)
			}
		}

		hLatency := prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "request_duration_seconds",
				Help:    "A histogram of latencies for requests.",
				Buckets: latencyBuckets,
				ConstLabels: prometheus.Labels{
					"handler": name,
					"path":    path,
				},
			},
			[]string{"method"},
		)
		if err := prometheus.Register(hLatency); err != nil {
			if aerr, ok := err.(prometheus.AlreadyRegisteredError); ok {
				hLatency = aerr.ExistingCollector.(*prometheus.HistogramVec)
			} else {
				panic(err)
			}
		}

		promhttp.InstrumentHandlerCounter(
			hRequests, promhttp.InstrumentHandlerDuration(
				hLatency, h)).ServeHTTP(w, r)
	})
}

// Metrics returns a http.HandlerFunc for the prometheus /metrics endpoint
func Metrics() http.HandlerFunc {
	return promhttp.Handler().(http.HandlerFunc)
}
