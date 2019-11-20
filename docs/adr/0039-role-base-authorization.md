# Add Role-Based Authorization

**NOTE:** This ADR updates and supersedes [ADR0024 Model Authorization and Handler Design](./0024-model-authorization-and-handler-design.md).

[Additional background can be found in this design document](https://docs.google.com/document/d/1-CZx-hqDr7VtGtn0pr-UJ-ZH-fHiiBrn5Q_WX0CDGJY/edit).

**User Story:**

As the MilMove system takes on more types of users (including TOOs, TIOs, and
the GHC prime contractor), it is important to ensure that only the correct users
are able to access each API endpoint. For example, an office user should be able to
index a list of moves, but a service member should not; they should only be able
to view their own move. Currently, permissions by role are handled individually in API
handlers. Since this information is scattered, it poses risk of being incorrect and
not easily identified. A structured Role-Based Authorization scheme will unify
this information. This will make it easier for engineers to understand and edit the
permissions.

More specific than filtering by role, there is also a desire to be intentional about
the way we scope database calls to avoid users having access to data they shouldn't.
For example, while a service member should have access to the Fetch Move
endpoint, they should only be able to access their own move, not another SM's.

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

* **Create an Enforcer Object to Scope Queries and Perform Fine-grained Access Control**

  This approach involves creating an object to encapsulate two operations:
  1. Determining if an operation too granular for a purely role-based system is allowed. for example, can a service member access a specific move?
  2. Scoping database access by adding `WHERE` clauses to queries. This allows access control to be delegated to the database for data-loading operations where filtering data in code would be inefficient.

  There is a [prototype](https://github.com/transcom/mymove/pull/new/auth-spike) of this approach.

## Decision Outcome

* Chosen Alternative: **Implement Role-Based Access Control in Middleware**
and **Create a Service Object to Scope Queries**
* The main motivators of this decision are:
  * It accomplishes our desired feature set
  * Much of the work to determine how to implement this has already been prototyped
  * We are assuming that we will only need finer-grained access control in some cases and that blanket role-based checks will be sufficient in enough cases to justify having authorization at two levels.

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

### *Create an Enforcer Object to Scope Queries and Perform Fine-grained Access Control*

* `+` Codifies process for scoping data in uniform way
* `+` Does not require fetching large amounts of data, to then be filtered
in code (this is slow and doesn't scale)
* `+` Similar handlers can use similar/the same enforcer, making it easier to
change multiple handlers in response to db migrations
* `+` Able to express row-level access rules
* `-` Requires explicit use in every handler
