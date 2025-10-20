package api

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humamux"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	errutil "github.com/AntonioBR9998/go-nats-simulator/errors"
	"github.com/AntonioBR9998/go-nats-simulator/gan/api/dtos"
	"github.com/AntonioBR9998/go-nats-simulator/gan/config"
	"github.com/AntonioBR9998/go-nats-simulator/gan/domain"
	"github.com/AntonioBR9998/go-nats-simulator/humamw"
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

type APIResponse[T any] struct {
	Body T `contentType:"application/json"`
}

type APIResponseWithoutBody struct{}

func NewAPI(cfg config.Config, service domain.Service) Server {
	a := &api{service: service}

	log.Traceln("creating new *mux.Router")
	r := mux.NewRouter()
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
	huma.Post(ganApi, SENSORS_ENDPOINT, a.createSensor)
	huma.Put(ganApi, SENSORS_ENDPOINT, a.modifySensor)
	huma.Get(ganApi, SENSORS_ENDPOINT, a.getSensorList, humamw.UseMiddlewares(
		humamw.UsePagination(humamw.PaginationOptions(humamw.SetMaxLimit(3000))),
		humamw.SetHeaderUsingCallback("Total"),
		humamw.UseFilter(
			ganApi,
			map[string]humamw.FilterDefinition{
				"id":    {Type: humamw.STRING},
				"type":  {Type: humamw.STRING},
				"alias": {Type: humamw.STRING},
			},
			[]string{"id", "type", "alias", "UpdatedAt"},
		),
	))
	huma.Delete(ganApi, SENSORS_ENDPOINT+"/{id:"+UUID_REGEX+"}", a.deleteSensor)

	// Metrics endpoints

	a.router = r
	return a
}

func (a *api) Router() http.Handler {
	return a.router
}

// Sensors handlers
func (a *api) createSensor(ctx context.Context, req *dtos.SensorBaseRequest) (*APIResponse[*dtos.SensorResponseBody], error) {
	// Validating type of sensor
	if !dtos.ValidateSensorType(req.Body.Type) {
		return nil, huma.NewError(400, "validation error: type must be one of temperature, humidity or pressure")
	}

	res, err := a.service.CreateSensor(ctx, req.Body.ID, req.Body.Type, req.Body.Alias, req.Body.Rate, req.Body.MaxThreshold, req.Body.MinThreshold)

	if err != nil {
		log.Errorf("error in createSensor endpoint: %v", err)
		apiErr := errutil.APIErrorHandler(err)
		return nil, huma.NewError(apiErr.GetStatus(), apiErr.Error())
	}

	return &APIResponse[*dtos.SensorResponseBody]{
		Body: dtos.ToSensorResponseDto(res),
	}, nil
}

func (a *api) modifySensor(ctx context.Context, req *dtos.SensorBaseRequest) (*APIResponse[*dtos.SensorResponseBody], error) {
	res, err := a.service.ModifySensor(ctx, req.Body.ID, req.Body.Type, req.Body.Alias, req.Body.Rate, req.Body.MaxThreshold, req.Body.MinThreshold)

	if err != nil {
		log.Errorf("error in modifySensor endpoint: %v", err)
		apiErr := errutil.APIErrorHandler(err)
		return nil, huma.NewError(apiErr.GetStatus(), apiErr.Error())
	}

	return &APIResponse[*dtos.SensorResponseBody]{
		Body: dtos.ToSensorResponseDto(res),
	}, nil
}

func (a *api) getSensorList(ctx context.Context, req *struct{}) (*APIResponse[[]*dtos.SensorResponseBody], error) {
	res, err := a.service.GetSensors(ctx)

	if err != nil {
		log.Errorf("error in getSensorList endpoint: %v", err)
		apiErr := errutil.APIErrorHandler(err)
		return nil, huma.NewError(apiErr.GetStatus(), apiErr.Error())
	}

	var sensorDtoList []*dtos.SensorResponseBody
	for _, sensor := range res {
		sensorDto := dtos.ToSensorResponseDto(sensor)
		sensorDtoList = append(sensorDtoList, sensorDto)
	}

	return &APIResponse[[]*dtos.SensorResponseBody]{
		Body: sensorDtoList,
	}, nil
}

func (a *api) deleteSensor(ctx context.Context, request *dtos.SensorRequestById) (*APIResponseWithoutBody, error) {
	err := a.service.DeleteSensor(ctx, request.Id)

	if err != nil {
		log.Errorf("error in deleteSensor endpoint: %v", err)
		apiErr := errutil.APIErrorHandler(err)
		return nil, huma.NewError(apiErr.GetStatus(), apiErr.Error())
	}

	return &APIResponseWithoutBody{}, nil
}

// Metrics handlers
