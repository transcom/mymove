# Use a consistent file structure for cypress tests

In the MilMove code base there is a previous iteration of PPM code and related cypress tests. Those older cypress tests were in two files `cypress/integration/mymove/ppm.js` and `cypress/integration/mymove/ppmCloseout.js`. These tests, before they were commented out, were long and later portions were dependent on earlier ones which occasionally caused failures that were difficult to find.

In thinking through how to improve things during our work on the new iteration of PPM features, it seemed better to break the tests up into individual tests per page of the on-boarding flow instead of trying have one test that went through the whole flow. To organize them further they are all under a `ppm` folder since they are all tied to the PPM features directly.

An update here is that the current tests live in a `mymove` folder but since that name no longer applies to the customer app we created a new one `milmove` with the intention of moving other tests that still apply to the customer portion into the new directory.

## Considered Alternatives (bold denotes chosen)

- Keep every test for the new PPM flow in a `ppm.js` file
- Create new files in the same location
- Create a single new test file
- **Create a new folder with new test files**

## Decision Outcome

_Chosen Alternative:_ Create a new folder with new test files

The recommendation is to keep the test files small and easy to follow, testing specific parts of a feature or page. If there are many files for one large feature a directory can be created to group the files together.

Example of this is the new PPM customer flow cypress tests in `cypress/integration/milmove/ppms` directory.

Recommend short test files over long ones. Long ones have been harder to debug in the past and have sometimes setup dependent tests.

It was a choice at the time of the first PPM tests for the new flow to not migrate all the tests. It was not possible given the time constraints at the time. However, to make the distinction between the old PPM tests and the new PPM tests it seemed prudent to create the new `milmove` directory.

This approach increases the likely hood that our tests will have repeated sections. This can be mitigated by the creation of functions that can be reused instead of repeating the code. This ADR doesn't dictate an approach and leaves it to the best judgement of those working with the tests. General rule of thumb though is if you need to repeate a section of code for a third time it's a good time to consider pulling that code into a reusable function.

## Pros and Cons of the Alternatives

### Keep every test for the new PPM flow in a `ppm.js` file

- `+` No new files
- `+` Files with old tests are updated
- `-` Only one file for entire PPM on-boarding flow
- `-` Long tests or lots of tests in one file can be hard to understand, follow, and debug
- `-` Customer cypress tests in `mymove` directory instead of `milmove`

### Create new files in the same location

- `+` Maintains current pattern of tests
- `-` Old PPM tests in same space as new PPM tests will lead to more confusion
- `-` Customer cypress tests in `mymove` directory instead of `milmove`

### Create a single new test file

- `+` Maintains current pattern of tests
- `-` Long cypress tests files get confusing quick
- `-` Can be harder to debug
- `-` Customer cypress tests in `mymove` directory instead of `milmove`

### Create a new folder with new test files

- `+` New pattern for short cypress tests avoids long test files
- `+` Short cypress files avoids long hard to diagnose tests
- `+` Clear distinction of old and new PPM tests
- `+` Smaller tests make it easier to write cypress tests as the larger feature is being developed
- `-` Opportunity to rename `mymove` to `milmove` in the cypress test folder
- `-` If all tests are not moved right away will have competing locations for storing new customer cypress tests
