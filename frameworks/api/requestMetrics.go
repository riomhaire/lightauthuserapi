package api

import (
	"net/http"
	"sync"
	"time"
)

// A simple way of recording per-second and per minute call tallies

type Metrics struct {
	PerSecond float32
	PerMinute float32
	PerHour   float32
}

type PerTimePeriodMetrics interface {
	Current() Metrics
	Increment()
}

type MetricsRegistry struct {
	Minutes       [60]int    // Minutes in an hour
	CurrentMinute int        // WHich is current minute (used to reset)
	mux           sync.Mutex // Prevents stepping on our feet
}

// Increment - add one to the current minute resetting tally if changes
func (r *MetricsRegistry) Increment() {
	// Secure lock
	r.mux.Lock()
	bucket := time.Now().Minute()

	// Has minute changed since last time ?
	if bucket != r.CurrentMinute {
		// Yes ... reset
		r.Minutes[bucket] = 0
		r.CurrentMinute = bucket
	}
	// Now increment
	r.Minutes[bucket] = r.Minutes[bucket] + 1

	// Release lock
	r.mux.Unlock()
}

// Current - returns current rates/tallies in calls-per-hour, per-minute, per-second
func (r *MetricsRegistry) Current() Metrics {
	bucket := time.Now().Minute()
	// Get the value
	r.mux.Lock()
	// Has minute changed since last time ?
	if bucket != r.CurrentMinute {
		// Yes ... reset
		r.Minutes[bucket] = 0
		r.CurrentMinute = bucket
	}
	tally := r.Minutes[bucket]
	// Release lock
	r.mux.Unlock()

	// OK fill out the structures
	metrics := Metrics{0.0, 0.0, 0.0}
	// If we have some work to do
	if tally > 0 {
		metrics.PerMinute = float32(tally)
		metrics.PerSecond = metrics.PerMinute / 60.0
		metrics.PerHour = metrics.PerMinute * 60.0
	}

	return metrics
}

// Handler is a MiddlewareFunc makes Stats implement the Middleware interface.
func (r *RestAPI) RecordCall(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	// Increment call
	r.MetricsRegistry.Increment()
	if next != nil {
		next(rw, req)
	}
}
