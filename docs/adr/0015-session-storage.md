# Session storage/handling

**User Story:** [155140012](https://www.pivotaltracker.com/story/show/155140012)

We want our users to be able to log in to our site and have that authentication backed by <http://www.login.gov>. We also want that authentication to spawn a session that persists until the user logs out or until a designated amount of time elapses.

## Considered Alternatives

* Encode our session data in a JSON web token that's stored on the client
* Encode session data in a secure cookie using gorilla/sessions
* Store session data server-side in Redis or an equivalent cache store

## Decision Outcome

* Chosen Alternative: **Encode our session data in a JSON web token that's stored on the client**

* Ultimately the real choice comes down to whether we want to store session information in a stateful or a stateless manner. A conventional stateful store backed by Redis or Memcached creates both a natural performance bottleneck and a single point of failure for session management. Alternatively, in a stateless system the server is only responsible for creating, reading, and verifying information that is securely stored on the client as a cookie. This makes stateless session storage, through JSON web tokens or some other means, an obvious choice.

* Storing sessions in encrypted cookies adds a bit of computational latency for encrypting and decrypting tokens, but this is easily outweighed by its advantages in contributing to a reliable, scalable architecture.

* We looked at a couple Go libraries that would allow us to easily implement client-side session storage: [jwt-go](https://github.com/dgrijalva/jwt-go) and [gorilla/sessions](https://github.com/gorilla/sessions) cookie stores. Gorilla adds some convenience around setting/retrieving client-side tokens, but the requisite data storage implementation (one-time-read flash messages) felt clumsy and wasn't necessary for our uses. jwt-go leaves us to set/retrieve cookies on our own, but the JWT claims data structure and expiration handling, including combined parsing and claims validation, doesn't leave many loopholes for accidentally dealing with an expired session. Since cookie handling is relatively simple through the `http` library, we opted for the jwt-go library.

## Pros and Cons of the Alternatives <!-- optional -->

### Encode our session data in a JSON web token that's stored on the client

* `+` Stateless session management
* `+` Provides structure for token expiration
* `+` Data is self contained in token
* `-` Vulnerable to replay attacks

### Encode session data in a secure cookie using gorilla/sessions

* `+` Stateless session management
* `+` Handles cookie setting/retrieval for us
* `-` Flash data storage is awkward

### Store session data server-side in Redis or an equivalent cache store

* `+` Maintain control over session state
* `+` Less data being transmitted to/from client
* `-` Single point of failure
* `-` Requires maintenance to avoid ballooning
* `-` Not terribly scalable
