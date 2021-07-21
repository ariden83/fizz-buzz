package endpoint

import (
	"ariden/fizz-buzz/config"
	"ariden/fizz-buzz/internal/metrics"
	middle "ariden/fizz-buzz/internal/middleware"
	"ariden/fizz-buzz/internal/zap-graylog/logger"
	"context"
	"github.com/dimfeld/httptreemux"
	"github.com/gofrs/uuid"
	"github.com/karlseguin/ccache"
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
	log                  *zap.Logger
	metrics              *metrics.Metrics
	conf                 *config.Config
	server               *http.Server
	fetching             map[string]struct{}
	cache                *ccache.Cache // cache for valid entries
	negCache             *ccache.Cache // cache for valid entries
	cacheSize            int
	cacheTTL             int
	cacheMaxSizeAccepted int
	negCacheSize         int
	negCacheTTL          int
	queuedLock           sync.Mutex
	queued               map[string]struct{}
	fetchQueue           chan string
	fetchLock            sync.Mutex
	fetchCond            *sync.Cond
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

func New(input EndPointInput) *Endpoint {

	if input.Config.CacheSize < 10 {
		input.Config.CacheSize = 10
	}

	e := &Endpoint{
		log:                  input.Log.With(zap.String("component", "http")),
		metrics:              input.Metrics,
		conf:                 input.Config,
		cache:                ccache.New(ccache.Configure().MaxSize(int64(input.Config.CacheSize)).ItemsToPrune(uint32(input.Config.CacheSize/20) + 1)),
		negCache:             ccache.New(ccache.Configure().MaxSize(int64(input.Config.NegCacheSize)).ItemsToPrune(uint32(input.Config.NegCacheSize/20) + 1)),
		negCacheTTL:          input.Config.NegCacheTTL,
		negCacheSize:         input.Config.NegCacheSize,
		cacheSize:            input.Config.NegCacheSize,
		cacheTTL:             input.Config.NegCacheTTL,
		cacheMaxSizeAccepted: input.Config.CacheMaxSizeAccepted,

		fetchQueue: make(chan string, 1000),
		fetching:   make(map[string]struct{}),
		queued:     make(map[string]struct{}),
	}
	e.fetchCond = sync.NewCond(&e.fetchLock)
	return e
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
