package handlers

import (
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
}

// errResponse creates errResponse with default headers values
func newErrResponse(code int) *errResponse {

	return &errResponse{code: code}
}

// WriteResponse to the client
func (o *errResponse) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(o.code)
}

func responseForError(logger *zap.Logger, err error) middleware.Responder {
	switch errors.Cause(err) {
	case models.ErrFetchNotFound:
		return newErrResponse(http.StatusNotFound)
	case models.ErrFetchForbidden:
		return newErrResponse(http.StatusForbidden)
	case models.ErrInvalidPatchGate:
		return newErrResponse(http.StatusBadRequest)
	case models.ErrInvalidTransition:
		return newErrResponse(http.StatusBadRequest)
	default:
		logger.Error("Unexpected db error", zap.Error(err))
		return newErrResponse(http.StatusInternalServerError)
	}
}

func responseForVErrors(logger *zap.Logger, verrs *validate.Errors, err error) middleware.Responder {
	if verrs.HasAny() {
		logger.Error("Encountered validaton error", zap.Any("Validation errors", verrs.String()))
		return newErrResponse(http.StatusBadRequest)
	}
	return responseForError(logger, err)
}
