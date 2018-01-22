# Put mymove into the standard gopath, eliminte server and client directories

The server component of mymove is written in Go. Go is very particular about where source code can live to allow for a standard way of fetching and building dependencies. We need to decide how we want to build the server, and how to fit into the go ecosystem to do so.

## Considered Alternatives

* Putting the server source directly into the go path
* Using `make` to set the gopath as being the ./server directory
* getting rid of the ./server directory all together and making the server exist at the top level of mymove

## Decision Outcome

mymove should reside in your gopath in the directory go expects: `$GOPATH/src/github.com/transcom/mymove`. Developers may put their gopath wherever they wish, but if they put it somewhere other than the default of ~/go, they should set $GOPATH correctly in their shell's profile. (and either way should add $GOPATH/bin to their path)

Within the mymove directory, server and client code exist together at the top level. This means that all internal go imports will look like this one for api: `github.com/transcom/mymove/pkg/api`

## Pros and Cons of the Alternatives

### Putting the server source directly into the gopath

* `+` This is what is widely expected of go code and is what all of Go's tooling expects
* `+` `go install` will install into the standard $GOBIN which makes it easy to access by adding that to your path in you shell's profile
* `-` All developers *must* checkout the source code into a very specific directory, and using symlinks to get to that directory can cause some go tooling to fail
* `-` It is possible (though this is mitigated by proper use of dep) for your code to accidentally rely on code that has been installed via `go get` which is not pinned to any version and could be potentially missing from another developer's machine. (though that would be caught by CI)

### Using Make to set a custom GOPATH for all go-related commands

* `+` Developers can checkout mymove anywhere they like and have building work
* `+` All our dependencies are fully isolated from any other go code on the system
* `-` Every invocation of go build and go install leaves detritus in mymove subdirectories
* `-` All go code has to be put under an `src` directory in the custom GOPATH
* `-` Import paths would be shorter, but non-standard. Go expects you to include github in your import paths.

### Putting the client and the server code together at the top level of mymove

* `+` Tooling that expects to be able to find vendored dependencies at the top level of a project now can
* `+` Dependencies (and tools) are now clearly for the entire project, not just one piece of it. i.e. It should be easy to use node based tools on our go code, now
* `+` The Makefile doesn't require cd'ing into different directories to do it's work
* `+` You can use all our tools directly, outside of the Makefile, without GOPATH magic
* `-` A layer of hierarchy has been removed, so it is less clear what certain directories mean

### Segregating server and client code into their own directories

* `+` It is more clear how the client and the server have their code grouped together
* `-` Even if we'd moved node_modules and package.json up to the top level to remedy the tooling problem, we'd have to have a second package.json in the client directory and we'd end up with duplicated commands between the two
* `-` If we didn't pull their depdency management up to the top level, it would limit our tool choices somewhat
