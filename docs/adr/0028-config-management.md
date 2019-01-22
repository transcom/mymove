# Config Management

* Status: accepted
* Deciders: @dynamike, @pjdufour-truss
* Date: 2018-09-21

## Context and Problem Statement

Our webserver (`github.com/transcom/mymove/cmd/webserver`) currently supports config variables defined as flags or environment variables (via `github.com/namsral/flag`), but our use of config throughout the application is not managed in a cohesive way.  While the use of [direnv](https://direnv.net/) and `.envrc` provides some basic external validation of the environment variables, we have a need to bring config parsing and validation logic into the application itself.  Beyond whether a config variable is set or not, we need to validate the values of these variables.  For example:

* is a port in the range from `8000 - 8999`,
* is the hostname valid,
* is the certificate valid,
* is the storage backend `filesystem` or `s3`,
* does s3 bucket name point to an existing bucket, and
* more.

The use of a more robust config framework with standard patterns will enable the seamless integration of new options and application contexts as we add new features to `mymove`.  Better management of config will enable the following:

* turn features on and off,
* debug startup errors,
* local docker server,
* per-branch test environments, and
* more.

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
