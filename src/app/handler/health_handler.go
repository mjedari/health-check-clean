package handler

import (
	"context"
	"encoding/json"
	"github.com/mjedari/health-checker/app/contract"
	"github.com/mjedari/health-checker/domain"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strconv"
)

type HealthHandler struct {
	service contract.IHealthService
}

func NewHealthHandler(service contract.IHealthService) *HealthHandler {
	return &HealthHandler{service: service}
}

func (h *HealthHandler) Index(w http.ResponseWriter, r *http.Request) {
	endpoints, err := h.service.FetchAllEndpoints(context.Background())
	if err != nil {
		// todo: handle errors properly
		logrus.Errorf("failed to get data from service: %s", err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	JsonResponse(w, endpoints, http.StatusOK)
}

func (h *HealthHandler) Create(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)

	if err != nil {
		NewHttpError(BadRequest).
			WithErr(err).WithMessage("something went wrong").Handle(w, r)
		return
	}

	var endpoint domain.Endpoint
	err = json.Unmarshal(body, &endpoint)
	if err != nil {
		NewHttpError(InternalError).
			WithErr(err).WithMessage("something went wrong").Handle(w, r)
		return
	}

	if err := endpoint.Validate(); err != nil {
		JsonResponse(w, ValidationResponse{
			Message: "invalid data.",
			Errors:  err.Converted(),
		}, http.StatusUnprocessableEntity)
		return
	}

	err = h.service.CreateEndpoint(context.Background(), endpoint)
	if err != nil {
		// todo: handle errors properly
		logrus.Errorf("failed to get data from service: %s", err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	JsonResponse(w, endpoint, http.StatusOK)

}

func (h *HealthHandler) Start(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background() // todo: or background()
	strID := r.PathValue("id")
	id, err := strconv.Atoi(strID)
	if err != nil {
		NewHttpError(InvalidEntity).WithMessage("provided id is not valid").Handle(w, r)
		return
	}

	endpoint, err := h.service.FetchEndpoint(ctx, uint(id))
	if err != nil {
		NewHttpError(InvalidEntity).WithMessage("resource not found").Handle(w, r)
		return
	}

	// start lunching a goroutine to check its health in its interval
	err = h.service.StartWatching(ctx, endpoint)
	if err != nil {
		NewHttpError(InternalError).WithMessage("can not start watching!").Handle(w, r)
		return
	}

	JsonResponse(w, Response{Message: "watching started successfully!"}, http.StatusOK)
}

func (h *HealthHandler) Stop(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	strID := r.PathValue("id")
	id, err := strconv.Atoi(strID)
	if err != nil {
		NewHttpError(InvalidEntity).WithMessage("provided id is not valid").Handle(w, r)
		return
	}

	endpoint, err := h.service.FetchEndpoint(ctx, uint(id))
	if err != nil {
		NewHttpError(InvalidEntity).WithMessage("resource not found").Handle(w, r)
		return
	}

	err = h.service.StopWatching(ctx, endpoint)
	if err != nil {
		NewHttpError(NotFoundError).WithErr(err).Handle(w, r)
		return
	}

	JsonResponse(w, Response{Message: "watching stopped!"}, http.StatusOK)
}

func (h *HealthHandler) Delete(w http.ResponseWriter, r *http.Request) {
	// delete endpoint form database
	// remove stop its running goroutine
	strID := r.PathValue("id")
	id, err := strconv.Atoi(strID)
	if err != nil {
		NewHttpError(InvalidEntity).
			WithErr(err).WithMessage("invalid id provided").Handle(w, r)
		return
	}

	err = h.service.DeleteEndpoint(context.Background(), uint(id))
	if err != nil {
		NewHttpError(InternalError).
			WithErr(err).WithMessage("something went wrong").Handle(w, r)
		return
	}

	JsonResponse(w, Response{Message: "operation succeed!"}, http.StatusOK)
}
