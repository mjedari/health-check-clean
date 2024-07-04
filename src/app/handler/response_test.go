package handler

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestJsonResponse(t *testing.T) {
	// arrange
	rr := httptest.NewRecorder()
	data := struct {
		Name   string `json:"name"`
		Family string `json:"family"`
	}{
		"Mahdi",
		"Jedari",
	}

	invalidData := struct {
		Channel chan int `json:"channel"`
	}{
		make(chan int),
	}

	// what are edge cases?

	cases := []struct {
		name           string
		body           any
		status         int
		expectedBody   any
		expectedStatus int
	}{
		{
			"nil body",
			nil,
			http.StatusOK,
			nil,
			http.StatusOK,
		},
		{
			"invalid body",
			invalidData,
			http.StatusOK,
			"something went wrong",
			http.StatusInternalServerError,
		},
		{
			"bad request status",
			data,
			http.StatusBadRequest,
			data,
			http.StatusBadRequest,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// act
			JsonResponse(rr, c.body, c.status)

			// assert
			//body, _ := io.ReadAll(rr.Body)
			//byteData, err := json.Marshal(c.expectedBody)
			//if err != nil {
			//	//byteData = []byte(c.expectedBody.(string))
			//}

			// assert body
			//assert.Equal(t, byteData, body)
			// assert status
			assert.Equal(t, rr.Code, c.expectedStatus)
			// assert header
			assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
		})
	}
}
