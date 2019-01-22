package main

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	awssession "github.com/aws/aws-sdk-go/aws/session"
	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"
	"github.com/namsral/flag"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/paperwork"
	"github.com/transcom/mymove/pkg/storage"
	"github.com/transcom/mymove/pkg/uploader"
)

func main() {
	config := flag.String("config-dir", "config", "The location of server config files")
	env := flag.String("env", "development", "The environment to run in, which configures the database.")
	storageBackend := flag.String("storage_backend", "local", "Storage backend to use, either filesystem or s3.")
	s3Bucket := flag.String("aws_s3_bucket_name", "", "S3 bucket used for file storage")
	s3Region := flag.String("aws_s3_region", "", "AWS region used for S3 file storage")
	s3KeyNamespace := flag.String("aws_s3_key_namespace", "", "Key prefix for all objects written to S3")
	moveID := flag.String("move", "", "The move ID to generate advance paperwork for")
	build := flag.String("build", "build", "the directory to serve static files from.")
	flag.Parse()

	// DB connection
	err := pop.AddLookupPaths(*config)
	if err != nil {
		log.Fatal(err)
	}
	db, err := pop.Connect(*env)
	if err != nil {
		log.Fatal(err)
	}

	logger, err := zap.NewDevelopment()

	if err != nil {
		log.Fatalf("Failed to initialize Zap logging due to %v", err)
	}

	if *moveID == "" {
		log.Fatal("Usage: paperwork -move <29cb984e-c70d-46f0-926d-cd89e07a6ec3>")
	}

	var storer storage.FileStorer
	if *storageBackend == "s3" {
		zap.L().Info("Using s3 storage backend")
		if len(*s3Bucket) == 0 {
			log.Fatalln(errors.New("must provide aws_s3_bucket_name parameter, exiting"))
		}
		if *s3Region == "" {
			log.Fatalln(errors.New("Must provide aws_s3_region parameter, exiting"))
		}
		if *s3KeyNamespace == "" {
			log.Fatalln(errors.New("Must provide aws_s3_key_namespace parameter, exiting"))
		}
		aws := awssession.Must(awssession.NewSession(&aws.Config{
			Region: s3Region,
		}))

		storer = storage.NewS3(*s3Bucket, *s3KeyNamespace, logger, aws)
	} else {
		zap.L().Info("Using filesystem storage backend")
		fsParams := storage.NewFilesystemParams("tmp", "storage", logger)
		storer = storage.NewFilesystem(fsParams)
	}
	uploader := uploader.NewUploader(db, logger, storer)
	generator, err := paperwork.NewGenerator(db, logger, uploader)
	if err != nil {
		log.Fatal(err)
	}

	id := uuid.Must(uuid.FromString(*moveID))
	outputPath, err := paperwork.GenerateAdvancePaperwork(generator, id, *build)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(outputPath)
}
