## [Jira ticket](tbd) for this change

## Summary

Is there anything you would like reviewers to give additional scrutiny?

[this article](tbd) explains more about the approach used.

## Setup to Run Your Code

<details>
<summary>💻 You will need to use three separate terminals to test this locally.</summary>

##### Terminal 1

Start the Storybook locally.

```sh
make storybook
```

##### Terminal 2

Start the UI locally.

```sh
make client_run
```

##### Terminal 3

Start the Go server locally.

```sh
make server_run
```

</details>

### Additional steps

<!-- Fill out the next section as you see fit, these are just suggestions. -->

1. Access the
2. Login as a
3.

## Verification Steps for Author

These are to be checked by the author.

- [ ] Tested in the Experimental environment (for changes to containers, app startup, or connection to data stores)
- [ ] Request review from a member of a different team.
- [ ] Have the Jira acceptance criteria been met for this change?

## Verification Steps for Reviewers

These are to be checked by a reviewer.

<!-- Fill out the next sections as you see fit, these are just suggestions. -->

### Frontend

- [ ] User facing changes have been reviewed by design.
- [ ] There are no aXe warnings for UI.
- [ ] This works in [Supported Browsers and their phone views](https://github.com/transcom/mymove/tree/main/docs/adr/0016-Browser-Support.md) (Chrome, Firefox, IE, Edge).
- [ ] There are no new console errors in the browser devtools
- [ ] There are no new console errors in the test output
- [ ] If this PR adds a new component to Storybook, it ensures the component is fully responsive, OR if it is intentionally not, a wrapping div using the `officeApp` class or custom `min-width` styling is used to hide any states the would not be visible to the user.

### Backend

- [ ] Code follows the guidelines for [Logging](https://transcom.github.io/mymove-docs/docs/dev/contributing/backend/Backend-Programming-Guide#logging)
- [ ] The requirements listed in [Querying the Database Safely](https://transcom.github.io/mymove-docs/docs/dev/contributing/backend/Backend-Programming-Guide/#querying-the-database-safely) have been satisfied.

### Database

#### Any new migrations/schema changes:

- [ ] Follows our guidelines for [Zero-Downtime Deploys](https://transcom.github.io/mymove-docs/docs/dev/contributing/database/Database-Migrations#zero-downtime-migrations)
- [ ] Have been communicated to #g-database
- [ ] Secure migrations have been tested following the instructions in our [docs](https://transcom.github.io/mymove-docs/docs/dev/contributing/database/Database-Migrations#secure-migrations)

## Screenshots
