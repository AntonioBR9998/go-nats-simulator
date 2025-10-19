package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humamux"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"github.com/AntonioBR9998/go-nats-simulator/gan/config"
	"github.com/AntonioBR9998/go-nats-simulator/gan/domain"
)

const (
	API_CONTEXT      = "/api"
	API_V1           = "/v1"
	API_V1_BASE      = API_CONTEXT + API_V1
	SENSORS_ENDPOINT = "/sensors"
	UUID_REGEX       = "[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}"
)

type api struct {
	router  http.Handler
	service domain.Service
}

type Server interface {
	Router() http.Handler
}

type Logger func(string)

type logResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func LogRequest(logger Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			startTime := time.Now()
			logRespWriter := &logResponseWriter{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}
			next.ServeHTTP(logRespWriter, r)
			logger(fmt.Sprintf(`"%s %s %s" - %d "%s" : %s`, r.Method, r.URL, r.Proto, logRespWriter.statusCode, r.UserAgent(), time.Since(startTime).String()))
		})
	}
}

func NewAPI(cfg config.Config, service domain.Service) Server {
	a := &api{service: service}

	log.Traceln("creating new *mux.Router")
	r := mux.NewRouter()
	log.Traceln("registering LogRequest middleware")
	r.Use(LogRequest(func(msg string) {
		log.Infoln(msg)
	}))
	log.Traceln("creating a new Subrouter for path:", API_V1_BASE)
	apiV1 := r.PathPrefix(API_V1_BASE).Subrouter()

	humaConfig := huma.DefaultConfig("GAN", "1.0.0")
	humaConfig.Servers = []*huma.Server{
		{URL: cfg.API.GetURL() + API_V1_BASE},
	}

	// This configuration allows to remove "$schema" link in JSON response
	humaConfig.CreateHooks = nil

	// V1 API Definition
	ganApi := humamux.New(apiV1, humaConfig)

	ganApi.UseMiddleware()

	// Sensors endpoints

	// Metrics endpoints

	a.router = r
	return a
}

func (a *api) Router() http.Handler {
	return a.router
}

// Handlers
