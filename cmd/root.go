package cmd

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/willfantom/goverseerr"
	"github.com/opspotes/jellyseerr-exporter/collector"
)

// persistent flags
var (
	logLevel           string
	jellyseerrAddress   string
	jellyseerrAPIKey    string
	jellyseerrAPILocale string
	fullData     bool
)

// instance to use
var jellyseerr *goverseerr.Overseerr

var RootCmd = &cobra.Command{
	Use:   "jellyseerr-exporter",
	Short: "Export request metrics from Jellyseerr",
	Long:  `Export request metrics from an Jellyseerr for Prometheus`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		setupLogger()
		logrus.WithFields(logrus.Fields{
			"command": cmd.Name(),
			"args":    args,
		}).Debugln("running command")
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		setJellyseerr()
	},
	Run: func(cmd *cobra.Command, args []string) {
		prometheus.MustRegister(prometheus.NewBuildInfoCollector())
		prometheus.MustRegister(collector.NewRequestCollector(jellyseerr, fullData))
		prometheus.MustRegister(collector.NewUserCollector(jellyseerr))

		// Handle Metrics endpoint
		promHandler := promhttp.Handler()
		http.Handle("/metrics", promHandler)

		// Default exporter redirect message on /
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`<html>
		<head><title>Jellyseerr Exporter</title></head>
		<body>
		<h1>Jellyseerr Exporter</h1>
		<p><a href="/metrics">Metrics</a></p>
		</body>
		</html>`))
		})

		if err := http.ListenAndServe(":9850", nil); err != nil {
			logrus.WithField("err msg", err.Error()).Fatalln("http server failed: exiting")
		}
	},
}

func setupLogger() {
	if level, err := logrus.ParseLevel(logLevel); err != nil {
		logrus.SetLevel(logrus.FatalLevel)
	} else {
		logrus.SetLevel(level)
	}
}

func setJellyseerr() {
	if o, err := goverseerr.NewKeyAuth(jellyseerrAddress, nil, jellyseerrAPILocale, jellyseerrAPIKey); err != nil {
		logrus.WithField("message", err.Error()).Fatalln("Could not connect to Jellyseerr")
	} else {
		jellyseerr = o
	}
}

func init() {
	RootCmd.PersistentFlags().StringVar(&logLevel, "log", "fatal", "set the log level (fatal, error, info, debug, trace)")

	// jellyseerr setup
	RootCmd.PersistentFlags().StringVar(&jellyseerrAddress, "jellyseerr.address", "", "Address at which Jellyseerr is hosted.")
	RootCmd.PersistentFlags().StringVar(&jellyseerrAPIKey, "jellyseerr.apiKey", "", "API key for admin access to the Jellyseerr instance.")
	RootCmd.PersistentFlags().StringVar(&jellyseerrAPILocale, "jellyseerr.locale", "en", "Locale of the Jellyseerr instance.")
	RootCmd.PersistentFlags().BoolVar(&fullData, "fullData", false, "Reduce scraping and cardinality on requests count metric.")
	RootCmd.MarkPersistentFlagRequired("jellyseerr.address")
	RootCmd.MarkPersistentFlagRequired("jellyseerr.api-key")
}
