package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/lib/pq"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/gobuffalo/validate"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/route"
	uploaderpkg "github.com/transcom/mymove/pkg/uploader"
)

// SQLErrMessage represents string value to represent generic sql error to avoid leaking implementation details
const SQLErrMessage string = "Unhandled SQL error encountered"

// NotFoundMessage string value to represent sql not found
const NotFoundMessage string = "Not Found Error"

// NilErrMessage indicates an uninstantiated error was passed
const NilErrMessage string = "Nil error passed"

// ConflictErrMessage indicates that there was a conflict with input values
const ConflictErrMessage string = "Conflict Error"

// PreconditionErrMessage indicates that the IfMatch header (eTag) was stale
const PreconditionErrMessage string = "Precondition Failed"

// BadRequestErrMessage indicates that the request was malformed
const BadRequestErrMessage string = "Bad Request"

// ValidationErrMessage indicates that some fields were invalid
const ValidationErrMessage string = "Validation Error"

// ValidationErrorListResponse maps field names to a list of errors for the field
type ValidationErrorListResponse struct {
	Errors map[string][]string `json:"errors,omitempty"`
}

// NewValidationErrorListResponse returns a new validations error list response
func NewValidationErrorListResponse(verrs *validate.Errors) *ValidationErrorListResponse {
	errorList := make(map[string][]string)
	for _, key := range verrs.Keys() {
		errorList[key] = verrs.Get(key)
	}
	return &ValidationErrorListResponse{Errors: errorList}
}

// ValidationErrorsResponse is a middleware.Responder for a set of validation errors
type ValidationErrorsResponse struct {
	Errors map[string]string `json:"errors,omitempty"`
}

// NewValidationErrorsResponse returns a new validations errors response
func NewValidationErrorsResponse(verrs *validate.Errors) *ValidationErrorsResponse {
	errors := make(map[string]string)
	for _, key := range verrs.Keys() {
		errors[key] = strings.Join(verrs.Get(key), " ")
	}
	return &ValidationErrorsResponse{Errors: errors}
}

// WriteResponse to the client
func (v *ValidationErrorsResponse) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {
	rw.WriteHeader(http.StatusBadRequest)
	errNewEncoder := json.NewEncoder(rw).Encode(v)
	if errNewEncoder != nil {
		log.Panic("Unable to encode and write response")
	}
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
	errNewEncoder := json.NewEncoder(rw).Encode(clientMessage{o.Err.Error()})
	if errNewEncoder != nil {
		log.Panic("Unable to encode and write response")
	}
}

// ResponseForError logs an error and returns the expected error type
func ResponseForError(logger Logger, err error) middleware.Responder {
	// AddCallerSkip(1) prevents log statements from listing this file and func as the caller
	skipLogger := logger.WithOptions(zap.AddCallerSkip(1))

	// Some code might pass an uninstantiated error for which we should throw a 500
	// instead of throwing a nil pointer dereference.
	if err == nil {
		skipLogger.Error("unexpected error")
		return newErrResponse(http.StatusInternalServerError, errors.New(NilErrMessage))
	}

	cause := errors.Cause(err)
	switch e := cause.(type) {
	case route.Error:
		skipLogger.Info("Encountered error using route planner", zap.Error(e))
		// Handle RouteError codes
		switch e.Code() {
		case route.UnsupportedPostalCode, route.UnroutableRoute:
			return newErrResponse(http.StatusUnprocessableEntity, err)
		case route.ShortHaulError:
			return newErrResponse(http.StatusConflict, err)
		default:
			return newErrResponse(http.StatusInternalServerError, err)
		}
	case *pq.Error:
		skipLogger.Info(SQLErrMessage, zap.Error(e))
		return newErrResponse(http.StatusInternalServerError, errors.New(SQLErrMessage))
	default:
		return responseForBaseError(skipLogger, err)
	}
}

func responseForBaseError(logger Logger, err error) middleware.Responder {
	skipLogger := logger.WithOptions(zap.AddCallerSkip(1))

	switch errors.Cause(err) {
	case models.ErrFetchNotFound:
		skipLogger.Info("not found", zap.Error(err))
		return newErrResponse(http.StatusNotFound, err)
	case models.ErrFetchForbidden:
		skipLogger.Info("forbidden", zap.Error(err))
		return newErrResponse(http.StatusForbidden, err)
	case models.ErrWriteForbidden:
		skipLogger.Info("forbidden", zap.Error(err))
		return newErrResponse(http.StatusForbidden, err)
	case models.ErrWriteConflict:
		skipLogger.Info("conflict", zap.Error(err))
		return newErrResponse(http.StatusConflict, err)
	case models.ErrUserUnauthorized:
		skipLogger.Info("unauthorized", zap.Error(err))
		return newErrResponse(http.StatusUnauthorized, err)
	case uploaderpkg.ErrZeroLengthFile:
		skipLogger.Info("uploaded zero length file", zap.Error(err))
		return newErrResponse(http.StatusBadRequest, err)
	case models.ErrInvalidPatchGate:
		skipLogger.Info("invalid patch gate", zap.Error(err))
		return newErrResponse(http.StatusBadRequest, err)
	case models.ErrInvalidTransition:
		skipLogger.Info("invalid transition", zap.Error(err))
		return newErrResponse(http.StatusBadRequest, err)
	case models.ErrDestroyForbidden:
		skipLogger.Info("invalid deletion", zap.Error(err))
		return newErrResponse(http.StatusBadRequest, err)
	default:
		skipLogger.Error("unexpected error", zap.Error(err))
		return newErrResponse(http.StatusInternalServerError, err)
	}
}

// ResponseForVErrors checks for validation errors
func ResponseForVErrors(logger Logger, verrs *validate.Errors, err error) middleware.Responder {
	skipLogger := logger.WithOptions(zap.AddCallerSkip(1))
	if verrs.HasAny() {
		skipLogger.Error("Encountered validation error", zap.Any("Validation errors", verrs.String()))
		return NewValidationErrorsResponse(verrs)
	}
	return ResponseForError(skipLogger, err)
}

// ResponseForCustomErrors checks for custom errors and returns a custom response body message
func ResponseForCustomErrors(logger Logger, err error, httpStatus int) middleware.Responder {
	skipLogger := logger.WithOptions(zap.AddCallerSkip(1))
	skipLogger.Error("Encountered error", zap.Error(err))

	return newErrResponse(httpStatus, err)
}

// ResponseForConflictErrors checks for conflict errors
func ResponseForConflictErrors(logger Logger, err error) middleware.Responder {
	skipLogger := logger.WithOptions(zap.AddCallerSkip(1))
	skipLogger.Error("Encountered conflict error", zap.Error(err))

	return newErrResponse(http.StatusConflict, err)
}
