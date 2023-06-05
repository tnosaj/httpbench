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
}

func (st SimpleStrategy) curl() {
	timer := prometheus.NewTimer(st.S.Metrics.RequestDuration)

	client := http.Client{
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout: time.Duration(st.S.Timeout),
			}).DialContext,
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(st.S.Timeout)*time.Second)
	defer cancel()

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, st.S.Url, nil)
	res, err := client.Do(req)
	if res != nil {
		defer res.Body.Close()
	}
	// res, err := http.Get(st.S.Url)

	if e, ok := err.(net.Error); ok && e.Timeout() {
		st.S.Metrics.ErrorRequests.WithLabelValues("timeout").Inc()
		logrus.Errorf("Error getting url: %s - %s", st.S.Url, err)
	} else if err != nil {
		st.S.Metrics.ErrorRequests.WithLabelValues(strconv.Itoa(res.StatusCode)).Inc()
		logrus.Errorf("Error getting url: %s - %s", st.S.Url, err)
	}
	timer.ObserveDuration()
}
