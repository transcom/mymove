package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/transcom/mymove/pkg/migrate"
)

func main() {

	root := cobra.Command{
		Use:                   "clean-migration INPUT_FILE|- [OUTPUT_FILE]",
		DisableFlagsInUseLine: true,
		Short:                 "Clean migration file",
		Long:                  "Clean migration file",
		RunE: func(cmd *cobra.Command, args []string) error {

			if len(args) == 0 {
				return errors.New("must provide input file or stdin as -")
			}

			errParse := cmd.ParseFlags(args)
			if errParse != nil {
				return errors.Wrap(errParse, "Could not parse flags")
			}
			flag := cmd.Flags()

			v := viper.New()
			errBind := v.BindPFlags(flag)
			if errBind != nil {
				return errors.Wrap(errBind, "Could not bind flags")
			}
			v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
			v.AutomaticEnv()

			var inputReader io.Reader
			if args[0] == "-" {
				inputReader = os.Stdin
			} else {
				f, err := os.Open(args[0])
				if err != nil {
					return errors.Wrap(err, fmt.Sprintf("error reading from file %q", args[0]))
				}
				inputReader = f
				defer f.Close()
			}

			out := make(chan string, 1000)

			go migrate.CleanMigraton(inputReader, out)

			for line := range out {
				_, err := fmt.Fprintln(os.Stdout, line)
				if err != nil {
					return err
				}
			}

			return nil
		},
	}

	if err := root.Execute(); err != nil {
		panic(err)
	}
}
