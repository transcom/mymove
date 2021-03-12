package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/cli"
	"github.com/transcom/mymove/pkg/logging"
)

func main() {
	flag := pflag.CommandLine
	cli.InitLoggingFlags(flag)

	v := viper.New()
	bindErr := v.BindPFlags(flag)
	if bindErr != nil {
		log.Fatal("failed to bind flags", zap.Error(bindErr))
	}

	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()

	dbEnv := v.GetString(cli.DbEnvFlag)

	logger, err := logging.Config(logging.WithEnvironment(dbEnv), logging.WithLoggingLevel(v.GetString(cli.LoggingLevelFlag)))

	if err != nil {
		log.Fatalf("Failed to initialize Zap logging due to %v", err)
	}

	if len(os.Args) != 2 && len(os.Args) != 3 {
		fmt.Println("Usage: big-cat <path> [limit]")
		os.Exit(1)
	}
	files, err := filepath.Glob(os.Args[1])
	if err != nil {
		panic(err)
	}
	limit := -1
	if len(os.Args) == 3 {
		l, err := strconv.Atoi(os.Args[2])
		if err != nil {
			panic(err)
		}
		limit = l
	}
	count := 0
	for _, file := range files {
		f, err := os.Open(filepath.Clean(file))
		if err != nil {
			panic(err)
		}
		if _, err = io.Copy(os.Stdout, bufio.NewReader(f)); err != nil {
			panic(err)
		}

		defer func() {
			if closeErr := f.Close(); closeErr != nil {
				logger.Debug("Failed to close filepath", zap.Error(closeErr))
			}
		}()

		count++
		if limit >= 0 && count == limit {
			break
		}
	}
}
