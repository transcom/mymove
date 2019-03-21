# *Working With Mocks: Generation and Assertion*

Currently, there is no way to approach mocking the behavior of dependencies in the backend API. Mocking allows you to simulate the behavior of an interface, in Golang, in a controlled way. Mock objects are helpful within the realm of unit testing. Oftentimes while testing, the piece of functionality that is being tested calls other dependencies - external services or methods within it. In order to exert full control over the test, mocking these dependencies allows one to define how they should behave. In the MilMove project, there are several instances that we integrate with services that we do not control, for example, the login and GEX services. Having the ability to mock these dependencies and others like it enables better testing. It would be helpful to have some way to work with mocks, such that they can be integrated easily with minimal developer overhead.

## Decision Drivers

* Allows easy control and testing of mock methods by allowing full control over parameter and return values
* Allows easy generation of mock methods
* Minimal intrusion on integrating into codebase
* Ease of implementation
* High usage and support
* Quick learning curve

## Considered Options

* Use Mockery for mock generation and Testify for mock assertion testing
* Write mock implementation methods in development without Mockery and only use Testify for mock assertion testing
* Use a test-specific struct implementation

## Decision Outcome

Using Mockery for mock generation and Testify for mock assertion testing is the chosen solution. Both [Mockery](https://github.com/vektra/mockery) and [Testify/Mock](https://godoc.org/github.com/stretchr/testify) are both actively maintained and are relatively easy to install and integrate. Using these two tools also will remove much of the complexity around creating and working with mocks. Using a code generation tool means that we don't have to write those mocked methods, increasing development speed. Integrating Testify's Mock assertion library allows us to do further complex mock assertion testing, such as being able to assert how many times the method was called, with what parameters and parameter types it was called with, along with denoting the return value and type.

Resources:

* [Mockery Mock Generation Library](https://github.com/vektra/mockery)
* [Testify/Mock Assertion Library](https://godoc.org/github.com/stretchr/testify)
* [How To Generate Mocks with Mockery in MilMove Project](https://github.com/transcom/mymove/docs/how-to/generate-mocks-with-mockery.md)
* [Mocking Dependencies in Go](https://medium.com/agrea-technogies/mocking-dependencies-in-go-bb9739fef008)
* [Improving Your Go Tests and Mocks with Testify](https://tutorialedge.net/golang/improving-your-tests-with-testify-go/)

## Pros and Cons of the Alternatives

### *Use Mockery for mock generation and Testify for mock assertion testing*

* `+` Implementation mock methods generated for use
* `+` Exerts full control allowing complex testing and assertion of mocked methods
* `+` Mockery and Testify/Mock Libraries actively maintained
* `+` Allows us to move quicker in development with code generation tools
* `+` Relatively simple API for mock generation and assertion
* `+` Testify is already apart of our testing suite and present in our codebase
* `-` Steeper learning curve given that these are new libraries that the team has not yet worked with
* `-` Mocks must be generated
* `-` Limited to what we can do by the mocking API or we have to contribute to it.

### *Write mock implementation methods in development without Mockery and only use Testify for mock assertion testing*

* `+` Exerts full control allowing complex testing and assertion of mocked methods
* `+` Testify/Mock Library actively maintained
* `-` Does not scale well as number of structs increases
* `-` Mocks cannot be generated
* `-` Minor learning curve for working with new mock assertion testing library

### *Test Specific Struct Implementation*

* `+` Lightweight and succinct
* `+` Low learning curve
* `-` Slower development
* `-` Does not scale well
* `-` Mocks cannot be generated
* `-` Must write test struct implementation and mock methods for everything that needs to be tested with mocking
* `-` No complex mocking assertion testing
* `-` No easy to use library API for checking things like if mocked method was called, how many times it was called, or what parameter and return values and types are