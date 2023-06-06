package simple

import (
	"context"
	"net"
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
	logrus.Debugf("pong %s", st.S.Url)
}

func (st SimpleStrategy) curl() {
	timer := prometheus.NewTimer(st.S.Metrics.RequestDuration)

	client := http.Client{}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(st.S.Timeout)*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, st.S.Url, nil)
	if err != nil {
		logrus.Errorf("error creating request %s", err)
		return
	}

	res, err := client.Do(req)
	if res != nil {
		defer res.Body.Close()
	}
	// res, err := http.Get(st.S.Url)

	if e, ok := err.(net.Error); ok && e.Timeout() {
		st.S.Metrics.ErrorRequests.WithLabelValues("timeout").Inc()
		logrus.Errorf("Error getting url timeout: %s - %s", st.S.Url, err)
	} else if err == nil && res.StatusCode != http.StatusOK {
		st.S.Metrics.ErrorRequests.WithLabelValues(strconv.Itoa(res.StatusCode)).Inc()
		logrus.Errorf("Error getting url reponsecode: %s - %s", st.S.Url, err)
	} else if err != nil {
		st.S.Metrics.ErrorRequests.WithLabelValues("unknown").Inc()
		logrus.Errorf("Error from an unknown source: %s", err)
	}
	timer.ObserveDuration()
}
