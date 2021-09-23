package endpoint

import (
	"context"
	"github.com/ariden83/fizz-buzz/config"
	"github.com/ariden83/fizz-buzz/internal/metrics"
	middle "github.com/ariden83/fizz-buzz/internal/middleware"
	"github.com/ariden83/fizz-buzz/internal/xcache"
	"github.com/ariden83/fizz-buzz/internal/zap-graylog/logger"
	"github.com/dimfeld/httptreemux"
	"github.com/gofrs/uuid"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/urfave/negroni"
	"go.uber.org/zap"
	"net/http"
	"strings"
	"sync"
	"time"
)

type Endpoint struct {
	log        *zap.Logger
	metrics    *metrics.Metrics
	conf       *config.Config
	server     *http.Server
	fetching   map[string]struct{}
	xcache     *xcache.Cache // cache for valid entries
	queuedLock sync.Mutex
	queued     map[string]struct{}
	fetchQueue chan string
	fetchLock  sync.Mutex
	fetchCond  *sync.Cond
}

const (
	RequestIDHeaderKey = "X-Request-ID"
	RequestIDKey       = "RequestID"
	ContentTypeJSON    = "application/json"
)

type EndPointInput struct {
	Config  *config.Config
	Log     *zap.Logger
	Metrics *metrics.Metrics
}

// Option is the type of option passed to the constructor.
type Option func(e *Endpoint)

func New(input EndPointInput, opts ...Option) *Endpoint {

	if input.Config.CacheSize < 10 {
		input.Config.CacheSize = 10
	}

	e := &Endpoint{
		log:        input.Log.With(zap.String("component", "http")),
		metrics:    input.Metrics,
		conf:       input.Config,
		fetchQueue: make(chan string, 1000),
		fetching:   make(map[string]struct{}),
		queued:     make(map[string]struct{}),
	}
	e.fetchCond = sync.NewCond(&e.fetchLock)

	for _, o := range opts {
		o(e)
	}

	return e
}

func WithXCache() Option {
	return func(s *Endpoint) {
		var err error

		s.xcache, err = xcache.New(
			xcache.WithSize(int32(s.conf.CacheSize)),
			xcache.WithTTL(time.Duration(s.conf.CacheTTL)*time.Second),
			xcache.WithNegSize(int32(s.conf.NegCacheSize)),
			xcache.WithNegTTL(time.Duration(s.conf.NegCacheTTL)*time.Second),
			xcache.WithStale(true),
			xcache.WithPruneSize(int32(s.conf.CacheSize/20)+1))

		if err != nil {
			s.log.Error("fail to init xcache", zap.Error(err))
		}
	}
}

func (s *Endpoint) RequestIDHeader(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var reqID string
	if r.Header.Get(RequestIDHeaderKey) == "" {
		u, _ := uuid.NewV4()
		reqID = u.String()
	} else {
		u2, err := uuid.FromString(r.Header.Get(RequestIDHeaderKey))
		if err != nil {
			u, _ := uuid.NewV4()
			reqID = u.String()
		} else {
			reqID = u2.String()
		}
	}

	w.Header().Set(RequestIDHeaderKey, reqID)
	ctx := context.WithValue(r.Context(), RequestIDKey, reqID)
	ctx = logger.ToContext(ctx, s.log.With(zap.String(RequestIDKey, reqID)))
	next(w, r)
}

func (s *Endpoint) Shutdown(ctx context.Context) {
	s.log.Debug("Gracefully pausing down the HTTP server", zap.String("address", s.server.Addr))
	s.server.Shutdown(ctx)
}

func (s *Endpoint) LoadHttpTreeMux() *negroni.Negroni {
	mux := httptreemux.New()

	mux.Handle("GET", "/fizz-buzz", s.GetFizzBuzz)

	n := negroni.New(negroni.HandlerFunc(middle.DefaultHeader))
	n.UseFunc(s.RequestIDHeader)

	n.Use(negroni.HandlerFunc(func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		route := strings.ToLower(r.Method)

		jsonHandler := promhttp.InstrumentHandlerInFlight(
			s.metrics.InFlight,

			promhttp.InstrumentHandlerResponseSize(
				s.metrics.ResponseSize.MustCurryWith(prometheus.Labels{"service": route}),

				promhttp.InstrumentHandlerRequestSize(
					s.metrics.RequestSize.MustCurryWith(prometheus.Labels{"service": route}),

					promhttp.InstrumentHandlerCounter(
						s.metrics.RouteCountReqs.MustCurryWith(prometheus.Labels{"service": route}),

						promhttp.InstrumentHandlerDuration(
							s.metrics.ResponseDuration.MustCurryWith(prometheus.Labels{"service": route}),
							next)))))

		jsonHandler.ServeHTTP(rw, r)
	}))

	n.UseHandler(mux)

	return n
}

func (s *Endpoint) Listen(address string) error {
	n := s.LoadHttpTreeMux()

	s.log.Info("Listening HTTP server", zap.String("address", address))
	s.server = &http.Server{
		Addr:         address,
		Handler:      n,
		ReadTimeout:  s.conf.APIReadTimeout * time.Second,
		WriteTimeout: s.conf.APIWriteTimeout * time.Second,
	}
	if err := s.server.ListenAndServe(); err != nil {
		return err
	}
	return nil
}
