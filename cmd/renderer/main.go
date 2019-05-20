package main

import (
	"errors"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"text/template"
)

func main() {

	templateFile := ""
	variablesFile := ""

	flag.StringVar(&templateFile, "t", "", "template file")
	flag.StringVar(&variablesFile, "v", "", "variables file")

	flag.Parse()

	// If no template file given, then error out
	if len(templateFile) == 0 {
		log.Fatal(errors.New("error: no template file given"))
	}

	// Read contents of template file into tmpl
	tmpl, err := ioutil.ReadFile(templateFile)
	if err != nil {
		log.Fatal(errors.New("error reading template file"))
	}

	ctx := map[string]string{}

	if len(variablesFile) > 0 {
		// Read contents of variables file into vars
		vars, readFileErr := ioutil.ReadFile(variablesFile)
		if readFileErr != nil {
			log.Fatal(errors.New("error reading variables file"))
		}

		// Adds variables from file into context
		for _, x := range strings.Split(string(vars), "\n") {
			// If a line is empty or starts with #, then skip.
			if len(x) > 0 && x[0] != '#' {
				// Split each line on the first equals sign into [name, value]
				pair := strings.SplitAfterN(x, "=", 2)
				ctx[pair[0][0:len(pair[0])-1]] = pair[1]
			}
		}
	}

	// Adds environment vairables to context
	// os.Environ() returns a copy of strings representing the environment, in the form "key=value".
	// https://golang.org/pkg/os/#Environ
	for _, x := range os.Environ() {
		// Split each environment variable on the first equals sign into [name, value]
		pair := strings.SplitAfterN(x, "=", 2)
		// Add to context
		ctx[pair[0][0:len(pair[0])-1]] = pair[1]
	}

	// Adds command line arguments to context
	for _, x := range flag.Args() {
		// Split each argument on the first equals sign into [name, value]
		pair := strings.SplitAfterN(x, "=", 2)
		// Add to context
		ctx[pair[0][0:len(pair[0])-1]] = pair[1]
	}

	t, err := template.New("main").Option("missingkey=error").Parse(string(tmpl))
	if err != nil {
		log.Fatal(err)
	}

	// If template uses variable that does not exist in context, then errors out.
	err = t.Execute(os.Stdout, ctx)
	if err != nil {
		log.Fatal(err)
	}
}
