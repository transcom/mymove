package main

import (
	"io/ioutil"
	"os"

	"github.com/gobuffalo/fizz"
	"github.com/gobuffalo/fizz/translators"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func main() {

	root := cobra.Command{
		Use:                   "fizz-validate [path]...",
		Short:                 "Validate fizz files",
		Long:                  "Validate fizz files",
		DisableFlagsInUseLine: true,
		SilenceUsage:          true,
		SilenceErrors:         true,
		RunE: func(cmd *cobra.Command, args []string) error {
			trans := translators.NewPostgres()
			bubbler := fizz.NewBubblerWithDisabledExec(trans)
			for _, path := range args {
				file, err := os.Open(path)
				if err != nil {
					return errors.Wrapf(err, "error opening file at path %q", path)
				}
				b, err := ioutil.ReadAll(file)
				if err != nil {
					return errors.Wrapf(err, "error reading file at path %q", path)
				}
				_, err = bubbler.Bubble(string(b))
				if err != nil {
					return errors.Wrapf(err, "error translating file at path %q to sql", path)
				}
			}
			return nil
		},
	}

	if err := root.Execute(); err != nil {
		panic(err)
	}
}
