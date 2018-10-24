# Service Layer and Dependency injection

## Context and Problem Statement

Currently the web service is built as two layers, Web Handlers (`pkg/handlers`) which implement interfaces based on the
swagger definitions of the services provided by the server and Model Objects (`pkg/models`) which marshal object representations
of data in and out of the database.

We are currently coming across a number of issues which suggest that we have reached the limits of what such a naive,
two-layer design can easily support, viz:

* it is not clear where Authorization code should live, i.e. code which enforces that logged in users only see and can access the data pertinent to them. Currently this is in the models (see ADR 0024)  but that means that models cannot be used for tools applications with different authorization controls, e.g. bulk loaders or admin interfaces.
* there is no place for code which touched multiple models and but is used by more than one handler, e.g. enforcing coherent state for multiple object relating to the same move (aka 'state machines') or making sure invoices line items are consistent between the GBL and the invoice.
* there is little or no encapsulation in the layers, so details of pop (database ORM) usage are in the handlers and equally swagger details appear in the model code. This makes testing and refactoring painful.

These variously lead to discussion around [Business Logic](https://en.wikipedia.org/wiki/Business_logic) and
[Service Layers](https://en.wikipedia.org/wiki/Service_layer_pattern). [Jim](https://github.com/jim) drew the teams attention
 to the [Service Object](https://medium.com/selleo/essential-rubyonrails-patterns-part-1-service-objects-1af9f9573ca1)
 pattern from rails. Looking for a similar pattern for go, it was suggested that we simply implement the Service Object
 pattern describe in the medium article in go.

This, in turn, lead to a search for a [Dependency Injection](https://en.wikipedia.org/wiki/Dependency_injection) framework
for golang which could be used in place of the global state used in Rails.

This ADR explains the choice of DI framework [DIG](https://github.com/uber-go/dig) and details conventions for naming and
using objects in the new 3-layer design.

## Decision Drivers

* Maintained (new commits less than 6 months ago)
* Support environment variables, and command line flags.
* Supports integer, duration, and time variables.
* Supports complex config or JSON, e.g., `map[string]string`.
* Mark variables as required and implement sanity checks

## Considered Options

* Built-in flag package
* Viper & pflag
* github.com/namsral/flag
* github.com/jessevdk/go-flags

## Decision Outcome

Chosen option: "Viper & pflag".  This option has the most community support and will give us continued flexibility as the code base grows over time.

## Pros and Cons of the Options

### Built-in flag package

Go ships with a built-in [flag](https://godoc.org/flag) package that provides support for command line flags.

* Good, no additional dependencies.
* Good, maintained but shouldn't receive any improvements either.
* Good, supports bool, (u)int, (u)int64, (u)float64, time.Duration, and string.
* Bad, no support for JSON variables or complex config, e.g., `map[string]string`.
* Bad, cannot mark variables as required (only provide defaults)
* Bad, invalid flag values cause panic (making custom sanity checks impossible)

### Viper & pflag

[Viper](https://github.com/spf13/viper) and [pflag](https://github.com/spf13/pflag) are 2 packages that are used together to enable 12-factor applications in Go.  VIPER is used by some of the most widely used Go programs, including `Hugo`, `go-swagger`, and `jfrog-cli-go`.

* Good, viper and pflag each have over 50 contributors and are actively maintained.
* Good, viper and pflag are "owned" by a [Steve Francia](https://github.com/spf13/), a Google employee and the creator of Hugo.
* Good, supports aliases to enable non-breaking improvements.
* Good, supports bool, int, int64, float64, duration, string, map[string]string, []string, map[string][]string, and time.Time.
* Bad, no support for JSON variables.
* Good, can unmarshal flag values into structs.
* Bad, cannot mark flag as required  (but can do defaults).
* Good, doesn't panic on bad values and can retrieve errors from pflag if needed.
* Good, supports json, toml, yaml, properties, and hcl config file formats.

#### Examples

Bind to config defined through command line flags (via `pflag`) and environment variables.

```go

flag := pflag.CommandLine

v := viper.New()

v.BindPFlags(flag) // bind to command line flags

// viper by default binds to upper case and
// supports a custom environment key replacer,
// but let's just use a typical one that replaces - with _
v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
// AutomaticEnv turns on binding for all environment variables
v.AutomaticEnv()
```

Custom validation of config variable and error handling during program startup.

```go
type InvalidConfigPort struct {
  Name  string
  Value interface{}
  Start int
  End   int
}

func (c *InvalidConfigPort) Error() string {
  return "config variable " + c.Name + " has invalid value " + fmt.Sprintf("%#v", c.Value) + ", must be in range " + strconv.Itoa(c.Start) + " - " + strconv.Itoa(c.End)
}

...
func main() {
  ...
  if v.IsSet("http-port-tls-none") {
    if value := v.GetInt("http-port-tls-none"); value < 8000 || value > 8999 {
      return &InvalidConfigPort{Name: "http-port-tls-none", Value: value, Start: 8000, End: 8999}
    }
  }
  ...
}
```

### github.com/namsral/flag

[flag](github.com/namsral/flag) is a drop-in replacement for Go's flag package that adds support for environment variables.  Currently used by our webserver (`github.com/transcom/mymove/cmd/webserver`).

* Bad, not maintained (the last code update was December 28, 2016).
* Good, supports bool, (u)int, (u)int64, (u)float64, time.Duration, and string (drop in replacement for built-in flag package)
* Bad, no support for JSON variables.
* Good, supports environment variables.
* Bad, only supports `name=value` and `name value` config file formats
* Bad, cannot mark variables as required.
* Bad, invalid values cause panic (making custom sanity checks impossible)

### github.com/jessevdk/go-flags

[go-flags](https://github.com/jessevdk/go-flags) enhances the functionality of the builtin `flag` package with support for many useful features.  You define your config as a single struct using fields and [struct tags](https://medium.com/golangspec/tags-in-golang-3e5db0b8ef3e).  Currently used by [truss-aws-tools](https://github.com/trussworks/truss-aws-tools).

* Borderline, last updated on March 31, 2018.
* Good, supports a variety of integer, float, string, and maps, including `[]*string{}`
* Bad, no support for JSON variables.
* Good, supports environment variables
* Bad, must unmarshal values into a single config struct.  Creates some baseline structure and increases readability, but reduces flexibility for responsively handling multiple contexts.
* Good, can mark variables as required
* Bad, hard to make custom sanity checks, since the config is parsed into a struct all at once.
