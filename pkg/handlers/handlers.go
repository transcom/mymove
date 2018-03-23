package handlers

import (
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/markbates/pop"
	"go.uber.org/zap"
)

// HandlerContext contains dependencies that are shared between all handlers.
// Each individual handler is declared as a type alias for HandlerContext so that the Handle() method
// can be declared on it. When wiring up a handler, you can create a HandlerContext and cast it to the type you want.
type HandlerContext struct {
	db     *pop.Connection
	logger *zap.Logger
}

// NewHandlerContext returns a new HandlerContext with its private fields set.
func NewHandlerContext(db *pop.Connection, logger *zap.Logger) HandlerContext {
	return HandlerContext{
		db:     db,
		logger: logger,
	}
}

type S3Puter interface {
	PutObject(input *s3.PutObjectInput) (*s3.PutObjectOutput, error)
}

type S3HandlerContext struct {
	HandlerContext
	s3 S3Puter
}

func NewS3HandlerContext(handlerContext HandlerContext, s3Client S3Puter) S3HandlerContext {
	return S3HandlerContext{
		HandlerContext: handlerContext,
		s3:             s3Client,
	}
}
