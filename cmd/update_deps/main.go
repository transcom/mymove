package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"time"

	"github.com/rogpeppe/go-internal/modfile"
	"github.com/rogpeppe/go-internal/semver"
)

const (
	// Some dependencies take longer than 10 seconds to load, this increases the timeout
	// to deal with network latency
	cmdTimeout time.Duration = 30 * time.Second
)

// Use a custom branch for the following dependencies
var customBranches = map[string]string{
	"github.com/trussworks/pdfcpu": "afero",
}

// This program exists so that we can work around go mod's MVS behavior.
// For each dependency listed in go.mod, this program will update that dependency
// to either the latest released patch (using go get -u=patch) or the latest commit
// on the master branch based on if we're currently using a tagged version or a
// commit.
//
// There is a special case for pdfcpu, where we need to pull in the latest commit on a
// non-master branch.
func main() {
	content, readErr := ioutil.ReadFile("go.mod")
	if readErr != nil {
		log.Fatal(readErr)
	}

	file, parseErr := modfile.Parse("go.mod", content, nil)
	if parseErr != nil {
		log.Fatal(parseErr)
	}

	for _, req := range file.Require {
		fmt.Printf("%s", req.Mod.Path)

		ctx, cancel := context.WithTimeout(context.Background(), cmdTimeout)
		args := updateArgs(req)

		out, cmdErr := exec.CommandContext(ctx, "go", args...).CombinedOutput()
		if cmdErr != nil {
			fmt.Println(" ×")
			if ctx.Err() == context.DeadlineExceeded {
				log.Fatalf("timed out trying trying to run %s %s", "go", args)
			} else {
				log.Fatalf("failed to update %s: ran %s %v, got %s", req.Mod.Path, "go", args, string(out))
			}
		}

		cancel()
		fmt.Println(" ✓")
	}

	if output, err := modTidy(); err != nil {
		log.Fatalf("failed to run go mod tidy: got %s, error was %s", output, err)
	}
}

func updateArgs(req *modfile.Require) []string {
	if semver.Prerelease(req.Mod.Version) == "" {
		// Use the latest patch release if we're already using a tagged release
		return []string{"get", "-u=patch", req.Mod.Path}
	}

	branch := "master"
	customBranch, ok := customBranches[req.Mod.Path]
	if ok {
		branch = customBranch
	}

	return []string{"get", req.Mod.Path + "@" + branch}
}

func modTidy() ([]byte, error) {
	return exec.Command("go", "mod", "tidy").CombinedOutput()
}
