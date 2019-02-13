# How To Generate Mocks with Mockery

[Mockery](https://github.com/vektra/mockery) provides the ability to easily generate mocks for golang interfaces. It removes the boilerplate coding required to use mocks.

 To generate a mock for testing purposes, you must use the mockery command line interface tool to do so.

 `$GOPATH/bin/mockery -name <nameOfInterface> -dir $GOPATH/src/github.com/transcom/mymove/<directoryInterfaceIsLocatedIn>`

 After you run this command a mock, corresponding to the interface name, will be created in the `/mocks` directory.

 *In Golang, mocks can only be created on interfaces - not structs. So, it is important that for whichever mock you are trying to generate, it should correspond to the appropriate interface.*