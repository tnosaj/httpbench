package simple

import (
	"context"
	"net/http"
	"strconv"
	"time"

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

	req, err := http.NewRequest("GET", st.S.Url, nil)
	if err != nil {
		st.S.Metrics.ErrorRequests.WithLabelValues("create").Inc()
		logrus.Errorf("Error creating url: %s - %s", st.S.Url, err)
	}

	ctx, cancel := context.WithTimeout(req.Context(), time.Duration(st.S.Timeout)*time.Second)
	defer cancel()

	req = req.WithContext(ctx)

	client := http.DefaultClient
	res, err := client.Do(req)

	if err != nil {
		st.S.Metrics.ErrorRequests.WithLabelValues(strconv.Itoa(res.StatusCode)).Inc()
		logrus.Errorf("Error getting url: %s - %s", st.S.Url, err)
	}
	timer.ObserveDuration()
}
