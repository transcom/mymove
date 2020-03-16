# How to Run Loki Tests Against Storybook

MilMove uses [Loki](https://loki.js.org/) for testing stories in storybook to ensure they have not changed. You can easily run the tests locally for verification using this document.

## Running tests locally

### Prereqs

* You will need storybook running locally already. See [How to Run Storybook](docs/how-to/run-storybook.md) for details.
* You will need Docker for Mac running locally as well. You can install the latest stable version from [here](https://download.docker.com/mac/stable/Docker.dmg).

### Running Loki Tests

Once the local storybook server is started you can run the tests.

```sh
make storybook_tests
```

Sample output

```sh
❯ make storybook_tests
yarn run loki test
yarn run v1.19.0
$ /Users/john/projects/dod/mymove/node_modules/.bin/loki test
loki test v0.18.1
  ✔ Chrome (docker)
✨  Done in 12.02s.
```

### If there are expected failures

If you are working on a storybook story and have finished making changes to the story or the components used the above command will fail. You should review the results files stored in `.loki/current`, `.loki/reference`, and `.loki/difference` directories. If the changes are as expected run the following command to update the reference files.

```sh
make loki_approve_changes
```
