# Personal Property Prototype

[![Build status](https://img.shields.io/circleci/project/github/transcom/ppp/master.svg)](https://circleci.com/gh/transcom/ppp/tree/master)

This repository contains the application source code for the Personal Property Prototype, a possible next generation version of the Defense Personal Property System (DPS). DPS is an online system managed by the U.S. [Department of Defense](https://www.defense.gov/) (DoD) [Transportation Command](http://www.ustranscom.mil/) (USTRANSCOM) and is used by service members and their families to manage household goods moves.

This prototype was built by a [Defense Digital Service](https://www.dds.mil/) team in support of USTRANSCOM's mission.

## Development

### Prerequisites

Run `bin/prereqs` and install everything it tells you to. Then run `make client deps` and `make server deps`.

### Setup: Server

`make server_run`: installs dependencies and builds both the client and the server, then runs the server.

For faster development, use `make server_run_dev`. This builds and runs the server but skips updating dependences and re-building the client. Those tasks can be accomplished as needed with `make server_deps` and `make client_build`

You can verify the server is working as follows:

`> curl http://localhost:8080/api/v1/issues --data "{ \"issue\": \"This is a test issue\"}"`

from which the response should be

`{"id":1}`

Dependencies are managed by glide. To add a new dependency:
`GOPATH=/path/to/ProtoWebapp glide get new/dependency`

### Setup: Client

`make server_run`
`make client_run_dev`

The above will start the server running and starts the webpack dev server, proxied to our running go server.

Dependencies are managed by yarn
