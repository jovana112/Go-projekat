package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

var (
	httpHits               = createCounter("my_app_http_hit_total", "Total number of http hits.")
	httpHitsConfigPost     = createCounter("my_app_http_hit_config_post", "Total number of http hits for create config.")
	httpHitsConfigGet      = createCounter("my_app_http_hit_config_get", "Total number of http hits for get all configs.")
	httpHitsConfigDelete   = createCounter("my_app_http_hit_config_delete", "Total number of http hits for delete config.")
	httpHitsConfigGetForId = createCounter("my_app_http_hit_config_get_one", "Total number of http hits for get config.")
	httpHitsConfigUpdate   = createCounter("my_app_http_hit_config_update", "Total number of http hits for update config.")

	httpHitsGroupPost         = createCounter("my_app_http_hit_group_post", "Total number of http hits for create group.")
	httpHitsGroupGet          = createCounter("my_app_http_hit_group_get", "Total number of http hits for get all group.")
	httpHitsGroupDelete       = createCounter("my_app_http_hit_group_delete", "Total number of http hits for delete group.")
	httpHitsGroupGetForId     = createCounter("my_app_http_hit_group_get_one", "Total number of http hits for get group")
	httpHitsGroupUpdate       = createCounter("my_app_http_hit_group_update", "Total number of http hits for update group")
	httpHitsGroupExtend       = createCounter("my_app_http_hit_group_extend", "Total number of http hits for extend group")
	httpHitsGroupSearchConfig = createCounter("my_app_http_hit_group_search_config", "Total number of http hits for search config in group")

	// Add all metrics that will be resisted
	metricsList = []prometheus.Collector{
		httpHits,
		httpHitsConfigPost,
		httpHitsConfigGet,
		httpHitsConfigDelete,
		httpHitsConfigGetForId,
		httpHitsConfigUpdate,
		httpHitsGroupPost,
		httpHitsGroupGet,
		httpHitsGroupDelete,
		httpHitsGroupGetForId,
		httpHitsGroupUpdate,
		httpHitsGroupExtend,
		httpHitsGroupSearchConfig,
	}

	prometheusRegistry = prometheus.NewRegistry()
)

func createCounter(name string, help string) prometheus.Counter {
	return prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: name,
			Help: help,
		},
	)
}

func metricsHandler() http.Handler {
	prometheusRegistry.MustRegister(metricsList...)
	return promhttp.HandlerFor(prometheusRegistry, promhttp.HandlerOpts{})
}

func count(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		httpHits.Inc()
		f(w, r)
	}
}
