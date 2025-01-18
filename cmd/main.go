package cmd

import (
	"errors"
	"net/http"
	"pstrobl96/prusa_metrics_handler/handler"
	"strconv"
	"strings"

	"github.com/alecthomas/kingpin/v2"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	metricsPath         = kingpin.Flag("exporter.metrics-path", "Path where to expose metrics.").Default("/metrics").String()
	metricsPort         = kingpin.Flag("exporter.metrics-port", "Port where to expose metrics.").Default("10011").Int()
	syslogListenAddress = kingpin.Flag("listen-address", "Address where to expose port for gathering metrics. - format <address>:<port>").Default("0.0.0.0:8514").String()
	influxEnabled       = kingpin.Flag("influx.enabled", "Enable sending data to InfluxDB").Default("false").Bool()
	influxOrg           = kingpin.Flag("influx.org", "Influx organization").Default("prusa").String()
	influxBucket        = kingpin.Flag("influx.bucket", "Influx Bucket").Default("prusa").String()
	influxToken         = kingpin.Flag("influx.token", "Token for influx").Default("loremipsumdolorsitmaet").String()
	influxURL           = kingpin.Flag("influx.url", "url for influx - format http://<address>:<port>").Default("http://localhost:8086").String()
	otelEnabled         = kingpin.Flag("otel.enabled", "Enable sending data to OpenTelemetry").Default("false").Bool()
	otelEndpoint        = kingpin.Flag("otel.endpoint", "Endpoint for OpenTelemetry - format http://<address>:<port>").Default("http://localhost:4318").String()
	logLevel            = kingpin.Flag("log.level", "Log level for prusa_metrics_handler.").Default("info").String()
	prefix              = kingpin.Flag("prefix", "Prefix for metrics. Do not forget underscore!").Default("prusa_").String()
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

	configCheck, err := checkConfig()

	if !configCheck {
		log.Error().Msg("Error in configuration: " + err.Error())
		return
	}

	log.Info().Msg("Syslog metrics server starting at: " + *syslogListenAddress)

	if *influxEnabled {
		log.Info().Msg("InfluxDB enabled")
		handler.InitInfluxDB(*influxURL, *influxToken, *influxBucket, *influxOrg)
	}

	go handler.MetricsListener(*syslogListenAddress, *prefix)

	http.Handle(*metricsPath, promhttp.Handler())
	http.ListenAndServe(":"+strconv.Itoa(*metricsPort), nil)
}

func checkConfig() (bool, error) {
	if !*influxEnabled && !*otelEnabled {
		return false, errors.New("both influx and otel are disabled")
	}
	if *influxEnabled && *otelEnabled {
		return false, errors.New("both influx and otel are enabled")
	}
	if *otelEnabled && *otelEndpoint == "" {
		return false, errors.New("otel enabled but no endpoint provided")
	}
	if *influxEnabled {
		if *influxURL == "" {
			return false, errors.New("influx enabled but no URL provided")
		}
		if *influxToken == "" {
			return false, errors.New("influx enabled but no token provided")
		}
		if *influxBucket == "" {
			return false, errors.New("influx enabled but no bucket provided")
		}
		if *influxOrg == "" {
			return false, errors.New("influx enabled but no org provided")
		}
	}
	if *metricsPath == "" {
		return false, errors.New("no metrics path provided")
	}
	if *metricsPort < 1 || *metricsPort > 65535 {
		return false, errors.New("incorrect port provided")
	}
	if *syslogListenAddress == "" {
		return false, errors.New("no syslog listen address provided")
	}
	if !strings.Contains(*prefix, "_") || *prefix != "" {
		return false, errors.New("prefix without underscore")
	}
	return true, nil
}
