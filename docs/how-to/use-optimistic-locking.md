# How To Use Optimistic Locking

*Note: you'll probably want to use this on `PUT` or `PATCH` endpoints only.*

## Leaning on the query builder

Let's say you're building out a new endpoint that needs optimistic locking to avoid people updating stale data. You have a handler and a service object, let's call them `WidgetHandler` and `WidgetUpdater` respectively. The `WidgetUpdater` is going to take care of the business logic. Let's take a look at the handler first:

```go
type PatchWidgetHandler struct {
  //...
  services.WidgetUpdater
}


func (h PatchWidgetHandler) Handle(params widgetops.PatchWidgetParams) middleware.Responder {
  //...

  eTag := params.IfMatch

  widget := h.UpdateWidget(someArg, eTag)

  //...
}
```

We're grabbing the E-tag from the `If-Match` header, and passing it along to the
service object. Meanwhile, in the service object:

```go
type widgetUpdater struct {
  //...
  builder QueryBuilder
}

func (w *widgetUpdater) UpdateWidget(someArg interface{}, eTag string) {
  var widget models.Widget
  //...

  verrs, err := w.builder.UpdateOne(&widget, &eTag)
}

type PreconditionFailedError struct {
  id  uuid.UUID
  Err error
}

func (e PreconditionFailedError) Error() string {
  return fmt.Sprintf("widget with id: '%s' could not be updated due to the record being stale", e.id.String())
}
```

You'll notice that `UpdateOne`, a function that takes a model struct as an
`interface{}` in order to find the record and update it, now takes an optional
`string` argument for an E-tag. If the supplied E-tag is stale, `UpdateOne` will
return a `query.StaleIdentifierError` that you can then use to return a `412
Precondition Failed` in the handler:

```go
//...
widget, err := h.UpdateWidget(someArg, eTag)
if err != nil {
  logger.Error("error: ", zap.Error(err))

  switch e := err.(type) {
  case widget.NotFoundError:
    return widgetops.NewPatchWidgetNotFound()
  case widget.PreconditionFailedError:
    return widgetops.NewPatchWidgetPreconditionFailed()
  default:
    return widgetops.NewPatchWidgetInternalServerError()
  }
}
//...
```
