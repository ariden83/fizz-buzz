package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"context"
	"github.com/ariden83/fizz-buzz/config"
	httpEndpoint "github.com/ariden83/fizz-buzz/internal/endpoint"
	"github.com/ariden83/fizz-buzz/internal/metrics"
	"github.com/juju/errors"
	promhttp "github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/urfave/negroni"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http/pprof"
	"os"
	"strings"
	"time"
)

type Healthz struct {
	Result   bool     `json:"result"`
	Messages []string `json:"messages"`
	Version  string   `json:"version"`
}

type Server struct {
	log           *zap.Logger
	conf          *config.Config
	metrics       *metrics.Metrics
	httpServer    *httpEndpoint.Endpoint
	swaggerServer *http.Server
	metricsServer *http.Server
}

func (s *Server) startHTTPServer(stop chan error) {
	go func() {
		s.httpServer = httpEndpoint.New(httpEndpoint.EndPointInput{
			Config:  s.conf,
			Log:     s.log,
			Metrics: s.metrics,
		}, httpEndpoint.WithXCache())
		if err := s.httpServer.Listen(fmt.Sprintf("%s:%d", s.conf.Host, s.conf.Port)); err != nil {
			stop <- errors.Annotate(err, "cannot start server HTTP")
		}
	}()
}

func (s *Server) getPublicSwaggerURL() string {
	url := s.conf.PublicURL
	i := strings.Index(url, "http")
	if i == -1 {
		return url
	}

	re := regexp.MustCompile("^http(s)://")
	return re.ReplaceAllString(url, "")
}

func (s *Server) generateSwaggerFile(rootDir string) error {
	buff, err := ioutil.ReadFile(rootDir + "/swagger/swagger-template.json")
	if err != nil {
		return err
	}
	r := strings.Replace(string(buff), "{{url}}", s.getPublicSwaggerURL(), -1)
	r = strings.Replace(r, "{{name}}", s.conf.Name, -1)
	r = strings.Replace(r, "{{version}}", Version, -1)

	swaggerYaml := []byte(r)
	f, err := os.Create(rootDir + "/swagger/swagger.json")
	if err != nil {
		return err
	}
	if _, err := f.Write(swaggerYaml); err != nil {
		return err
	}
	return f.Sync()
}

func (s *Server) startSwaggerServer(stop chan error) {
	rootDir, _ := os.Getwd()
	err := s.generateSwaggerFile(rootDir)
	if err != nil {
		s.log.Fatal("Fail to generate swagger file", zap.Error(err))
		return
	}
	mux := http.NewServeMux()

	n := negroni.New()
	n.Use(negroni.NewStatic(http.Dir(rootDir + "/swagger")))
	n.UseHandler(mux)

	address := fmt.Sprintf("%s:%d", s.conf.Swagger.Host, s.conf.Swagger.Port)
	s.swaggerServer = &http.Server{
		Addr:           address,
		Handler:        n,
		ReadTimeout:    time.Duration(s.conf.Healthz.ReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(s.conf.Healthz.WriteTimeout) * time.Second,
		IdleTimeout:    120 * time.Second,
		MaxHeaderBytes: 1 << 12,
	}
	go func() {
		s.log.Info("Listening HTTP for swagger route", zap.String("address", address))
		if err := s.swaggerServer.ListenAndServe(); err != nil {
			s.log.Fatal("StartServer", zap.Error(err))
			stop <- errors.Annotate(err, "cannot start swagger server")
		}
	}()
}

func (s *Server) Shutdown(ctx context.Context) {
	s.httpServer.Shutdown(ctx)
	s.metricsServer.Shutdown(ctx)
	s.swaggerServer.Shutdown(ctx)
}

func (s *Server) startMetricsServer(stop chan error) {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		message := "The service " + s.conf.Name + " responds correctly"
		res := Healthz{Result: true, Messages: []string{message}, Version: Version}
		js, err := json.Marshal(res)
		if err != nil {
			s.log.Fatal("Fail to jsonify", zap.Error(err))
		}
		if _, err := w.Write(js); err != nil {
			s.log.Fatal("Fail to Write response in http.ResponseWriter", zap.Error(err))
			return
		}
	})

	mux.HandleFunc("/readiness", func(w http.ResponseWriter, r *http.Request) {
		result := true
		message := "The service " + s.conf.Name + " responds correctly"

		res := Healthz{Result: result, Messages: []string{message}, Version: Version}
		js, err := json.Marshal(res)
		if err != nil {
			s.log.Fatal("Fail to jsonify", zap.Error(err))
		}
		if _, err := w.Write(js); err != nil {
			s.log.Fatal("Fail to Write response in http.ResponseWriter", zap.Error(err))
			return
		}
	})

	mux.Handle("/metrics", promhttp.Handler())
	s.PProf(mux)

	addr := fmt.Sprintf("%s:%d", s.conf.Metrics.Host, s.conf.Metrics.Port)
	s.metricsServer = &http.Server{
		Addr:           addr,
		Handler:        mux,
		ReadTimeout:    time.Duration(s.conf.Healthz.ReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(s.conf.Healthz.WriteTimeout) * time.Second,
		IdleTimeout:    120 * time.Second,
		MaxHeaderBytes: 1 << 12,
	}
	go func() {
		s.log.Info("Listening HTTP for healthz route", zap.String("address", addr))
		if err := s.metricsServer.ListenAndServe(); err != nil {
			stop <- errors.Annotate(err, "cannot start healthz server")
		}
	}()
}

func (s *Server) PProf(mux *http.ServeMux) {
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.Handle("/debug/pprof/allocs", pprof.Handler("allocs"))
	mux.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	mux.Handle("/debug/pprof/mutex", pprof.Handler("mutex"))
	mux.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	mux.Handle("/debug/pprof/block", pprof.Handler("block"))
	mux.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
	mux.Handle("/debug/pprof/trace", pprof.Handler("trace"))
}
