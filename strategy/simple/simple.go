package simple

import (
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"github.com/tnosaj/httpbench/internals"
)

type SimpleStrategy struct {
	S internals.Settings
}

func MakeSimpleStrategy(s internals.Settings) SimpleStrategy {
	logrus.Info("Simple strategy")
	return SimpleStrategy{S: s}
}

func (st SimpleStrategy) RunCommand() {
	logrus.Debugf("ping %s", st.S.Url)
	st.curl()
}

func (st SimpleStrategy) curl() {
	timer := prometheus.NewTimer(st.S.Metrics.RequestDuration)
	res, err := http.Get(st.S.Url)
	if err != nil {
		st.S.Metrics.ErrorRequests.WithLabelValues(strconv.Itoa(res.StatusCode)).Inc()
		logrus.Errorf("Error getting url: %s - %s", st.S.Url, err)
	}
	timer.ObserveDuration()
}
