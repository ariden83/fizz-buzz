package main

import (
	"ariden/fizz-buzz/benches"
	"ariden/fizz-buzz/config"
	httpEndpoint "ariden/fizz-buzz/internal/endpoint"
	"ariden/fizz-buzz/internal/metrics"
	"ariden/fizz-buzz/internal/zap-graylog/logger"
	"fmt"
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
	})
	tts := &benches.Tests{
		Conf:          conf,
		HTTPEndRouter: httpHpt.LoadHttpTreeMux(),
	}
	return tts
}
