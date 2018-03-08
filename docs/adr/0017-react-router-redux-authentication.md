# Client side route restriction based on authentication

**User Story:** [#155131945](https://www.pivotaltracker.com/story/show/155131945)

    We want the loggedIn state of the user and basic user information (e.g. email address) available in redux state. We also want to be able to restrict access to client routes to logged in users.

## Considered Alternatives

* PrivateRoute component that is used for routes with restricted access
* HOC to wrap components that should be restricted to logged in users
* [EnsureLoggedInContainer](https://medium.com/the-many/adding-login-and-authentication-sections-to-your-react-or-react-native-app-7767fd251bd1) to create a section of routes that

## Decision Outcome

* Chosen Alternative: **PrivateRoute component that is used for routes with restricted access**

* This is how most of the samples examined worked. In particular, [react-router-redux example code](https://github.com/ReactTraining/react-router/blob/master/packages/react-router-redux/examples/AuthExample.js) worked this way. However, it is important to note that this only works if the routes are contained by a react router Switch.

There was also work done to build ducks for user state. Here we are using `js-cookie` to retrieve the jwt token and `jwt-decode` to extract the user email. Since login and logout are both performed by a server redirect, there is only one action at present to load the user and token. This action is called when the main route "/" is hit. A fair amount of exploration was done to find a way to fire this action whenever a new route is hit, but every example involved react router api that was removed in 4.0. We will need to do future work to handle expiration.

## Pros and Cons of the Alternatives <!-- optional -->

### PrivateRoute component that is used for routes with restricted access

* `+` Straightforward implementation
* `-` Dependency on Switch is not obvious from code (and behavior without it is strange and alarming)

### HOC to wrap components that should be restricted to logged in users

* `+` HOCs are the react way
* `-` Felt more complicated than needed for this use case. (Though a similar mechanism would work well for [authorization](https://hackernoon.com/role-based-authorization-in-react-c70bb7641db4). There is a stub of code for that for the future in `src/shared/User/Authorization.jsx`)

### EnsureLoggedInContainer to create a section of routes that

* `+` seemed like a nice way to group restricted routes
* `-` Doesn't work in react router 4.0
