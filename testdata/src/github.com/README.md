This whole folder exists because when we run analysistest, it overwrites GOPATH to be the directory we point it at.
In this case, that's the testdata directory. It's a problem when it comes to imports in our test files if we want things
to look close to how they would in the real code (e.g. `*pop.Connection`) so we are faking it out using this directory
structure that `go` would normally expect to find imports in.
