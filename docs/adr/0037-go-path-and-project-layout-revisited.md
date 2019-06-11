# Put mymove outside of standard GOPATH

**NOTE:** This ADR updates and supersedes [ADR0003 Go Path and Project Layout](./0003-go-path-and-project-layout.md).

MilMove adopted using the go module system as introduction in `go1.12` in [PR 1932](https://github.com/transcom/mymove/pull/1932).
This gave us the option of continuing to host our repository inside of `$GOPATH` or outside of `$GOPATH` and per
[ADR0003 Go Path and Project Layout](./0003-go-path-and-project-layout.md) the decision was to stay in `$GOPATH` and
set the environment variable `GO111MODULE=on` as a way of forcing go module behavior.

This choice had implications for other tools that the project uses like `pre-commit` hooks that use golang libraries.  As seen in
[PR 2236](https://github.com/transcom/mymove/pull/2236) the interaction made it impossible to develop the project and
thus people set `GO111MODULE=auto` and were asked to move their repository outside of `$GOPATH`. This also
brings the local development inline with our CI/CD pipeline choices as modified in [PR 2172](https://github.com/transcom/mymove/pull/2172).

## Considered Alternatives

* Maintain direction of ADR0003 and keep repository checkout inside `$GOPATH` with `GO111MODULE=on`
* Move repository outside of `$GOPATH` with `GO111MODULE=auto`

## Decision Outcome

* Move repository outside of `$GOPATH` and explicitly set `GO111MODULE=auto`.
* Falls inline with best practices for golang going forward while allowing compatibility with dependencies not ready for go modules
* Forces everyone to move their directories or be unable to develop.

## Pros and Cons of the Alternatives

### Maintain direction of ADR0003 and keep repository checkout inside `$GOPATH` with `GO111MODULE=on`

Inside of `GOPATH` it is necessary to set `GO111MODULE=on` to force go module support. One example where
this causes issues is that `pre-commit` is not installing its hooks into the `$GOPATH` but instead into `~/.cache/pre-commit`.
The effect is that golang modules installed by `pre-commit` think that they are inside of `$GOPATH` when they are not
and that causes bad interaction issues with the tool itself.

* `+` No changes needed to repository location
* `-` Installation of `pre-commit` hooks that use golang will fail if `GO111MODULE=on` is set in the environment
* `-` Golang dependencies need to be ready to use go modules or installation doesn't work
* `-` Not in line with CI/CD setup which does not uses `auto` mode

### Move repository outside of `$GOPATH` with `GO111MODULE=auto`

Outside of the `GOPATH` it is not necessary to set `GO111MODULE=auto` as it is the default. However, being explicit
is better than implicit and `.envrc` will override whatever developers set globally in their environments.

* `+` Installation of `pre-commit` hooks that use golang will not fail
* `+` Installation of dependencies that do not support go modules will not fail
* `+` Prepared for future roll out of `go1.13` and future changes
* `+` In line with CI/CD setup which does not uses `auto` mode
* `-` Must change repository location outside of `$GOPATH`

## Resources

* [Go Modules](https://github.com/golang/go/wiki/Modules)
* [Old behavior vs New behavior](https://github.com/golang/go/wiki/Modules#when-do-i-get-old-behavior-vs-new-module-based-behavior)
