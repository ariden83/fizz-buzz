package main

import (
	"fmt"
	"github.com/ariden83/fizz-buzz/benches"
	"github.com/ariden83/fizz-buzz/config"
	httpEndpoint "github.com/ariden83/fizz-buzz/internal/endpoint"
	"github.com/ariden83/fizz-buzz/internal/metrics"
	"github.com/ariden83/fizz-buzz/internal/zap-graylog/logger"
	"go.uber.org/zap"
	"testing"
)

func BenchmarkApi(b *testing.B) {
	tts := setUpBench()
	b.Run("Test GET /fizz-buzz?limit=100", tts.GetFizzBuzz100Bench)
	b.Run("Test GET /fizz-buzz?limit=1000", tts.GetFizzBuzz1000Bench)
	b.Run("Test GET /fizz-buzz?limit=10000", tts.GetFizzBuzz10000Bench)
	b.Run("Test GET /fizz-buzz?limit=100000", tts.GetFizzBuzz100000Bench)
}

type BenchServer struct {
	log     *zap.Logger
	conf    *config.Config
	metrics *metrics.Metrics
}

func setUpBench() *benches.Tests {
	conf := config.New()
	conf.Logger.Level = "ERROR"
	conf.CLILevel = "ERROR"
	l, err := logger.NewLogger(
		fmt.Sprintf("%s:%d", conf.Host, conf.Logger.Port),
		logger.Level(logger.LevelsMap[conf.Logger.Level]),
		logger.Level(logger.LevelsMap[conf.Logger.CLILevel]))
	if err != nil {
		l.Fatal("cannot setup logger")
	}

	l = l.With(zap.String("facility", conf.Name), zap.String("version", Version))
	defer l.Sync()

	m := metrics.New(conf, l)

	httpHpt := httpEndpoint.New(httpEndpoint.EndPointInput{
		Config:  conf,
		Log:     l,
		Metrics: m,
	}, httpEndpoint.WithXCache())
	tts := &benches.Tests{
		Conf:          conf,
		HTTPEndRouter: httpHpt.LoadHttpTreeMux(),
	}
	return tts
}
