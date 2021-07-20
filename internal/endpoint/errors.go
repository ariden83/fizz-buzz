package endpoint

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"net/http"
)

// GenericError Default response when we have an error
//
// swagger:response genericError
// nolint
type GenericError struct {
	// in: body
	Body ErrorResponse `json:"body"`
}

// ErrorResponse structure of error response
type ErrorResponse struct {
	// The status code
	Code int `json:"code"`
	// The error message
	Message string `json:"message"`
}

// fail Respond error to json format
func (m *Endpoint) fail(statusCode int, err error, w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(statusCode)

	if r.Header.Get("Content-type") == "application/json" {
		error := ErrorResponse{
			Message: err.Error(),
			Code:    statusCode,
		}
		js, err := json.Marshal(error)
		if err != nil {
			m.log.Error("Fail to json.Marshal in Patch method", zap.Error(err))
			return
		}
		if _, err := w.Write(js); err != nil {
			m.log.Error("Fail to Write response in http.ResponseWriter", zap.Error(err))
		}
		return
	}

	fmt.Fprint(w, err)
}
