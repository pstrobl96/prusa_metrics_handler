package syslog

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	metricsHandlerTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name:        "prusa_metrics_handler_syslog_message_total",
		Help:        "The total number of processed events",
		ConstLabels: prometheus.Labels{"type": "syslog"},
	})
)
