package server

import (
	"github.com/sirupsen/logrus"
	"github.com/tnosaj/httpbench/internals"
)

type HttpbenchServer struct {
	Settings internals.Settings
}

type HttpSettings struct {
	Strategy    string `json:"strategy"`
	Concurrency int    `json:"concurrency"`
	Duration    int    `json:"duration"`
	Rate        int    `json:"rate"`
	Url         string `json:"url"`
}

func NewHttpbenchServer(settings internals.Settings) HttpbenchServer {
	logrus.Info("started server")

	return HttpbenchServer{Settings: settings}
}
