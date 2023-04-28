## [Jira ticket](tbd)

## Summary

Is there anything you would like reviewers to give additional scrutiny?

[this article](tbd) explains more about the approach used.

## Setup to Run Your Code

- [Instructions for starting storybook](https://transcom.github.io/mymove-docs/docs/frontend/setup/run-storybook)
- [Instructions for starting the MilMove application](https://transcom.github.io/mymove-docs/docs/about/application-setup/milmove-local-client/)
- [Instructions for running tests](https://transcom.github.io/mymove-docs/docs/about/development/testing)

### How to test

1. Access the
2. Login as a
3.

## Verification Steps for Author

These are to be checked by the author.

- [ ] Tested in the Experimental environment (for changes to containers, app startup, or connection to data stores)
- [ ] Have the Jira acceptance criteria been met for this change?

## Verification Steps for Reviewers

These are to be checked by a reviewer.

### Frontend

- [ ] There are no aXe warnings for UI.
- [ ] This works in [Supported Browsers and their phone views](https://transcom.github.io/mymove-docs/docs/about/supported-browsers) (Chrome, Firefox, Edge).
- [ ] There are no new console errors in the browser devtools.
- [ ] There are no new console errors in the test output.
- [ ] If this PR adds a new component to Storybook, it ensures the component is fully responsive, OR if it is intentionally not, a wrapping div using the `officeApp` class or custom `min-width` styling is used to hide any states the would not be visible to the user.

### Backend

- [ ] Code follows the guidelines for [Logging](https://transcom.github.io/mymove-docs/docs/about/development/logging).
- [ ] The requirements listed in [Querying the Database Safely](https://transcom.github.io/mymove-docs/docs/backend/guides/golang-guide#querying-the-database-safely) have been satisfied.

### Database

#### Any new migrations/schema changes:

- [ ] Follows our guidelines for [Zero-Downtime Deploys](https://transcom.github.io/mymove-docs/docs/backend/setup/database-migrations#zero-downtime-migrations).
- [ ] Have been communicated to #g-database.
- [ ] Secure migrations have been tested following the instructions in our [docs](https://transcom.github.io/mymove-docs/docs/backend/setup/database-migrations#secure-migrations).

## Screenshots
