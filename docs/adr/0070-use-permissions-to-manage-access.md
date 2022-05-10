# _Use permissions to manage access_

**User Story:** [MB-12237](https://dp3.atlassian.net/browse/MB-12237)

## Context

More access restrictions and related complexity are on the way with the QAE/CSR role soon coming into play. We need to be able to manage access to specific buttons, pages, components of the application in a way that will be easy to understand and manage. Currently, we have `users` and `roles` and we look at roles to decide what a user can do or see (manage access). See [ADR-0040](0040-role-base-authorization.md) for information on how we landed on this approach originally.

## Decision drivers

- Frontend and Backend API should be using the same set of 'rules' to manage access.

- If we decide to make changes/refactor, we need to be able to update existing access restrictions piece-wise over time. There is quite a bit of 'role conditional' access restriction already in place that we want to keep working as we transition to using permissions.

## Considered Alternatives

- _No Change - Continue using roles only to manage access_

- _Use `Permissions` in addition to `Roles` to manage acess_

- _Use 3rd package to manage access_

## Decision Outcome

- Chosen Alternative: Use `Permissions` in addition to `Roles` to manage access

- The concept of `Permissions` can be added to the codebase as an extension of the `Roles` work that we already have in place to provide a common set of access rules that both the frontend and backend can utilize. We can add permissions over time that join to our existing roles without disruption to existing code that makes use of the roles. This feels like the natural progression from where we are at with roles currently and doesn't come with significant drawbacks. Also, we gain the benefit of being able to frame user access in the code around what the user _has permission to do_ rather than _what roles they have_.

## Pros and Cons of the Alternatives

### _No Change - Continue using `Roles` only_

- `+` No new pattern to learn/implement

- `-` Not a best practice, commonly refered to as the 'naive approach' in documentation

- `-` Hard to maintain, especially as the number roles and access restrictions increases

- `-` No common mapping between roles and access restrictions making it very difficult to tell who has access to what without combing though the codebase.

- `-` Hard to keep frontend and backend API in sync with respect to access restrictions. Need to keep conditional role checks in sync.

### _Use `Permissions` in addition to `Roles`_

A permission is something a user can do (add shipment, edit allowances, flag for financial review, etc). Users have Roles, Roles grant Permissions.

- `+` Allows framing of user access in the code around what the user _has permission to do_ rather than _what roles they have_. This is a much more intuitive way to frame user access. It also makes it easier to understand what the user has access to and what they don't.

- `+` Significantly easier to maintain and refactor as complexity of access restrictions increases.

- `+` Helps limit developer scope when making changes to access restrictions as permissions are much more granular.

- `+` Can utilize the same permissions on both the frontend and backend so we use consistent access restrictions across the two.

- `-` We have a fair amount of 'role conditionals' in the codebase currently that we would want to refactor to match this pattern and use permissions rather than checking roles.

### _Use 3rd package to manage access_

I did not spend all that much time looking into the options here. We have already decided not to go this route in the past and I dont want to re-hash old discussions but I do think it is at least worth mentioning as an option.

- `+` 3rd party libraries are more tested and receive more frequent updates

- `-` Would require significant refactoring work, we have already decided not to go this route in the past when first adding roles (see [ADR-0040](0040-role-base-authorization.md)). Going back on this now would require refactoring all our role work up to this point.

- `-` Most packages seem to be geared towards either frontend or backend, we need both to be in sync using the same rules.

## Resources

- [ADR-0040 Role Based Authorization](0040-role-base-authorization.md)

- [What are permissions?](https://documentation.n-able.com/N-central/userguide/Content/User_Management/Role%20Based%20Permissions/role_based_permissions_what_are_permissions.htm)

- [Managing Access Control](https://levelup.gitconnected.com/access-control-in-a-react-ui-71f1df60f354)

- [How to conditionally render base on user permissions](https://medium.com/geekculture/how-to-conditionally-render-react-ui-based-on-user-permissions-7b9a1c73ffe2)

- [Clean patter for handling roles and permissions in large React apps](https://isamatov.com/react-permissions-and-roles/)
