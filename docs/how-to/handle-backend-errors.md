# How to handle back-end errors

## What we want

When users of Milmove make a mistake, we need a way to guide them back to the
happy path. This means returning not only coherent response codes for our
front-end to consume, but more helpful detail around what went wrong. This will
become increasingly important as we need to support users that will be using
Milmove programmatically as opposed to using our own front-end.

## How we're currently doing it

We currently have two ways of generating error responses in our APIs:

1. Using the generated Swagger response methods that conform to our API contract
   as defined in our Swagger yaml file.
2. Using `handlers.ResponseForError` and `handlers.ResponseForVErrors`, convenience functions
   which allow us to pass in any errors that are returned to us from model or
   service code. They take on the responsibility of deciding which error code to
   return based on the error values we pass in.

The approach in #2 is ergonomic, but makes it easy for our error handling to
drift from the API spec since the convenience functions write their own response
headers.

## Defining error models in Swagger

We typically define error responses in our Swagger definition like this:

```yaml
400:
  description: invalid request
401:
  description: request requires user authentication
404:
  description: office not found
500:
  description: server error
```

This will tell `go-swagger` to generate error response methods for each
respective status code, but we can take it further in order to reach our goal of
giving more robust, detailed responses to API consumers.

Both Swagger [2.0](https://github.com/OAI/OpenAPI-Specification/blob/master/versions/2.0.md) and [3.0](https://swagger.io/specification/#responseObject) allow us to define error models like the following:

```yaml
definitions:
  ClientError:
    type: object
    properties:
      title:
        type: string
      detail:
        type: string
    required:
      - title
      - detail
  ValidationError:
    allOf:
      - $ref: '#/definitions/ClientError'
      - type: object
    properties:
      invalid_fields:
        type: object
        additionalProperties:
          type: string
    required:
      - invalid_fields
```

With this approach, we can use generated code to add any data we'd like to an
error response.

### [RFC #7807](https://tools.ietf.org/html/rfc7807): Problem details for HTTP APIs

This RFC proposes some interesting ways of standardizing a way of providing
better descriptions for API errors. We will use the concepts of an error
`title`, `detail`, and extension fields in the rest of this document.

### 422 vs. 400

We currently return a `400 Bad Request` for validation errors. [It is recommended](https://tools.ietf.org/html/rfc4918#section-11.2) that we use `422 Unprocessable Entitity` instead.

## Example setups

### Validation errors

```yaml
# ...
post:
  summary: create an office user
  description: creates and returns an office user record
  operationId: createOfficeUser
  tags:
    - office
  parameters:
    - in: body
      name: officeUser
      description: Office user information
      schema:
        $ref: '#/definitions/OfficeUserCreatePayload'
  responses:
    201:
      description: Successfully created Office User
      schema:
        $ref: '#/definitions/OfficeUser'
    422:
      description: validation error
      schema:
        $ref: '#/definitions/ValidationError' #=> the interesting part
    500:
      description: internal server error
```

```go
func (h CreateOfficeUserHandler) Handle(params officeuserop.CreateOfficeUserParams) middleware.Responder {
  // ...

  createdOfficeUser, verrs, err := h.OfficeUserCreator.CreateOfficeUser(&officeUser, transportationIDFilter)
  if verrs != nil {
    payload := &adminmessages.ValidationError{
      InvalidFields: handlers.NewValidationErrorsResponse(verrs).Errors,
    }

    payload.Title = handlers.FmtString(handlers.ValidationErrMessage)
    payload.Detail = handlers.FmtString("The information you provided is invalid.")

    return officeuserop.NewCreateOfficeUserUnprocessableEntity().WithPayload(payload)
  }

  if err != nil {
    return officeuserop.NewCreateOfficeUserInternalServerError()
  }

  returnPayload := payloadForOfficeUserModel(*createdOfficeUser)
  return officeuserop.NewCreateOfficeUserCreated().WithPayload(returnPayload)
}
```

### Move is not in a state to be approved

```yaml
/moves/{moveId}/submit:
  post:
    summary: Submits a move for approval
    description: Submits a move for approval by the office. The status of the move will be updated to SUBMITTED
    operationId: submitMoveForApproval
    tags:
      - moves
    parameters:
      - name: moveId
        in: path
        type: string
        format: uuid
        required: true
        description: UUID of the move
      - name: submitMoveForApprovalPayload
        in: body
        required: true
        schema:
          $ref: '#/definitions/SubmitMoveForApprovalPayload'
    responses:
      200:
        description: returns updated (submitted) move object
        schema:
          $ref: '#/definitions/MovePayload'
      400:
        description: invalid request
      401:
        description: must be authenticated to use this endpoint
      403:
        description: not authorized to approve this move
      409:
        description: the move is not in a state to be approved
        schema:
          $ref: '#/definitions/ClientError' #=> the interesting part
      500:
        description: server error
```

```go
func (h SubmitMoveHandler) Handle(params moveop.SubmitMoveForApprovalParams) middleware.Responder {
  // ...

  submitDate := time.Time(*params.SubmitMoveForApprovalPayload.PpmSubmitDate)
  err = move.Submit(submitDate)
  if err != nil {
    payload := &internalmessages.ClientError{
      Title:  handlers.FmtString("This move is not in a state to be approved"),
      Detail: handlers.FmtString("Make sure the move is in state x before
      attempting to approve..."),
    }

    return moveop.NewSubmitMoveForApprovalConflict().WithPayload(payload)
  }

  // ...
}
```
