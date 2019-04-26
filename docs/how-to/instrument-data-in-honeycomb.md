# How to Instrument Data in Honeycomb

[Honeycomb](https://honeycomb.io) is a hosted service used to debug requests flowing through the live AWS environments. The MyMove application is configured to send events to the HoneyComb API using the [beeline-go](https://docs.honeycomb.io/getting-data-in/beelines/go-beeline/) library. The library ties into the [HTTP handler](https://golang.org/pkg/net/http/#Handler) and sends structured events to the Honeycomb API. By default, beeline-go derives fields from the incoming HTTP request. Each HTTP request maps to 1 event in Honeycomb. Below is a sample of some of the fields that are captured.

| field name           | field value             |
|----------------------|-------------------------|
| request.path         | /internal/duty_stations |
| duration_ms          | 848.190559              |
| response.status_code | 200                     |

The standard HTTP fields are a good start, but Honeycomb is more useful when we add zap logs, traces and new fields specific to MyMove.

## Tracing an API Request

Honeycomb tracing is a powerful tool for breaking down an API request into individual segments/spans. Each new span will include a name (usually the calling function name), a duration specifying how much time was spent in the span, and a series of unique identifiers that act as breadcrumbs to drill down into a particular request. [Beeline-go](https://docs.honeycomb.io/getting-data-in/beelines/go-beeline/) has some simple methods for instrumenting traces into the MyMove codebase.

An example of instrumenting SubmitMoveHandler, would be to to call `beeline.StartSpan(ctx, reflect.TypeOf(h).Name())` in the beginning of the handler and immediately add a `defer span.Send()` to defer sending the span to honeycomb until after the handler completes.

```golang
func (h SubmitMoveHandler) Handle(params moveop.SubmitMoveForApprovalParams) middleware.Responder {
    ctx, span := beeline.StartSpan(params.HTTPRequest.Context(), reflect.TypeOf(h).Name())
    defer span.Send()
```

To instrument subordinate functions, add `ctx context.Context` as the first parameter and pass the function name as the span name.

```golang
func (s *SocialSecurityNumber) SetEncryptedHash(ctx context.Context, unencryptedSSN string) (*validate.Errors, error) {
    ctx, span := beeline.StartSpan(ctx, "SetEncryptedHash")
    defer span.Send()
```

Be sure to pass the derived context from `beeline.StartSpan(...)` rather than the original context when passing context deeper into the function stack.  Reuse the variable name `ctx` when possible rather than allocating a new variable.

Useful fields can be added to the span that would help with debugging. To do this you can use [span.AddField](https://github.com/honeycombio/beeline-go/blob/master/trace/trace.go#L173).

```golang
    err = move.Submit(time.Now())
    span.AddField("move-status", string(move.Status))
```

You can add fields that apply to the entire function stack using [span.AddTraceField](https://github.com/honeycombio/beeline-go/blob/master/trace/trace.go#L206).  For example, the application name should be traced down through all the function calls.

```golang
    span.AddTraceField("auth.application_name", session.ApplicationName)
```

## Adding Zap Logs

Logs from Zap can also be included as part of the traces sent to Honeycomb. For example zap errors can be tied to the corresponding API request that generated the error. The hnyzap library is a wrapper around [zap](https://github.com/uber-go/zap). It includes additional functions TraceDebug, TraceInfo, TraceWarn, TraceError, TraceFatal to send the logs to honeycomb. Each function requires a additional context value. An example of tracing an error that's in the SubmitMoveHandler.

```golang
    if err != nil {
        h.HoneyZapLogger().TraceError(ctx, "Failed to change move status to submit",
            zap.String("move_id", moveID.String()),
            zap.String("move_status", string(move.Status)))
```

Note: Honeycomb supports bool, numbers and strings, so passing more complex types like zap.Object will result in values being set to "unsupported field type".

A typical use case is to use this along with standard error response code.  For instance, instead of using
`ResponseForErrors` you can use `h.RespondAndTraceError` and provide a message and additional parameters to trace with.
Similarly, validation errors can be captured via `h.RespondAndTraceVErrors` instead of using the standard
`ResponseForVErrors`.

## What Fields Should be Added to Honeycomb

_Do_:

* Try to come up with good field names that are descriptive to the data being represented.
  * A good naming convention to follow is `{go pkg name}.{field name}`. Session ids in the auth package look like auth.office_user_id
* More fields are better. The Honeycomb dataset is only as powerful as the data we decide to send, so if you think it would be useful to query requests on a particular field name, add it sooner rather than later.
* Honeycomb does a particularly good job of filtering events with unique fields like UUIDs.

_Don't_:

* Any fields that could potentially contain Personally Identifiable Information (PII) are a no-go.
  * Examples of PII include Name, email, social security number, date of birth etc.
