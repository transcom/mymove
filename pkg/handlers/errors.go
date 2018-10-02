package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/gobuffalo/validate"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/models"
	uploaderpkg "github.com/transcom/mymove/pkg/uploader"
)

// ValidationErrorsResponse is a middleware.Responder for a set of validation errors
type ValidationErrorsResponse struct {
	Errors map[string]string `json:"errors,omitempty"`
}

func newValidationErrorsResponse(errors map[string]string) *ValidationErrorsResponse {
	return &ValidationErrorsResponse{Errors: errors}
}

// WriteResponse to the client
func (v *ValidationErrorsResponse) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {
	rw.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(rw).Encode(v)
}

// ErrResponse collect errors and error codes
type ErrResponse struct {
	Code int
	Err  error
}

type clientMessage struct {
	Message string `json:"message"`
}

// ErrResponse creates ErrResponse with default headers values
func newErrResponse(code int, err error) *ErrResponse {
	return &ErrResponse{Code: code, Err: err}
}

// WriteResponse to the client
func (o *ErrResponse) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {
	rw.WriteHeader(o.Code)
	json.NewEncoder(rw).Encode(clientMessage{o.Err.Error()})
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
	case models.ErrUserUnauthorized:
		logger.Debug("unauthorized", zap.Error(err))
		return newErrResponse(http.StatusUnauthorized, err)
	case uploaderpkg.ErrZeroLengthFile:
		logger.Debug("uploaded zero length file", zap.Error(err))
		return newErrResponse(http.StatusBadRequest, err)
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
		errors := make(map[string]string)
		for _, key := range verrs.Keys() {
			errors[key] = strings.Join(verrs.Get(key), " ")
		}
		return newValidationErrorsResponse(errors)
	}
	return ResponseForError(logger, err)
}

// ResponseForConflictErrors checks for conflict errors
func ResponseForConflictErrors(logger *zap.Logger, err error) middleware.Responder {
	logger.Error("Encountered conflict error", zap.Error(err))

	return newErrResponse(http.StatusConflict, err)
}
