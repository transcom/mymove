# NOTE

This is a copy from [opentelemetry-go-contrip](https://github.com/open-telemetry/opentelemetry-go-contrib/tree/main/instrumentation)

The README there says, in part

> The last place instrumentation should be hosted is here in this repository. Maintaining instrumentation here hampers the development of OpenTelemetry for Go and therefore should be avoided.

There have been multiple PRs to improve the net/http instrumentation
that have not been merged for months and so, in typical go fashion, we
have to re-implement it ourselves.

This directory is a copy from

[otelhttp](https://github.com/open-telemetry/opentelemetry-go-contrib/tree/main/instrumentation/net/http/otelhttp)

as of [v.0.26.1](https://github.com/open-telemetry/opentelemetry-go-contrib/commit/7876cd14dc5f09765205caa0fb420fafe23141aa)

The `*_example_test.go` files were removed because they did not follow
the go convention of a single package per directory (they were in the
`otelhttp_test` package).
