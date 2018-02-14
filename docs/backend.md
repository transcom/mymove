# Backend Programming Guide

## Table of Contents

<!-- toc -->

* [Go](#go)
  * [Style and Conventions](#style-and-conventions)
  * [Libraries](#libraries)
    * [Pop](#pop)
  * [Learning](#learning)
  * [Testing](#testing)
    * [General](#general)
    * [Models](#models)
    * [Miscellaneous Tips](#miscellaneous-tips)

_Regenerate with `bin/docs-toc.sh`_

<!-- tocstop -->

## Go

### Style and Conventions

Generally speaking, we will follow the recommendations laid out in [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments). By its own admission, this page:
> _...collects common comments made during reviews of Go code, so that a single detailed explanation can be referred to by shorthands. This is a laundry list of common mistakes, not a style guide._

Despite not being an official style guide, it covers a good amount of scope in a concise format, and should be able to keep our project code fairly consistent.

Beyond what is described above, the following contain additional insights into how to write better Go code.

* [What's in a name?](https://talks.golang.org/2014/names.slide#1) (how to name things in Go)
* [Go best practices, six years in](https://peter.bourgon.org/go-best-practices-2016/)
* [A theory of modern Go](https://peter.bourgon.org/blog/2017/06/09/theory-of-modern-go.html)

### Libraries

#### Pop

We use Pop as the ORM(-ish) to mediate access to the database. [The Unofficial Pop Book](https://andrew-sledge.gitbooks.io/the-unofficial-pop-book/content/) is a valuable companion to Pop’s [GitHub documentation](https://github.com/markbates/pop).

### Learning

If you are new to Go, you should work your way through all of these resources (in this order, ideally):

1. [A Tour of Go](https://tour.golang.org) (in-browser interactive language tutorial)
1. [How to Write Go Code](https://golang.org/doc/code.html) (info about the Go environment, testing, etc.)
1. [Effective Go](https://golang.org/doc/effective_go.html) (how to do things “the Go way”)
1. [Daily Dep documentation](https://golang.github.io/dep/docs/daily-dep.html) (common tasks you’ll encounter with our dependency manager)
1. [Exercism](http://exercism.io/languages/go/about) offers a series of exercises with gradually increasing complexity

Additional resources:

* [GoDoc](https://godoc.org/) (where you can read the docs for nearly any Go package)
* Check out the [Go wiki](https://github.com/golang/go/wiki/Learn)
* *Video*: [Advanced Testing with Go](https://www.youtube.com/watch?v=yszygk1cpEc). (great overview of useful techniques, useful for all Go programmers)
* *Book*: [The Go Programming Language](http://www.gopl.io/)
* *Article*: [Copying data from S3 to EBS 30x faster using Golang](https://medium.com/@venks.sa/copying-data-from-s3-to-ebs-30x-faster-using-go-e2cdb1093284)

### Testing

Knowing what deserves a test and what doesn’t can be tricky, especially early on when a project’s testing conventions haven’t been established. Use the following guidelines to determine if and how some code should be tested.

#### General

* Use table-driven tests where appropriate.
* Make judicious use of helper funcs so that the intent of a test is not lost in a sea of error checking and boilerplate.

#### Models

In general, focus on testing non-trivial behavior.

* Structs do not need to be tested as they have no behavior of their own.
* Struct methods warrant a unit test if they contain important behavior, e.g. validations.
* Avoid testing functionality of libraries, e.g. model saving and loading (which is provided by Pop)
* Try to leverage the type system to ensure that components are “hooked up correctly” instead of writing integration tests.

#### Miscellaneous Tips

* Use `golang` instead of `go` in Google searches.
* Try to use the standard lib as much as possible, especially when learning.
