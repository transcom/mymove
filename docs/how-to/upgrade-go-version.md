# How to Upgrade Go Version

Upgrading the Go version that we use happens in roughly these steps:

1. Update [trussworks/circleci-docker-primary](https://github.com/trussworks/circleci-docker-primary) to point at an updated Go binary
2. Upgrade local Go version (`brew upgrade go`)
3. Create a PR for the `transcom/mymove` repo with the updated git hash created in step 1
4. Create a PR for the `transcom/ppp-infra` repo with the updated git hash created in step 1
5. Notify everyone that we're updating Go around the time your PR lands

For more details read the following sections.

## Updating Our Docker Image

- Grab the download URL and SHA256 checksum for the latest 64-bit Linux Go release from `https://golang.org/dl/`
  - The file name should be something like `gox.xx.x.linux-amd64.tar.gz`
- Update the Dockerfile and README with the new URL and checksum
  - See [this PR](https://github.com/trussworks/circleci-docker-primary/pull/10/files) as an example
- Open a PR and ask someone from the Truss Infra Team (not necessarily the MilMove Infra Team) to approve it

## Upgrade Local Go Version

- `brew upgrade go`
  - If you've done some PATH sorcery to point to a specific Go version (as detailed [here](https://github.com/transcom/mymove#setup-prerequisites)), you'll have to update that as well
  - If you've previously installed a version-specific go formula (e.g., go@1.12), you will need to uninstall it with
    `brew uninstall go@1.12` and then `brew install go` to get the current go version.
- `go version` to check it worked

## Update `transcom/mymove` Repo

- After your Docker image PR lands, grab the git hash from [Docker](https://hub.docker.com/r/trussworks/circleci-docker-primary/tags) that corresponds with your merged code
- Update the following files with the updated Docker image git hash and Go version:
  - `.circleci/config.yml`
  - `Dockerfile.local`
  - `Dockerfile.tasks_local`
  - `Dockerfile.tools_local`
  - `docker-compose.circle.yml`
  - Update `go-version` in `.github/workflows/go-auto-approve.yml`
  - See [this PR](https://github.com/transcom/mymove/pull/3524/files) as an example
- For minor go version changes (but not patch changes) update `scripts/check-go-version`.
- Rerun the Go formatter on the codebase with `pre-commit run --all-files go-fmt`
- Run `make e2e_test_docker` to test that the `Dockerfile.*local` files work  with the new image.
- Commit the above changes and any reformatted code and make sure everything builds correctly on CircleCI

## Notify Folks

- It can be jarring when everything suddenly breaks after pulling from master, so it's a nice courtesy to notify folks in #dp3-engineering that the official Go version will be updated shortly and their local Go version should be upgraded as well
- If `go-fmt` has changed as well, then in-flight PRs will need to be formatted before they are merged, lest they break master
