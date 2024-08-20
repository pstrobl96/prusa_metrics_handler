package cmd

import (
	"net/http"
	"strconv"

	"pstrobl96/prusa_metrics_handler/syslog"

	"github.com/alecthomas/kingpin/v2"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	metricsPath         = kingpin.Flag("exporter.metrics-path", "Path where to expose metrics.").Default("/metrics").String()
	metricsPort         = kingpin.Flag("exporter.metrics-port", "Port where to expose metrics.").Default("10010").Int()
	syslogListenAddress = kingpin.Flag("listen-address", "Address where to expose port for gathering metrics. - format <address>:<port>").Default("0.0.0.0:8514").String()
	influxDatabase      = kingpin.Flag("influx-database", "Database name in influx").Default("prusa").String()
	influxToken         = kingpin.Flag("influx-token", "Token for influx").Default("prusa").String()
	influxHost          = kingpin.Flag("influx-host", "Host for influx").Default("http://localhost:8086").String()
	logLevel            = kingpin.Flag("log.level", "Log level for prusa_metrics_handler.").Default("info").String()
)

// Run function to start the metrics handler
func Run() {
	kingpin.Parse()
	log.Info().Msg("Prusa metrics handler starting")

	logLevel, err := zerolog.ParseLevel(*logLevel)

	if err != nil {
		logLevel = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(logLevel)
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixNano

	log.Info().Msg("Syslog logs server starting at: " + *syslogListenAddress)
	go syslog.MetricsListener(*syslogListenAddress, *influxHost, *influxToken, *influxDatabase)

	http.Handle(*metricsPath, promhttp.Handler())
	http.ListenAndServe(":"+strconv.Itoa(*metricsPort), nil)

}
