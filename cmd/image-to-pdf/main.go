package main

import (
	"fmt"
	"log"

	"github.com/namsral/flag"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/paperwork"
	"github.com/transcom/mymove/pkg/storage"
	"github.com/transcom/mymove/pkg/uploader"
)

type stringSlice []string

func (i *stringSlice) String() string {
	return ""
}

func (i *stringSlice) Set(value string) error {
	*i = append(*i, value)
	return nil
}

var inputFiles stringSlice

func main() {
	flag.Var(&inputFiles, "input", "Image to add to PDF")
	flag.Parse()

	logger, err := zap.NewDevelopment()

	if err != nil {
		log.Fatalf("Failed to initialize Zap logging due to %v", err)
	}
	storer := storage.NewMemory(storage.NewMemoryParams("", ""))
	userUploader, err := uploader.NewUserUploader(storer, uploader.MaxCustomerUserUploadFileSizeLimit)
	if err != nil {
		log.Fatalf("could not instantiate uploader due to %v", err)
	}
	generator, err := paperwork.NewGenerator(userUploader.Uploader())
	if err != nil {
		log.Fatal(err.Error())
	}

	if len(inputFiles) == 0 {
		log.Fatal("Must specify at least one input file")
	}

	appCtx := appcontext.NewAppContext(nil, logger, nil, nil)
	path, err := generator.MergeImagesToPDF(appCtx, inputFiles, "")
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Printf("File written to %s\n", path)
}
