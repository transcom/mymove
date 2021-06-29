# OpenTelemetry

Currently on MilMove, all of the OpenTelemetry configuration is set up to be
opt-in, meaning that you can use it if you'd like to, but it is not assumed to
be the default way of running in local development environments.

## Collector

The OpenTelemetry Collector is set up to run as a "sidecar" in the local
development environment, i.e. it runs as its own Docker container. This allows
the server and clients to talk to just one instance, and the configuration
(for where the telemetry information is sent to) is managed with the sidecar
itself.

## Backend - Honeycomb

For the initial development in the efforts for distributed tracing, the config
is set up for using Honeycomb, mostly stemming from ease-of-use, familiarity,
and a generous free tier.

To effectively send information to Honeycomb, you need to set an API key and
optionally a "dataset" identifier in your environment variables, respectively
`HNY_API_KEY` and `HNY_DATASET`. (If unset, there is a default dataset
identifier already configured.)

The two main ways to get an API key are to: a) create your own account on
Honeycomb; or b) get invited to an existing team on Honeycomb. MilMove "best
practices" for this are currently TBD.

## Running the Collector

The top level `Makefile` has commands for managing the Collector.

To start (or re-start after a config change) the collector:

```sh
$ make otel_collector
```

If you want to verify the collector is up, you can then hit the health check:

```sh
$ curl http://localhost:13133/
{"status":"Server available","upSince":"2021-05-24T16:47:48.0230928Z","uptime":"24m17.6960986s"}
```

If you want to watch the logs:

```sh
$ docker logs -f otelc
...
2021-05-24T16:47:48.022Z info otlpreceiver/otlp.go:149 Setting up a second GRPC listener on legacy endpoint 0.0.0.0:55680 {"kind": "receiver", "name": "otlp"}
2021-05-24T16:47:48.022Z info otlpreceiver/otlp.go:87 Starting GRPC server on endpoint 0.0.0.0:55680 {"kind": "receiver", "name": "otlp"}
2021-05-24T16:47:48.022Z info otlpreceiver/otlp.go:105 Starting HTTP server on endpoint 0.0.0.0:55681 {"kind": "receiver", "name": "otlp"}
2021-05-24T16:47:48.023Z info builder/receivers_builder.go:75 Receiver started. {"kind": "receiver", "name": "otlp"}
2021-05-24T16:47:48.023Z info healthcheck/handler.go:128 Health Check state change {"kind": "extension", "name": "health_check", "status": "ready"}
2021-05-24T16:47:48.023Z info service/application.go:201 Everything is ready. Begin running and processing data.
```

To shutdown/kill/cleanup the Collector:

```sh
$ make otel_collect_kill
```

## Emitting to the Collector

To have the server emit tracing information to the Collector, set the
environment variable `OTEL_CONFIG` to `LOCAL_COLLECTOR`.

```sh
$ OTEL_CONFIG=LOCAL_COLLECTOR make server_run
```
