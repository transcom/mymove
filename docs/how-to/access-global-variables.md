# How to Access a Global Application Variable

## Overview

In this project, we access application variables (environment variables or other variables we set for the application) by adding them to the handler context.  Through the context we pass the variables to the functions that require them.

## Why we do it this way

Environment variables should only be accessed in the main `serve.go` file and turned into real variables for passing around at that point. Accessing environment vars in other parts of the code increases the scope of our problems if there is something wrong with the environment vars. Also it increases problems with security if people are using the `os` package directly to get them instead of using the `spf13/viper` package which reads both environment vars and command line flags.

We use [spf13/viper](https://github.com/spf13/viper) and [spf13/pflags](https://github.com/spf13/pflag) to access environment variables today. It replaces using the `os` package and the `flag` package because it does both. The pattern is the 12-factor-app pattern.

## Getting Environment Variables

We use command line flags to get the environment variables. The flags are set in the `cli` package.  Viper can take the flag and gets the value associated with that flag.  For example:

`dbEnv := v.GetString(cli.DbEnvFlag)` returns the database environment name

`loginGovSecretKey := v.GetString(cli.LoginGovSecretKeyFlag))` grabs the `LOGIN_GOV_SECRET_KEY` from the `.envrc`

## Setting up global variables in the Handler Context

To add an application variable to the handler context, we create essentially a getter and setter in the handler context.
(Ex. `SetUseSecureCookie` and `UseSecureCookie`)
Follow the pattern in [pkg/handlers/contexts.go](https://github.com/transcom/mymove/blob/master/pkg/handlers/contexts.go)

Then, in the [cmd/milmove/serve.go](https://github.com/transcom/mymove/blob/master/cmd/milmove/serve.go) file, in the function `serveFunction` set the value using the setter.
For example:

```go
dbEnv := v.GetString(cli.DbEnvFlag)
isDevOrTest := dbEnv == "development" || dbEnv == "test"
useSecureCookie := !isDevOrTest
handlerContext.SetUseSecureCookie(useSecureCookie)
```

In your handler, you should now be able to access the value through the handler context by calling the getter (ex. `h.HandlerContext.UseSecureCookie()`)
