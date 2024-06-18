package handler

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
)

type ErrorKind int

// gateway errors
const (
	Timeout ErrorKind = iota + 9001
	BadRequest
	InvalidEntity
	InternalError
	NotFoundError
	MethodNotAllowed
)

var kindToHTTPStatus = map[ErrorKind]int{
	Timeout:          http.StatusGatewayTimeout,
	InvalidEntity:    http.StatusUnprocessableEntity,
	BadRequest:       http.StatusBadRequest,
	InternalError:    http.StatusInternalServerError,
	NotFoundError:    http.StatusNotFound,
	MethodNotAllowed: http.StatusMethodNotAllowed,
}

var kindDefaultMessages = map[ErrorKind]string{
	Timeout:          "gateway timeout",
	BadRequest:       "bad request error",
	InternalError:    "internal error",
	MethodNotAllowed: "method is not allowed",
}

type IHttpError interface {
	Handle(w *http.ResponseWriter, r *http.Request)
}

type JsonError struct {
	Message string `json:"error_message"`
	Code    int    `json:"error_code"`
	Detail  string `json:"error_detail"`
}

type HttpError struct {
	Kind    ErrorKind
	Message string
	Err     error
}

func (h *HttpError) Handle(w http.ResponseWriter, r *http.Request) {
	logrus.Errorf("%v %v => Error code: \"%v\" internal error: %v", r.Method, r.URL, h.Kind, h.Error())

	httpErr := JsonError{
		Message: h.getMessage(),
		Code:    int(h.Kind),
	}

	// add detail to error
	//if !configs.Config.IsProduction() {
	//	httpErr.Detail = h.Error()
	//}

	marshaledErr, _ := json.Marshal(httpErr)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(h.getHTTPStatusCode())
	w.Write(marshaledErr)
}

func (h *HttpError) Error() string {
	if h.Err != nil {
		return h.Err.Error()
	}

	return "no internal error provided"
}

func (h *HttpError) getHTTPStatusCode() int {
	code, ok := kindToHTTPStatus[h.Kind]
	if !ok {
		return http.StatusInternalServerError
	}

	return code
}

func (h *HttpError) getMessage() string {
	if h.Message != "" {
		return h.Message
	}

	msg, ok := kindDefaultMessages[h.Kind]
	if !ok {
		// todo: check this out for nil panic
		if h.Err != nil && h.Err.Error() != "" {
			return h.Err.Error()
		}
		return "error message is not defined"
	}

	return msg
}

func NewHttpError(kind ErrorKind) *HttpError {
	return &HttpError{Kind: kind}
}

func (h *HttpError) WithErr(err error) *HttpError {
	h.Err = err
	return h
}

func (h *HttpError) WithMessage(message string) *HttpError {
	h.Message = message
	return h
}
