# Add Role-Based Authorization

**NOTE:** This ADR updates and supersedes [ADR0024 Model Authorization and Handler Design](./0024-model-authorization-and-handler-design.md).

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

[Prototype](https://github.com/transcom/mymove/pull/2824/files)

Swagger yaml files include a new ‘roles’ tag on each endpoint,
indicating the list of roles that can access the endpoint.
The `UserAuthMiddleware` is edited to compare this list with the roles the user has
(by checking the session).

This yields an opportunity to adjust the way role associations
are handled by the database. Currently, each role
(office, service member, TOO, TIO, etc) has its own table in the db,
which has a foreign key into the users table. We move to a model that includes
a Roles table, which has a many-to-many relationship
with the users table. This makes it easy to add new roles, should the need arise.

`Users` table:

```text
Id
SM id (key to SM table)
TIO id (key to TIO table)
TOO id (key to TOO table)
Etc.
```

`UserRoles` table:

```text
Id
User ID (key to Users table)
Role ID (key to Roles table)
```

`Roles` table:

```text
Id
Role name (values will include “Service Member”, “TOO”, etc)
```

Along with this, we modify the session object to reflect the new
shape of the database. Specifically, it can include a field ‘Roles’
which lists the roles associated with the user.

* **Use a Third-Party RBAC Library**

We move to using an existing third-party library, specifically
[casbin](https://github.com/casbin/casbin).

Casbin requires a list of associations, which should include information of two types:

1. Which users are assigned to each role.
2. Which endpoints are accessible by each role.

TODO: fill this out.

* **Use our existing system to prevent unauthorized data access**

The MilMove system currently has checks in place to prevent unauthorized data access.
For example, consider a Service Member accessing their PPM. The handler
currently fetches the PPM based on ID, and then checks whether the SM associated
with that PPM is indeed the one making the request. An error is returned if not.
Note that this process takes only one database call.

We can continue using this pattern as we expand the MilMove system.

* **Create a Service Object to Scope Queries**

[Prototype](https://github.com/transcom/mymove/pull/new/auth-spike)

We can create a service object that will add Where clauses to queries.
Before a database query is made in a handler, this "Enforcer" will
add any relevant Where clauses to prevent unauthorized fetches of data.

```sql
SELECT * FROM moves
  WHERE moves.id = ?;
```

Might become:

```sql
SELECT * FROM moves
  LEFT JOIN service_members ON moves.service_members_id = service_members.id
  WHERE moves.id = ? AND service_members.id = ?;
```

Each handler will need either its own Enforcer object, or at least its own method, since
different queries will require different clauses. Note that this process still
only takes one database call.

## Decision Outcome

* Chosen Alternative: **Implement Role-Based Access Control in Middleware**
and **Create a Service Object to Scope Queries**
* The main motivators of this decision are:
  * It accomplishes our desired feature set
  * Much of the work to determine how to implement this has already been prototyped
  * It makes our codebase more maintainable

## Pros and Cons of the Alternatives

### *Implement Role-Based Access Control in Middleware*

* `+` Role/Endpoint associations live in Swagger yaml file, which is easy to read/modify
* `+` Flexible to addition of new role types
* `-` Requires some rework of database

### *Use a Third-Party RBAC Library*

* `+` 3rd party libraries are more tested and receive more frequent updates
* `-` Casbin has poor documentation
* `-` Unclear how to produce raw list of endpoint/role and role/user pairings for Casbin
* `-` Casbin doesn't scale well, according to community members

### *Use our existing system to prevent unauthorized data access*

* `+` Requires no up-front labor
* `+` Without rigid patterns, this is very flexible for different handlers' needs
* `-` Hard to manage, hard to spot errors before they arise
* `-` Difficult to edit many data requests quickly, as each may have different
implementation for data authorization

### *Create a Service Object to Scope Queries*

* `+` Codifies process for scoping data in uniform way
* `+` Does not require fetching large amounts of data, to then be filtered
in code (this is slow and doesn't scale)
* `+` Similar handlers can use similar/the same enforcer, making it easier to
change multiple handlers in response to db migrations
