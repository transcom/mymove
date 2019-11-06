# Use dep to manage go dependencies

**NOTE:** Golang has explicitly moved to `go mod` and this project has as well, making this ADR obsolete.

## Considered Alternatives

* glide
* dep
* virtualgo

## Decision Outcome

The official Golang package manager `dep` will be used to manage server dependencies. `dep ensure` will be used to install all dependencies, `Godep.toml` will only be used to add required packages and to pin versions when incompatibilities are found. In the normal course of things, most dependencies will not need to be added to `Godep.toml` and we will rely on ensure's ability to automatically detect dependencies from imports.

## Pros and Cons of the Alternatives

### Glide

* `+` It is fairly mature
* `+` It uses the go-standard ./vendor library so it works fine with all of the build and install tooling
* `+` We've used it before and it has worked fine
* `-` It is not the official dependency management tool from the language authors
* `-` It has no way of encoding dependencies on non-imported tools (like linters or our migrator)

### Dep

* `+` It is the official dependency management tool
* `+` `dep ensure` is nicely designed, automatically doing what you expect for all the code that you * import
* `+/-` It has some support for encoding dev dependencies that are not imported
  * This is still not ideal, it will fetch version-pinned sources for those dependencies but will not build/* install them independently
* `-` It is still fairly new and unfinished. The docs say that it is ready for production use but it has a long roadmap

### VirtualGo

* `+` It works with dep and adds more isolation to your development environment
* `+` It allows you to actually install version-locked tools that are not imported by your code
* `-` It is another tool that would have to be installed and managed by all developers
