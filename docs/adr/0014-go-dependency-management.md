# Go Dependency Management

Our code depends on other code. Handlers need database connections, the auth handlers need to know our URL. There has to be a way for code to get access to its dependencies, we have chosen to go down an explicit route where code is initialized with all its dependencies passed into the initializer. This makes our dependency tree very explicit and makes testing individual components easier.

## Considered Alternatives

* Explicitly encode dependencies in initializers for code that needs it. Wire up the resulting object graph in main.go
* Create a package for every dependency and use a package function to get access to the global you need
* Use package globals to store dependencies for each package that needs them, set those in main.go

## Decision Outcome

* Chosen Alternative: **Explicitly encode dependencies in initializers for code that needs it**

For things like handlers, we create an initializer for them that takes their dependencies as arguments (db and logger, at this time). Then, when they run, they have access to their dependencies they were created with. The purpose of main.go is to wire up all of our dependencies. There you create the db connection and give it to the handlers when you create them. This is mirrored in package tests, when you test a handler, you configure it when you create it by passing in its dependencies.

In our first pass, the swagger handlers all require the same dependencies so we have created a single struct that encodes them and type alias the individual handlers to it. In the future if different handlers require different deps, they can be transparently turned into their own structs.

If there were a dependency we needed to mock out for tests, this strategy would allow for that quite nicely. Once could simply define an interface that the dependency is expected to follow and then pass either the real dependency or the mock one into the initializer.

## Pros and Cons of the Alternatives

### Explicitly encode dependencies in initializers for code that needs it

* `+` Instances are always created with all their dependencies explicitly handed to them.
* `+` It's very difficult to not wire up dependencies correctly
* `+` Different instances can be very easily configured differently.
* `+` Testing different components is made simpler by allowing you to configure them each at init time rather than relying on globals.
* `-` More verbose, for things like the logger that are true globals it's redundant to pass them in to each one.
* `-` Care needs to be taken in considering the thread safety of dependencies being given separately instances

### Create a package for every dependency and use a package function to get access to the global you need

* `+` This is basically how the zap logger works by default, so it's a common pattern
* `+` No need to configure some of our code at all, it will just work as is, calling the globals.
* `-` It is not clear when you initialize a function that it depends on one of these globals being initialized

### Use package globals to store dependencies for each package that needs them

* `+` Different packages can have different dependencies
* `-` Each package needs to declare its dependencies
* `-` Nothing prompts initialization of a package before it is used.
* `-` This is sort of the worst of both worlds.
