package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"reflect"
	"strings"

	openapierrors "github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/lib/pq"

	"github.com/gobuffalo/validate/v3"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/trace"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/route"
	uploaderpkg "github.com/transcom/mymove/pkg/uploader"
)

// The following are strings to be used in the title field of errors sent to the client

// SQLErrMessage represents string value to represent generic sql error to avoid leaking implementation details
const SQLErrMessage string = "Unhandled data error encountered"

// NotFoundMessage indicates a resource was not found
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

// MethodNotAllowedErrMessage indicates that the particular http request method was not allowed
const MethodNotAllowedErrMessage string = "Method Not Allowed"

// InternalServerErrMessage indicates that there was an internal server error
const InternalServerErrMessage string = "Internal Server Error"

// InternalServerErrDetail provides a default detail string
const InternalServerErrDetail string = "An internal server error has occurred"

// NotImplementedErrMessage indicates an endpoint has not been implemented
const NotImplementedErrMessage string = "Not Implemented"

// NotImplementedErrDetail indicates an endpoint has not been implemented
const NotImplementedErrDetail string = "This feature is in development"

// UnsupportedMediaTypeErrMessage indicates the server does not accept the media type sent
const UnsupportedMediaTypeErrMessage string = "Unsupported Media Type"

// NotAcceptableErrMessage indicates the server does not accept the media type requested
const NotAcceptableErrMessage string = "Not Acceptable"

// UnauthorizedErrMessage indicates the caller is not authorized
const UnauthorizedErrMessage string = "Unauthorized"

// ForbiddenErrMessage indicates the caller is forbidden
const ForbiddenErrMessage string = "Forbidden"

// ServiceUnavailableErrMessage indicates the service is currently unavailable
const ServiceUnavailableErrMessage string = "Service Unavailable"

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
func ResponseForError(logger *zap.Logger, err error) middleware.Responder {
	// AddCallerSkip(1) prevents log statements from listing this file and func as the caller
	skipLogger := logger.WithOptions(zap.AddCallerSkip(1))

	// Some code might pass an uninstantiated error for which we should throw a 500
	// instead of throwing a nil pointer dereference.
	if err == nil {
		skipLogger.Error("unexpected error")
		return newErrResponse(http.StatusInternalServerError, errors.New(NilErrMessage))
	}

	cause := errors.Cause(err)

	// pop/v6 now returns wrapped errors, so we have to unwrap if possible
	type unwrapper interface {
		Unwrap() error
	}
	unwrappable, ok := cause.(unwrapper)
	if ok {
		cause = unwrappable.Unwrap()
	}
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

func responseForBaseError(logger *zap.Logger, err error) middleware.Responder {
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
func ResponseForVErrors(logger *zap.Logger, verrs *validate.Errors, err error) middleware.Responder {
	skipLogger := logger.WithOptions(zap.AddCallerSkip(1))
	if verrs != nil && verrs.HasAny() {
		skipLogger.Error("Encountered validation error", zap.Any("Validation errors", verrs.String()))
		return NewValidationErrorsResponse(verrs)
	}
	return ResponseForError(skipLogger, err)
}

// ResponseForCustomErrors checks for custom errors and returns a custom response body message
func ResponseForCustomErrors(logger *zap.Logger, err error, httpStatus int) middleware.Responder {
	skipLogger := logger.WithOptions(zap.AddCallerSkip(1))
	skipLogger.Error("Encountered error", zap.Error(err))

	return newErrResponse(httpStatus, err)
}

// ResponseForConflictErrors checks for conflict errors
func ResponseForConflictErrors(logger *zap.Logger, err error) middleware.Responder {
	skipLogger := logger.WithOptions(zap.AddCallerSkip(1))
	skipLogger.Error("Encountered conflict error", zap.Error(err))

	return newErrResponse(http.StatusConflict, err)
}

// ServeCustomError is called by a hook provided by openapi to handle errors
// generated by that package. We override it below so that we can populate the
// title, detail, instance into the payload to match Milmove's desired
// error format.
func ServeCustomError(rw http.ResponseWriter, r *http.Request, err error) {

	rw.Header().Set("Content-Type", "application/json")
	var traceID = trace.FromContext(r.Context()).String()

	switch e := err.(type) {
	case *openapierrors.CompositeError:
		er := flattenComposite(e)
		// strips composite errors to first element only
		if len(er.Errors) > 0 {
			ServeCustomError(rw, r, er.Errors[0])
		} else {
			// guard against empty CompositeError (invalid construct)
			ServeCustomError(rw, r, nil)
		}
	case *(openapierrors.MethodNotAllowedError):
		rw.Header().Add("Allow", strings.Join(err.(*openapierrors.MethodNotAllowedError).Allowed, ","))
		rw.WriteHeader(asHTTPCode(int(e.Code())))
		if r == nil || r.Method != http.MethodHead {
			_, _ = rw.Write(errorAsJSON(e, traceID))
		}
	case openapierrors.Error:
		value := reflect.ValueOf(e)
		if value.Kind() == reflect.Ptr && value.IsNil() {
			rw.WriteHeader(http.StatusInternalServerError)
			_, _ = rw.Write(errorAsJSON(openapierrors.New(http.StatusInternalServerError, "Unknown error"), traceID))
			return
		}
		rw.WriteHeader(asHTTPCode(int(e.Code())))
		if r == nil || r.Method != http.MethodHead {
			_, _ = rw.Write(errorAsJSON(e, traceID))
		}
	case nil:
		rw.WriteHeader(http.StatusInternalServerError)
		_, _ = rw.Write(errorAsJSON(openapierrors.New(http.StatusInternalServerError, "Unknown error"), traceID))
	default:
		rw.WriteHeader(http.StatusInternalServerError)
		if r == nil || r.Method != http.MethodHead {
			_, _ = rw.Write(errorAsJSON(openapierrors.New(http.StatusInternalServerError, err.Error()), traceID))
		}
	}
}

func errorAsJSON(err openapierrors.Error, traceID string) []byte {

	// Turn all known openapi error codes into messages for title
	var title string
	switch err.Code() {
	case http.StatusMethodNotAllowed:
		title = MethodNotAllowedErrMessage
	case http.StatusNotFound:
		title = NotFoundMessage
	case http.StatusUnprocessableEntity:
		title = ValidationErrMessage
	case http.StatusBadRequest:
		title = BadRequestErrMessage
	case http.StatusInternalServerError:
		title = InternalServerErrMessage
	case http.StatusNotImplemented:
		title = NotImplementedErrMessage
	case http.StatusUnsupportedMediaType:
		title = UnsupportedMediaTypeErrMessage
	case http.StatusNotAcceptable:
		title = NotAcceptableErrMessage
	case http.StatusUnauthorized:
		title = UnauthorizedErrMessage
	case http.StatusForbidden:
		title = ForbiddenErrMessage
	case http.StatusServiceUnavailable:
		title = ServiceUnavailableErrMessage
	default:
		if err.Code() >= 600 {
			// All openapi validation errors are coded as 600+ errors
			title = ValidationErrMessage
		} else {
			title = "Unknown API Error"
		}
	}

	b, _ := json.Marshal(struct {
		Title    string `json:"title"`
		Instance string `json:"instance"`
		Detail   string `json:"detail"`
	}{title, traceID, err.Error()})
	return b
}

func flattenComposite(errs *openapierrors.CompositeError) *openapierrors.CompositeError {
	var res []error
	for _, er := range errs.Errors {
		switch e := er.(type) {
		case *openapierrors.CompositeError:
			if len(e.Errors) > 0 {
				flat := flattenComposite(e)
				if len(flat.Errors) > 0 {
					res = append(res, flat.Errors...)
				}
			}
		default:
			if e != nil {
				res = append(res, e)
			}
		}
	}
	return openapierrors.CompositeValidationError(res...)
}

func asHTTPCode(input int) int {
	// DefaultHTTPCode is used when the error Code cannot be used as an HTTP code.
	var DefaultHTTPCode = http.StatusUnprocessableEntity
	if input >= 600 {
		return DefaultHTTPCode
	}
	return input
}
