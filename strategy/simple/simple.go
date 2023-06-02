package simple

import (
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"github.com/tnosaj/httpbench/internals"
)

type SimpleStrategy struct {
	S       internals.Settings
	Metrics internals.Metrics
}

func MakeSimpleStrategy(s internals.Settings) SimpleStrategy {
	logrus.Info("Simple strategy")
	RequestDuration := prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "request_duration_seconds",
		Help:    "Histogram for the runtime of a simple method function.",
		Buckets: prometheus.LinearBuckets(0.00, 0.002, 75),
	})

	ErrorReuests := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "requerst_errors",
			Help: "The total number of failed requests",
		},
		[]string{"code"},
	)

	prometheus.MustRegister(RequestDuration)
	prometheus.MustRegister(ErrorReuests)

	metrics := internals.Metrics{
		RequestDuration: RequestDuration,
		ErrorRequests:   *ErrorReuests,
	}
	return SimpleStrategy{S: s, Metrics: metrics}
}

func (st SimpleStrategy) RunCommand() {
	logrus.Debugf("ping %s", st.S.Url)
}

func (st SimpleStrategy) curl() {
	timer := prometheus.NewTimer(st.Metrics.RequestDuration)
	res, err := http.Get(st.S.Url)
	if err != nil {
		st.Metrics.ErrorRequests.WithLabelValues(strconv.Itoa(res.StatusCode)).Inc()
		logrus.Errorf("Error getting url: %s - %s", st.S.Url, err)
	}
	timer.ObserveDuration()
}
