# Use go-swagger To Route, Parse, And Validate API Endpoints

## Considered Alternatives

* Buffalo
* Roll Our Own
* Swagger Go generation

## Decision Outcome

We use go-swagger to generate code for API request handling. The API will be defined in `swagger.yaml` and go-swagger will provide routing, JSON parsing, and whatever validation is expressed in `swagger.yaml`. While go-swagger can generate our `main.go` file, we disable that and take care of hooking up the generated API server to the rest of our routing ourselves. This means that we run our own server, we are not using the generated server file. Swagger code generation is taken care of by `./bin/gen_server.sh` and is invoked by Make anytime the `swagger.yaml` file changes.

### Notes

* Optional parameters should be marked with `x-nullable: true` so that go-swagger will represent them with pointers in structs. That way the handler will know if they have been provided or not
* Generated files all live under the `./pkg/gen` path and should never be checked into git
* All definitions that represent bodies of requests or responses should be named "Payloads" and should be described in the `definitions` section of `swagger.yaml`. This causes them to be named correctly in the generated code. "Payloads" is also what go-swagger calls the body "parameter" in the parameters struct passed that is generated for an API
* Use correct HTTP status codes for responding to success and errors. This is what expected by downstream tooling and is required by `swagger.yaml`, which is unable to represent multiple different 2XX responses.
* The generated payload objects should be used *exclusively* in the handler package. Everywhere else in our app, we should be using our model structs to represent data. The generated swagger code only exists to give us typed and validated structs for representing requests and responses. We do *not* want to rely on these structs for actual business logic.

## Pros and Cons of the Alternatives

### go-swagger

* `+` Translates `swagger.yaml` into go structs, giving us a single source of truth for our API definition
* `+` Runs validations before calling handlers, giving the handler some guarantees about what data it is given
* `+` Handles parsing of requests
* `+` It is maintained by responsive developers
* `-` Is a complex dependency we will be very dependent on
* `-` Some of its handling of `swagger.yaml` appears to be fragile, I managed to make it error with some valid `swagger.yaml` files, but I was able to get it to do what I wanted by re-organizing

### Buffalo

* `+` It is fairly mature
* `+` It integrates request handling with database model saving
* `-` It is ignorant of swagger
* `-` It integrates request handling with database model saving

### Roll Our Own

* `+` It would be more flexible, we could define requests and responses however we like
* `+` We would not have to learn the ins and outs of go-swagger which is a fairly complex bit of code
* `-` We would have to do our own routing, parsing, and validating
* `-` Our API could easily diverge from what we documented in `swagger.yaml`

### Official Swagger Go Generation

* `-` It appeared to generate garbage code
* `-` It didn't appear to do any validation
* `-` According to go-swagger, it didn't really work
* `-` I could never really get it to work
