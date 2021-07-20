package main

import (
	"ariden/fizz-buzz/config"
	"ariden/fizz-buzz/internal/metrics"
	"ariden/fizz-buzz/internal/zap-graylog/logger"
	"ariden/fizz-buzz/tests"
	"fmt"
	"go.uber.org/zap"
	"testing"
	"time"
)

// TestApi
func TestApi(t *testing.T) {
	c := setUpTest()
	tts := &tests.Tests{
		Conf:       c,
		DefaultURL: fmt.Sprintf("http://%s:%d", c.Host, c.Port),
	}
	time.Sleep(1500 * time.Millisecond)
	tts.StartFunctionnalTests(t)
	t.Run("HealthCheck", tts.HealthCheckTest)
	t.Run("Metrics", tts.MetricsTest)
	t.Run("Test GET /fizz-buzz", tts.GetFizzBuzzTest)
}

func setUpTest() *config.Config {
	conf := config.New()

	l, err := logger.NewLogger(
		fmt.Sprintf("%s:%d", conf.Host, conf.Logger.Port),
		logger.Level(logger.LevelsMap[conf.Logger.Level]),
		logger.Level(logger.LevelsMap[conf.CLILevel]))
	if err != nil {
		l.Fatal("cannot setup logger")
	}
	l = l.With(zap.String("facility", conf.Name), zap.String("version", Version))
	defer l.Sync()

	l.Info(fmt.Sprintf("%#v", conf))

	m := metrics.New(conf, l)

	server := &Server{
		log:     l,
		conf:    conf,
		metrics: m,
	}

	stop := make(chan error, 1)
	go server.startMetricsServer(stop)
	//	go server.startSwaggerRoutes(stop)
	go server.startHTTPServer(stop)

	return conf
}
