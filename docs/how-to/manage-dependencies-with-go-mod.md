# How to Manage Dependencies With go mod

[Go modules](https://github.com/golang/go/wiki/Modules) is the built-in dependency system provided by Go. It supersedes [dep](https://golang.github.io/dep/), which we previously used to manage Go dependencies.

It's important to note that go mod uses a [different dependency resolution algorithm](https://github.com/golang/go/wiki/Modules#version-selection) than many other packaging tools. It will install _oldest_ indirect
dependency (called _minimal version selection_) that will satisfy all direct dependencies, whereas other package managers will tend to install the _newest_.
You can read more about the rationale behind this approach [in the original proposal](https://github.com/golang/proposal/blob/master/design/24301-versioned-go.md#update-timing--high-fidelity-builds).

For the most part, a developer interacts with `go mod` using `go get`. The other go tools are likewise aware of how to work with go modules.

## Update all go dependencies

```console
$ go get -u
```

## Update a specific dependency

```console
$ go get -u github.com/pkg/errors
```

## Update a specific dependency to a specific branch

The following updates `github.com/pkg/errors` to the latest version available on the `master` branch:

```console
$ go get -u github.com/pkg/errors@master
```
