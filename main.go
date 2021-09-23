// Fizz buzz api
//
// the purpose of this application is to provide an application
// that is using plain go code to define an API
//
// This should demonstrate all the possible comment annotations
// that are available to turn go code into a fully compliant swagger 2.0 spec
//
// Github: https://github.com/ariden83/fizz-buzz
// Metrics: http://127.0.0.1:8081/metrics
//
// there are no TOS at this moment, use at your own risk we take no responsibility
//
//     Schemes: http, https
//     Host: {{url}}
//     Version: 1.0.0
//     Contact: adrienparrochia<adrienparrochia@gmail.com> http://www.citysearch-api.com
//
//     Consumes:
//     - application/json
//     - text/html
//
//     Produces:
//     - application/json
//     - text/html
//
// swagger:meta
package main

import (
	"context"
	"fmt"
	"github.com/ariden83/fizz-buzz/config"
	"github.com/ariden83/fizz-buzz/internal/metrics"
	"github.com/ariden83/fizz-buzz/internal/zap-graylog/logger"
	"os/signal"
	"time"
	// _ "github.com/go-swagger/go-swagger/cmd/swagger"
	// _ "gopkg.in/alecthomas/gometalinter.v1"
	"go.uber.org/zap"
	l "log"
	"os"
)

var Version = "0.0.0"

func main() {
	conf := config.New()
	log, err := logger.NewLogger(
		fmt.Sprintf("%s:%d", conf.Logger.Host, conf.Logger.Port),
		logger.Level(logger.LevelsMap[conf.Logger.Level]),
		logger.Level(logger.LevelsMap[conf.Logger.CLILevel]))
	if err != nil {
		l.Fatal("cannot setup logger")
	}

	hostname, err := os.Hostname()
	if err != nil {
		hostname = "N/A"
	}
	log = log.With(zap.String("facility", conf.Name), zap.String("version", Version), zap.String("instance", hostname))
	defer log.Sync()

	log.Info(conf.String())

	m := metrics.New(conf, log)

	server := &Server{
		log:     log,
		conf:    conf,
		metrics: m,
	}

	stop := make(chan error, 1)
	server.startMetricsServer(stop)
	server.startSwaggerServer(stop)
	server.startHTTPServer(stop)

	/**
	 * And wait for shutdown via signal or error.
	 */
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt)
		<-sig
		stop <- fmt.Errorf("received Interrupt signal")
	}()

	err = <-stop

	log.Error("Shutting down services", zap.Error(err))
	stopCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	server.Shutdown(stopCtx)
	log.Debug("Services shutted down")
}
