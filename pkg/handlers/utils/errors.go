package utils

import (
	"encoding/json"
	"net/http"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/gobuffalo/validate"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/models"
)

type errResponse struct {
	code int
	err  error
}

type clientMessage struct {
	Message string `json:"message"`
}

// errResponse creates errResponse with default headers values
func newErrResponse(code int, err error) *errResponse {
	return &errResponse{code: code, err: err}
}

// WriteResponse to the client
func (o *errResponse) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {
	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(o.code)
	json.NewEncoder(rw).Encode(clientMessage{o.err.Error()})
}

// ResponseForError logs an error and returns the expected error type
func ResponseForError(logger *zap.Logger, err error) middleware.Responder {
	switch errors.Cause(err) {
	case models.ErrFetchNotFound:
		logger.Debug("not found", zap.Error(err))
		return newErrResponse(http.StatusNotFound, err)
	case models.ErrFetchForbidden:
		logger.Debug("forbidden", zap.Error(err))
		return newErrResponse(http.StatusForbidden, err)
	case models.ErrInvalidPatchGate:
		logger.Debug("invalid patch gate", zap.Error(err))
		return newErrResponse(http.StatusBadRequest, err)
	case models.ErrInvalidTransition:
		logger.Debug("invalid transition", zap.Error(err))
		return newErrResponse(http.StatusBadRequest, err)
	default:
		logger.Error("Unexpected error", zap.Error(err))
		return newErrResponse(http.StatusInternalServerError, err)
	}
}

// ResponseForVErrors checks for validation errors
func ResponseForVErrors(logger *zap.Logger, verrs *validate.Errors, err error) middleware.Responder {
	if verrs.HasAny() {
		logger.Error("Encountered validation error", zap.Any("Validation errors", verrs.String()))
		if err == nil {
			err = errors.New("Validation error")
		}
		return newErrResponse(http.StatusBadRequest, err)
	}
	return ResponseForError(logger, err)
}

// ResponseForConflictErrors checks for conflict errors
func ResponseForConflictErrors(logger *zap.Logger, err error) middleware.Responder {
	logger.Error("Encountered conflict error", zap.Error(err))

	return newErrResponse(http.StatusConflict, err)
}
