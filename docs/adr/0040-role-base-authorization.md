# Add Role-Based Authorization

**NOTE:** This ADR updates and supersedes [ADR0024 Model Authorization and Handler Design](./0024-model-authorization-and-handler-design.md).
.

As the MilMove system takes on more types of users (including TOOs, TIOs, and
the GHC prime contractor), it is important to ensure that only the correct users
are able to access each API endpoint. For example, an office user should be able to
index a list of moves, but a service member should not; they should only be able
to view their own move. Currently, permissions are handled by code in the models package
that checks if the current user is able to access data at the time it is accessed in a
handler or service object. Having these checks baked into the model layer makes it difficult
to reuse functionality in different context and also hard to see what checks are being done in
 a handler or service object.

We want to move to a more structured role-based system, ideally one that allows us to declare
what roles can access an endpoints in the same place we define that end point: the swaggerfile.

There will need to be additional checks throughout the system to ensure that users are able to, for example, only see their own data. But having a role-based system in place will make it easier for that code to assume that a user has certain attributes that it can then use to make such determinations.

## Considered Alternatives

* **Implement Role-Based Access Control in Middleware**

  Swagger yaml files include a new ‘roles’ tag on each endpoint,
indicating the list of roles that can access the endpoint.
The `UserAuthMiddleware` is edited to compare this list with the roles the user has
(by checking the session).

  We expect that there will be database changes made to allow for there to be a many-to-many relationship between users and roles.

  There is a [prototype](https://github.com/transcom/mymove/pull/2824/files) of this approach.

* **Use a Third-Party RBAC Library**

  We move to using an existing third-party library, specifically
[casbin](https://github.com/casbin/casbin).

  Casbin requires a list of associations, which should include information of two types:

  1. Which users are assigned to each role.
  2. Which endpoints are accessible by each role.

* **Use our existing system to prevent unauthorized data access**

  The MilMove system currently has checks in place to prevent unauthorized data access.
For example, consider a Service Member accessing their PPM. The handler
currently fetches the PPM based on ID, and then checks whether the SM associated
with that PPM is indeed the one making the request. An error is returned if not.
Note that this process takes only one database call.

  We can continue using this pattern as we expand the MilMove system.

* **Create an Enforcer Object**

  This approach involves writing code to encapsulate the authorization logic and explicitly invoking it within handlers or service objects as needed. It could be seen as taking our existing system and just extracting the authorization checks into functions outside the models.

  The name for this object comes from what similar objects are called within casbin.

## Decision Outcome

* Chosen Alternative: **Implement Role-Based Access Control in Middleware**
* The main motivators of this decision are:
  * It allows us to easily audit which endpoints are accessible by which roles by examining the API descriptions in the swaggerfile.
  * Having the role check in a middleware makes it less likely that an authorization check is missed.
  * We have completed spikes that show that this approach is doable with a minimal amount of additional work.

## Pros and Cons of the Alternatives

### *Implement Role-Based Access Control in Middleware*

* `+` Role/Endpoint associations live in Swagger yaml file, which is easy to read/modify
* `+` Flexible to addition of new role types
* `-` Requires some rework of database to represent roles
* `-` Not able to express row-level access rules

### *Use a Third-Party RBAC Library*

* `+` 3rd party libraries are more tested and receive more frequent updates
* `+` Able to express row-level access rules
* `+` Flexible to addition of new role types
* `-` Casbin has poor documentation
* `-` Is optimized for defining rules without code changes (not a goal for our project)
* `-` Unclear how to produce raw list of endpoint/role and role/user pairings for Casbin
* `-` Casbin doesn't scale well, according to community members

### *Use our existing system to prevent unauthorized data access*

* `+` Requires no immediate changes to existing code or developer training
* `+` Able to express row-level access rules
* `-` Authorization in models is not visible from handlers
* `-` It isn't possible to check for a permission outside attempting to load data
* `-` Model code requires arguments for authorization that doesn't make sense in all contexts
* `-` Adding new roles requires considerable code and database schema changes

### *Create an Enforcer Object*

* `+` Allows for a place to put finer-grained access controls
* `-` Requires explicit use in every handler
* `-`
