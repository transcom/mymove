package migrate

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gobuffalo/fizz"
	"github.com/gobuffalo/fizz/translators"
	"github.com/gobuffalo/pop/v5"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// Builder is a builder for pop migrations.
type Builder struct {
	*pop.Match
	Path string
}

// Compile compiles the provided configration into a migration.
func (b *Builder) Compile(s3Client *s3.S3, wait time.Duration, logger *zap.Logger) (*pop.Migration, error) {

	if b.Type != "sql" && b.Type != "fizz" {
		return nil, &ErrInvalidFormat{Value: b.Type}
	}

	if b.Direction != "up" {
		return nil, &ErrInvalidDirection{Value: b.Direction}
	}

	if b.DBType != "all" && !pop.DialectSupported(b.DBType) {
		return nil, fmt.Errorf("unsupported dialect %s", b.DBType)
	}

	m := &pop.Migration{
		Version:   b.Version,
		Name:      b.Name,
		DBType:    b.DBType,
		Direction: b.Direction,
		Type:      b.Type,
	}

	if strings.HasPrefix(b.Path, "file://") {
		m.Path = b.Path
		m.Runner = func(m pop.Migration, tx *pop.Connection) error {

			defer func() {
				if r := recover(); r != nil {
					fmt.Fprintln(os.Stderr, "recovered:", r)
				}
			}()

			f, err := os.Open(m.Path[len("file://"):])
			if err != nil {
				return err
			}

			switch m.Type {
			case "sql":
				// have to use the tx from the Runner callback and
				// pass the logger
				errExec := Exec(f, tx, wait, logger)
				if errExec != nil {
					fmt.Fprintln(os.Stderr, errors.Wrapf(errExec, "error executing %s", m.Path).Error())
					return errors.Wrapf(errExec, "error executing %s", m.Path)
				}
			case "fizz":
				str, err := fizz.AFile(f, translators.NewPostgres())
				if err != nil {
					return errors.Wrapf(err, "could not fizz the migration %s", m.Path)
				}
				errExec := tx.RawQuery(str).Exec()
				if errExec != nil {
					return errors.Wrapf(errExec, "error executing statement: %q", str)
				}
			default:
				return &ErrInvalidFormat{Value: m.Type}
			}
			return nil
		}
		return m, nil
	} else if strings.HasPrefix(b.Path, "s3://") {
		m.Path = b.Path
		m.Runner = func(m pop.Migration, tx *pop.Connection) error {

			bucket := strings.SplitN(m.Path[len("s3://"):], "/", 2)[0]

			key := strings.SplitN(m.Path[len("s3://"):], "/", 2)[1]

			result, errGetObject := s3Client.GetObject(&s3.GetObjectInput{
				Bucket: aws.String(bucket),
				Key:    aws.String(key),
			})
			if errGetObject != nil {
				if aerr, ok := errGetObject.(awserr.Error); ok {
					logger.Error("AWS Error Code", zap.String("code", aerr.Code()), zap.String("message", aerr.Message()), zap.Any("AWSErr", aerr.OrigErr()))
				}
				return errors.Wrap(errGetObject, fmt.Sprintf("error reading migration file at %q", m.Path))
			}

			switch m.Type {
			case "sql":
				errExec := Exec(result.Body, tx, wait, logger)
				if errExec != nil {
					return errors.Wrapf(errExec, "error executing %s", m.Path)
				}
			case "fizz":
				str, err := fizz.AFile(result.Body, translators.NewPostgres())
				if err != nil {
					return errors.Wrapf(err, "could not fizz the migration %s", m.Path)
				}
				errExec := tx.RawQuery(str).Exec()
				if errExec != nil {
					return errors.Wrapf(errExec, "error executing statement: %q", str)
				}
			default:
				return &ErrInvalidFormat{Value: m.Type}
			}

			return nil
		}
		return m, nil
	}
	return nil, &ErrInvalidPath{Value: b.Path}
}
