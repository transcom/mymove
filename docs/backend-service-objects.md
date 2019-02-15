# Backend Service Objects Development Guide

## Table of Contents

<!-- Table of Contents auto-generated with `bin/generate-md-toc.sh` -->

<!-- toc -->

* [When a Service Object Makes Sense](#when-a-service-object-makes-sense)
* [Creating Service Objects](#creating-service-objects)
  * [Folder Structure And Naming](#folder-structure-and-naming)
  * [Naming And Defining Service Object Structs and Interfaces](#naming-and-defining-service-object-structs-and-interfaces)
  * [Naming and Defining Service Object Execution Method](#naming-and-defining-service-object-execution-method)
  * [Instantiating Service Objects](#instantiating-service-objects)
* [Testing Service Objects](#testing-service-objects)

Regenerate with "bin/generate-md-toc.sh"

<!-- tocstop -->

* [When a Service Object Makes Sense](#when-a-service-object-makes-sense)
* [Creating Service Objects](#creating-service-objects)
  * [Folder Structure And Naming](#folder-structure-and-naming)
  * [Naming And Defining Service Object Structs and Interfaces](#naming-and-defining-service-object-structs-and-interfaces)
  * [Naming and Defining Service Object Execution Method](#naming-and-defining-service-object-execution-method)
  * [Instantiating Service Objects](#instantiating-service-objects)
* [Testing Service Objects](#testing-service-objects)

## When a Service Object Makes Sense

When writing or refactoring a piece of business logic to adhere to the service object pattern, it is important that this business function truly is the responsibility of a service object. Overusing this pattern and not applying it when appropriate can lead to several problems. It is necessary that developers make sure they are using the service object layer pattern when appropriate.

When to use a service object?

* [ ] dedicated encapsulation of a single piece business logic
* [ ] focuses on one thing
* [ ] could possibly be reused
* [ ] could possibly be extended
* [ ] does this focus beyond parsing a request and rendering data
* [ ] does this singular piece of business logic use many different dependencies and/or different models

If you answered no to more than of these questions, then a service object may not be the appropriate design pattern to use in your use case.

## Creating Service Objects

Once you have analyzed and determined that a service object is appropriate the next step is to actually create it.

### Folder Structure And Naming

1. Find or create appropriate directory.

Find or create the appropriate directory, in `/services` where the service object will live. Oftentimes, this directory
will be related to the actual model entity that it is dealing with. If this is something that involves multiple models, or
does not necessarily easily map to a model entity name, then it might be best to create a new folder that has a relevant name.

 ```bash
 /mymove
   /pkg
     /services
       /paperwork
```

1. Create the appropriate file(s) for the service object file, service object test file, and service object directory file.

Create a file with a name that captures what the service object is responsible for. Choose this name carefully as it will also be
the name of the service object execution method.

```bash
/mymove
  /pkg
    /services
      /paperwork
        create_form.go
        create_form_test.go
      paperwork.go
```

### Naming And Defining Service Object Structs and Interfaces

1. Define a private struct with the same name as the service object file, making sure that it is a noun camel-cased.

The struct fields are the dependencies needed for the service. However, these dependencies must be defined as interfaces.
To implement an interface in Go, all we need to do is to implement all the methods in the interface. By using an interface here
we are able to easily do mock testing on this service object.

Instead of defining types of the fields on the struct, like this:

```go
// create_form.go
package paperwork

import (
  "github.com/spf13/afero"
  paperworkforms "github.com/transcom/mymove/pkg/paperwork"
  "io"
)

// DO NOT DEFINE STRUCT LIKE THIS
type createForm struct {
  File afero.File
  FormFiller *paperworkforms.FormFiller
}
```

 Make sure to define interfaces as the types value of the fields on the struct, like this:

 ```go
 // create_form.go
 package paperwork

import (
  "github.com/spf13/afero"
  paperworkforms "github.com/transcom/mymove/pkg/paperwork"
  "io"
)

 type FileStorer interface {
  Create(string) (afero.File, error)
 }

 // Filler is an interface for FormFiller
 type FormFiller interface {
  AppendPage(io.ReadSeeker, map[string]paperworkforms.FieldPos, interface{}) error
  Output(io.Writer) error
 }

// DEFINE STRUCT LIKE THIS
type createForm struct {
  fileStorer FileStorer
  formFiller FormFiller
}
```

As you must define interfaces as the field types, this also means you must actual define those interfaces.
All methods in an interface must be those that the service object is actually using. For example, above the `createForm`
struct is only using the `Create` method as a part of the `FileStorer` interface. That is because that is the
only method needed from that dependency.

In addition to these structs and interfaces, it is also expected to add an interface for the service, capturing the behavior
of the execution method.

```go
// paperwork.go
package services

import (
  "bytes"
  "github.com/spf13/afero"
  paperworkforms "github.com/transcom/mymove/pkg/paperwork"
)

// FormTemplate are the struct fields defined to call CreateForm service object
type FormTemplate struct {
  Buffer       *bytes.Reader
  FieldsLayout map[string]paperworkforms.FieldPos
  FormType
  FileName string
  Data     interface{}
}

// FormCreator is the service object interface for CreateForm
type FormCreator interface {
  CreateForm(template FormTemplate) (afero.File, error)
}

```

### Naming and Defining Service Object Execution Method

The service object execution method is responsible for kicking off the service object call. Ideally, the service object
should expose only one public function, with helper private functions, as needed and when it makes sense. Oftentimes,
smaller private functions are good to unit test smaller units of functionality. The service object execution method should be the same as the file name
and struct. The service object execution method should be a method of the service object struct, ideally taking only one parameter,
a struct of parameters that the service object requires, and returning two parameters - what the service object is responsible for creating, and an error.

```go
// create_form.go
package paperwork

import (
    "github.com/spf13/afero"
    "github.com/transcom/mymove/pkg/services"
)

type createForm struct {
  FileStorer Storer
  FormFiller Filler
}

func (c createForm) CreateForm(template services.FormTemplate) (afero.File, error) {
  ...
}
```

### Instantiating Service Objects

1. Create a `NewServiceObjectStruct` method that is responsible for creating a new service object. This method should be used whenever a new service object struct is needed.

```go
// create_form.go
package paperwork

import (
  "github.com/spf13/afero"
  "github.com/transcom/mymove/pkg/services"
)

type createForm struct {
  FileStorer Storer
  FormFiller Filler
}

func NewCreateForm(FileStorer Storer, FormFiller Filler) services.FormCreator {
  return &createForm{FileStorer: FileStorer, FormFiller: FormFiller}
}

```

1. Instantiate the service object and pass it in as a field value for the Handler struct in `NewAPIHandler` function call.

```go
// publicapi/api.go
package publicapi

func NewPublicAPIHandler(context handlers.HandlerContext) http.Handler {
  ...
  publicAPI.ShipmentsCreateGovBillOfLadingHandler = CreateGovBillOfLadingHandler{context, paperworkservice.NewCreateForm(context.FileStorer().TempFileSystem(), paperwork.NewFormFiller())}
  ...
  return publicAPI.Serve(nil)
}
```

## Testing Service Objects

1. Make sure the mock generation tool is installed by running `make server_deps`.
1. Generate the mock for the interface you'd like to test. See the [how-to doc](how-to/generate-mocks-with-mockery.md#how-to-generate-mocks-with-mockery)

```go
// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import afero "github.com/spf13/afero"
import mock "github.com/stretchr/testify/mock"

// FileStorer is an autogenerated mock type for the FileStorer type
type FileStorer struct {
  mock.Mock
}

// Create provides a mock function with given fields: _a0
func (_m *FileStorer) Create(_a0 string) (afero.File, error) {
  ret := _m.Called(_a0)

  var r0 afero.File
  if rf, ok := ret.Get(0).(func(string) afero.File); ok {
    r0 = rf(_a0)
  } else {
    if ret.Get(0) != nil {
      r0 = ret.Get(0).(afero.File)
    }
  }

  var r1 error
  if rf, ok := ret.Get(1).(func(string) error); ok {
    r1 = rf(_a0)
  } else {
    r1 = ret.Error(1)
  }

  return r0, r1
}

```

1. Properly mock all methods for interface, denoting the parameter types, along with the return value.
1. Check the proper assertions

```go
func (suite *CreateFormSuite) TestCreateFormServiceFormFillerAppendPageFailure() {
  FileStorer := &mocks.FileStorer{}
  FormFiller := &mocks.FormFiller{}

  gbl := suite.GenerateGBLFormValues()

  FormFiller.On("AppendPage",
    mock.AnythingOfType("*bytes.Reader"),
    mock.AnythingOfType("map[string]paperwork.FieldPos"),
    mock.AnythingOfType("models.GovBillOfLadingFormValues"),
  ).Return(errors.New("Error for FormFiller.AppendPage()")).Times(1)

  createForm := NewCreateForm(FileStorer, FormFiller)
  template, _ := MakeFormTemplate(gbl, "some-file-name", paperworkforms.Form1203Layout, services.GBL)
  file, err := createForm.CreateForm(template)

  assert.NotNil(suite.T(), err)
  assert.Nil(suite.T(), file)
  serviceErrMsg := errors.Cause(err)
  assert.Equal(suite.T(), "Error for FormFiller.AppendPage()", serviceErrMsg.Error(), "should be equal")
  assert.Equal(suite.T(), "Failure writing GBL data to form.: Error for FormFiller.AppendPage()", err.Error(), "should be equal")
  FormFiller.AssertExpectations(suite.T())
}
```

It is important to note that when using a mocked interface, the mock function call will be called, not the original. This helps
to minimize side affects and allows us as developers to focus on what we are truly testing.

*Use `MockedInterface.On()` to mock a method. See their [docs](https://godoc.org/github.com/stretchr/testify/mock#Call.On) for more information.*
*Use `MockedInterface.AssertExpectations` to validate expectations, such as parameter type and number of times the method was called.*

Click [here](TODO) to see the recorded conversation on service objects.