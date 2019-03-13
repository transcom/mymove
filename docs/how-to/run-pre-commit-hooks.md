# Run and troubleshoot pre-commit hook's

[Pre-commit](https://pre-commit.com/) is a powerful tool that automates validations, lint checks and adds to developer quality of life. The config file that determines the actions of pre-commit hook can be found [here](/path/.pre-commit-config.yaml)

Pre-commit can be run by simply running the following command in terminal:
`pre-commit`

If you would like to run an individual hook, for example if you want to only run *prettier*: `pre-commit run prettier`

## Current pre-commit hooks

| Hook  | Description | Notes |
| ------------- | ------------- |------------- |
| go-version  | Attempts to load go version and verify it  |
|  check-json  | Attempts to load all json files to verify syntax |
| check-merge-conflict  | Check for files that contain merge conflict strings |
| check-yaml  | Attempts to load all yaml files to verify syntax |
| detect-private-key  | Checks for the existence of private keys |
|  trailing-whitespace | Trims trailing whitespace |
| markdownlint  | Linting rules for markdown files | more information [here](http://github.com/igorshubovych/markdownlint-cli)
| shell-lint  |  |
|  prettier |  |
| eslint  |  |
| swagger  |  |
| markdown-toc  |  |
| go-imports  |  |
| go-lint |  |
| gosec |  |
| gen-docs |  |
| dep-version |  |
| dep-check |  |
