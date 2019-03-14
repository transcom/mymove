# Run and troubleshoot pre-commit hooks

[Pre-commit](https://pre-commit.com/) is a powerful tool that automates validations, lint checks and adds to developer quality of life. The config file that determines the actions of pre-commit hook can be found [here](/path/.pre-commit-config.yaml).

Pre-commit can be run by simply running the following command in terminal:
`pre-commit` or `make pre_commit_tests` which is similar to how CircleCI runs it.

*If pre-commit command is not found or errors out, please make sure you have the [prerequisites](README.md#setup-prerequisites) installed.*

## Testing

If you would like to run an individual hook, for example if you want to only run *prettier*: `pre-commit run prettier -a`

## Current pre-commit hooks

| Hook  | Description |
| ------------- | ------------- |
| go-version  | Attempts to load go version and verify it.
| check-json  | Attempts to load all json files to verify syntax. For more see [here](http://github.com/pre-commit/pre-commit-hooks).
| check-merge-conflict  | Check for files that contain merge conflict strings. For more see [here](http://github.com/pre-commit/pre-commit-hooks).
| check-yaml  | Attempts to load all yaml files to verify syntax. For more see [here](http://github.com/pre-commit/pre-commit-hooks).
| detect-private-key  | Checks for the existence of private keys. For more see [here](http://github.com/pre-commit/pre-commit-hooks).
| trailing-whitespace | Trims trailing whitespace. For more see [here](http://github.com/pre-commit/pre-commit-hooks).
| markdownlint  | Linting rules for markdown files. For more see [here](http://github.com/igorshubovych/markdownlint-cli).
| shell-lint  | Linter for shell files including spell check. For more see [here](http://github.com/detailyang/pre-commit-shell).
| prettier | Attempts to run [prettier](https://prettier.io/) hook against the code.
| eslint  | Attempts to run linting rules against the code base.
| swagger  | Attempts to run swagger validator for api, internal, order and dps endpoints.
| markdown-toc  | Wrapper script to generate table of contents on Markdown files.
| go-imports  | Attempts to run command `goimports` which updates your Go import lines, adding missing ones and removing unreferenced ones. For more see [here](https://godoc.org/golang.org/x/tools/cmd/goimports).
| go-lint | Attempts to run a linter against the go source code.
| go-vet | Attempts to examines Go source code and reports suspicious constructs, such as Printf calls whose arguments do not align with the format string.
| gosec | Inspects source code for security problems by scanning the Go AST. For more see [here](https://github.com/securego/gosec).
| gen-docs |Attempts to generate table of contents for the [docs/README](docs/README.md) file in doc folder.
| dep-version | checks the dep version.
| dep-check | Runs the command `dep check` to ensure dependencies are in sync with the `Gopkg.toml` file.
