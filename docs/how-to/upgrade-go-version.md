# How to Upgrade Go Version

Upgrading the Go version that we use happens in roughly these steps:

1. Update [trussworks/circleci-docker-primary](https://github.com/trussworks/circleci-docker-primary) to point at an updated Go binary
2. Upgrade local Go version
3. Update this repo with the updated git hash created in step 1
4. Notify everyone that we're updating Go around the time your PR lands

## Updating Our Docker Image

- Grab the download URL and SHA256 checksum for the latest 64-bit Linux Go release from `https://golang.org/dl/`
  - The file name should be something like `gox.xx.x.linux-amd64.tar.gz`
- Update the Dockerfile and README with the new URL and checksum
  - See [this PR](https://github.com/trussworks/circleci-docker-primary/pull/10/files) as an example
- Open a PR and ask someone from the infra team to approve it

## Upgrade Local Go Version

- `brew upgrade go`
  - If you've done some PATH sorcery to point to a specific Go version (as detailed [here](https://github.com/transcom/mymove#setup-prerequisites)), you'll have to update that as well
- `go version` to check it worked

## Update MilMove Repo

- After your Docker image PR lands, grab the git hash from [Docker](https://hub.docker.com/r/trussworks/circleci-docker-primary/tags) that corresponds with your merged code
- Update `.circleci/config.yml`, `README.md`, and `scripts/check-go-version` with the updated Docker image git hash and Go version
  - See [this PR](https://github.com/transcom/mymove/pull/1383/files) as an example
- Rerun the Go formatter on the codebase with `pre-commit run go-fmt --all-files`
- Commit the above changes and any reformatted code and make sure everything builds correctly on CircleCI

## Notify Folks

- It can be jarring when everything suddenly breaks after pulling from master, so it's a nice courtesy to notify folks in #dp3-engineering that the official Go version will be updated shortly and their local Go version should be upgraded as well
- If `go-fmt` has changed as well, then in-flight PRs will need to be formatted before they are merged, lest they break master
