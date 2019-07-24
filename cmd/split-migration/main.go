package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/transcom/mymove/pkg/migrate"
)

func main() {

	root := cobra.Command{
		Use:                   "split-migration INPUT_FILE|- [OUTPUT_FILE]",
		DisableFlagsInUseLine: true,
		Short:                 "Split migration file",
		Long:                  "Split migration file",
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

			lines := make(chan string, 1000)
			go func() {
				scanner := bufio.NewScanner(inputReader)
				for scanner.Scan() {
					lines <- scanner.Text()
				}
				close(lines)
			}()

			wait := 10 * time.Millisecond
			statements := make(chan string, 1000)
			go migrate.SplitStatements(lines, statements, wait)
			i := 0
			for stmt := range statements {
				fmt.Println("---------------------------------------------------")
				fmt.Println("Statement:", i)
				fmt.Println("---------------------------------------------------")
				fmt.Println(stmt)
				i++
			}

			return nil
		},
	}

	if err := root.Execute(); err != nil {
		panic(err)
	}
}
