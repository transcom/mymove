# Using Swagger to manage server route authentication

**User Story:** [155131958](https://www.pivotaltracker.com/story/show/155131958)

Swagger and OpenAPI allow for defining route security and authentication through the same YAML file we use to define the rest of our API. Doing so would enhance Swagger's usage benefit of having information about our API be in a single centralized place, and would also allow us to possibly rely on tooling to generate some amount of authentication-related code for us.

## Considered Alternatives

* *Stick with what we have*
* *Implement one of the security methods as defined in OpenAPI 2.0*
* *Upgrade to OpenAPI 3.0*

## Decision Outcome

* Chosen Alternative: *Stick with what we have*
* *After doing some research, trying to implement route security through OpenAPI version 2 or 3 came with too many consequences and the benefits weren't valuable enough to be worth it. OpenAPI 2.0 only supports a few authentication schemes: basic (user/pass), API key, and OAuth 2.0. OpenAPI 3.0 supports bearer auth and cookie auth, but is unsupported by go-swagger.*

## Pros and Cons of the Alternatives <!-- optional -->

### *Stick with what we have*

* `+` *Doesn't require shoehorning into some other auth model*
* `+` *Still lets us control route security, just not as defined in Swagger YAML files*
* `-` *Route security isn't centralized with other API descriptors*
* `-` *Need to do all auth middleware manually*

### *Implement one of the security methods as defined in OpenAPI 2.0*

* `+` *Get to define route security in our Swagger YAML file*
* `+` *Some amount of auth checking is taken care of by go-swagger*
* `-` *Doesn't directly support our current JWT auth system, would have to jury rig around `apiKey`*
* `-` *Much of the code around token creation/renewal would still be manual*

### *Upgrade to OpenAPI 3.0*

* `+` *Can define security in swagger without any fundamental changes to our current system*
* `------` *go-swagger, a library we rely on heavily, would be unusable*
