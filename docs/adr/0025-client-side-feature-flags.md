# Client Side Feature Flags using Custom JavaScript

**User Story:** Story [#158741324](https://www.pivotaltracker.com/story/show/158741324)

As new features are built out in the application, we need a way to prevent users from
seeing (and attempting to use) partially implemented functionality. At the same time,
we want to be able to selectively test features in specific environments.

While we eventually may wish to have such functionality available to the Go portions
of the application, our immediate needs can be satisfied through purely client-side
solutions.

## Considered Alternatives

* Use a third-party service such as [Launch Darkly](https://launchdarkly.com/implementation/)
* Detect current environment using the hostname and toggle features based on this environment
* Use environment variables to toggle features on and off, with or without a library

## Decision Outcome

* Chosen Alternative: **Detect current environment using `NODE_ENV` and fallback to using
  the hostname when `NODE_ENV === 'production'`. Toggle features based on this environment.**
* Simplest solution to the current set of known requirements.
* Requires the least investment in new libraries or services.
* Allows us to try out a simple approach with minimal investment, leaving open the
  possibility of moving to a more involved solution when it is warranted.
* This does require a commit and deploy to change the value of a flag, but this
  is OK given our intention of using this feature to gate large features.
* We can use the `AppContext` provider/consumer already used in the app for app-level
  settings to make the flags available throughout the app while preserving the ability
  to manipulate them in tests.

## Implementation Details

We will determine the environment using:

1. `NODE_ENV`
2. The page's hostname if `NODE_ENV` is `production`.

Each environment (production, staging, experimental, development, test) will have its own mapping of flags to values.

An example of checking the value of the `hhg` flag within JSX code is:

```jsx
<AppContext.Consumer>
  {settings => (
    <p>HHG is {settings.flags.hhg ? 'enabled' : 'disabled'}.</p>
  )}
</AppContext.Consumer>
```

## Pros and Cons of the Alternatives

### Launch Darkly

* `+` Provides a user-friendly web interface for controlling flags
* `+` Changes to flags do not require a code change or deployment
* `+` Provides server-side and client-side functionality
* `-` Allows for multivariate flags that contain values beyond simple boolean values
* `-` Adds a dependency on a new external service
* `-` Involves using the official Launch Darkly APIs

### Use environment variables to toggle features on and off, with or without a library

* `+` We already use environment variables to configure many parts of the application
* `-` Production environment variables are not available when the docker images are built on CircleCI
* `-` Many libraries exist to handle this, however, none that would work with both
  Go and JavaScript code without additional work on our part.
