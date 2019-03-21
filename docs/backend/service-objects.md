# Backend Service Objects Development Guide

## Table of Contents

<!-- Table of Contents auto-generated with `bin/generate-md-toc.sh` -->

<!-- toc -->

* [When a Service Object Makes Sense](#when-a-service-object-makes-sense)
* [Creating Service Objects](#creating-service-objects)
  * [Folder Structure And Naming](#folder-structure-and-naming)
  * [Parameters and Return Values](#parameters-and-return-values)
  * [Naming And Defining Service Object Structs and Interfaces](#naming-and-defining-service-object-structs-and-interfaces)
  * [Naming and Defining Service Object Execution Method](#naming-and-defining-service-object-execution-method)
  * [Instantiating Service Objects](#instantiating-service-objects)
* [Testing Service Objects with Mocks](#testing-service-objects-with-mocks)

Regenerate with "bin/generate-md-toc.sh"

<!-- tocstop -->

* [When a Service Object Makes Sense](#when-a-service-object-makes-sense)
* [Creating Service Objects](#creating-service-objects)
  * [Folder Structure And Naming](#folder-structure-and-naming)
  * [Parameters and Return Values](#parameters-and-return-values)
  * [Naming And Defining Service Object Structs and Interfaces](#naming-and-defining-service-object-structs-and-interfaces)
  * [Naming and Defining Service Object Execution Method](#naming-and-defining-service-object-execution-method)
  * [Instantiating Service Objects](#instantiating-service-objects)
* [Testing Service Objects](#testing-service-objects)

## When a Service Object Makes Sense

When writing or refactoring a piece of business logic to adhere to the service object pattern, it is important that this business function truly is the responsibility of a service object. Overusing this pattern and not applying it when appropriate can lead to several problems. It is necessary that developers make sure they are using the service object layer pattern when appropriate.

When to use a service object?

* [ ] dedicated encapsulation of a single piece business logic
* [ ] could possibly be re-purposed
* [ ] does this focus beyond parsing a request and rendering data
* [ ] does this singular piece of business logic use many different dependencies and/or different models

If you answered no to more than two of these questions, then a service object may not be the appropriate design pattern to use in your use case.

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

### Parameters and Return Values

**Parameters**
Remember that service objects should be reusable. Try to abstract as much out of the logic specific parameters to achieve this.
Pass as many parameters as make sense. Use your best judgement. In the following example from the codebase, we are only passing in one parameter to the `CreateForm` execution method, a
`template` variable with the type `FormTemplate`. This is because the `FormTemplate` is more complex than most service objects and this use case works for use here.
 Some service objects will only require only one or two parameters and a struct is not appropriate.`FormTemplate` only holds relatively abstract parameters such that the service
object can be reused if needed. Regarding `CreateForm`, this service object can be reused to generate another PDF by passing
different valid parameters.

```go
// paperwork.go
package services

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

**Return Values**
Service objects should return as many return values as appropriate. In the case of a service object like the CreateForm, the first parameter is whatever the service object is responsible for creating. If the first parameter returns a created entity, the second parameter should be an error.
In the case of a simple entity fetch by ID, the first parameter could be model validation errors. There are some situations that require more complex returns. For those, use your best judgement.

*Remember all `errors` should be Wrapped by using `errors.Wrap` so that the underlying error is propagated properly*

```go
// create_form.go
package paperwork

func (c createForm) CreateForm(template services.FormTemplate) (afero.File, error) {
  // Populate form fields with data
  err := c.FormFiller.AppendPage(template.Buffer, template.FieldsLayout, template.Data)
  if err != nil {
    return nil, errors.Wrap(err, fmt.Sprintf("Failure writing %s data to form.", template.FormType.String()))
  }
  ...
}
```

### Naming And Defining Service Object Structs and Interfaces

1. Define a private struct with the same name as the service object file, making sure that it is a noun camel-cased.

The struct fields are the dependencies needed for the service. To implement an interface in Go, all we need to do is to implement all the methods in the interface. By using an interface here
we are able to easily do mock testing on this service object. Adding these struct fields as interfaces will allow you to do testing with mocks; they are not required.

21 Add an interface for the service, that captures the behavior of the service object.

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
and struct. The service object execution method should be a method of the service object struct,
a struct of parameters that the service object requires, and returning values, as appropriate.

```go
// create_form.go
package paperwork

import (
    "github.com/spf13/afero"
    "github.com/transcom/mymove/pkg/services"
)

type createForm struct {
  fileStorer FileStorer
  formFiller FormFiller
}

func (c createForm) CreateForm(template services.FormTemplate) (afero.File, error) {
  ...
}
```

### Instantiating Service Objects

1. Create a `NewServiceObjectStruct` method that is responsible for creating a new service object. This method should be used whenever a new service object struct is needed. One of the main benefits of using service objects is abstracting implementation and returning an interface, then only using the interface in our codebase elsewhere. This allows us to separate interface from implementation.

```go
// create_form.go
package paperwork

import (
  "github.com/spf13/afero"
  "github.com/transcom/mymove/pkg/services"
)

type createForm struct {
  fileStorer Storer
  formFiller Filler
}

func NewFormCreator(FileStorer Storer, FormFiller Filler) services.FormCreator {
  return &createForm{FileStorer: FileStorer, FormFiller: FormFiller}
}

```

1. Add the service object as a field for the Handler struct of the handler that the service object will be executed in.

```go
// shipments.go
package publicapi
// CreateGovBillOfLadingHandler creates a GBL PDF & uploads it as a document associated to a move doc, shipment and move
type CreateGovBillOfLadingHandler struct {
  handlers.HandlerContext
  createForm services.FormCreator
}
```

1. Instantiate the service object while passing it in as a field for the Handler struct in `NewAPIHandler` function call.

```go
// publicapi/api.go
package publicapi

func NewPublicAPIHandler(context handlers.HandlerContext) http.Handler {
  ...
  publicAPI.ShipmentsCreateGovBillOfLadingHandler = CreateGovBillOfLadingHandler{
    context,
    paperworkservice.NewCreateForm(context.FileStorer().TempFileSystem(),
    paperwork.NewFormFiller(),
  )}
  ...
  return publicAPI.Serve(nil)
}
```

## Testing Service Objects with Mocks

1. Make sure the mock generation tool is installed by running `make server_deps`.
1. Generate the mock for the interface you'd like to test. See the [how-to doc](how-to/generate-mocks-with-mockery.md#how-to-generate-mocks-with-mockery)

```go
// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import io "io"
import mock "github.com/stretchr/testify/mock"
import paperwork "github.com/transcom/mymove/pkg/paperwork"


// FormFiller is an autogenerated mock type for the FormFiller type
type FormFiller struct {
  mock.Mock
}

// AppendPage provides a mock function with given fields: _a0, _a1, _a2
func (_m *FormFiller) AppendPage(_a0 io.ReadSeeker, _a1 map[string]paperwork.FieldPos, _a2 interface{}) error {
  ret := _m.Called(_a0, _a1, _a2)

  var r0 error
  if rf, ok := ret.Get(0).(func(io.ReadSeeker, map[string]paperwork.FieldPos, interface{}) error); ok {
    r0 = rf(_a0, _a1, _a2)
  } else {
    r0 = ret.Error(0)
  }

  return r0
}

```

1. Properly mock all methods for interface, denoting the parameter types, along with the return value.
1. Check the proper assertions

```go
// create_form_test.go
package paperwork

import (
  "github.com/pkg/errors"
  "github.com/spf13/afero"
  "github.com/stretchr/testify/assert"
  "github.com/stretchr/testify/mock"
  "github.com/stretchr/testify/suite"
  "github.com/transcom/mymove/mocks"
  paperworkforms "github.com/transcom/mymove/pkg/paperwork"
  "github.com/transcom/mymove/pkg/services"
)

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

  suite.NotNil(suite.T(), err)
  suite.Nil(suite.T(), file)
  serviceErrMsg := errors.Cause(err)
  suite.Equal(suite.T(), "Error for FormFiller.AppendPage()", serviceErrMsg.Error(), "should be equal")
  suite.Equal(suite.T(), "Failure writing GBL data to form.: Error for FormFiller.AppendPage()", err.Error(), "should be equal")
  FormFiller.AssertExpectations(suite.T())
}
```

It is important to note that when using a mocked interface, the mock function call will be called, not the original. This helps
to minimize side affects and allows us as developers to focus on what we are truly testing.

*Use `MockedInterface.On()` to mock a method. See their [docs](https://godoc.org/github.com/stretchr/testify/mock#Call.On) for more information.*
*Use `MockedInterface.AssertExpectations` to validate expectations, such as parameter type and number of times the method was called.*

Click [here](TODO) to see the recorded conversation on service objects.