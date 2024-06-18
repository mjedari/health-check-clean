package handler

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
)

type Response struct {
	Message string `json:"message"`
}

type ValidationResponse struct {
	Message string
	Errors  []string
}

func JsonResponse(w http.ResponseWriter, data any, status int) {

	res, err := json.Marshal(data)
	if err != nil {
		logrus.Errorf("failed to marshaling response: %s", err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(res)
}
