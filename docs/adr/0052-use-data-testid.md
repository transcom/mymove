# Use `data-testid` as an attribute for finding components in tests

*[Jira Story](https://dp3.atlassian.net/browse/MB-3386)*

In front end code you can use a data tag selector, such as `data-testid`, on your components to make finding them in tests easier. Right now we're using two different values for test selector: `data-testid` and `data-cy`. We should pick one, and then update them in the MilMove codebase so that usage is consistent.

## Considered Alternatives

* Do nothing
* Use only `data-testid`
* Use only `data-cy`

## Decision Outcome

* Chosen Alternative: Use only `data-testid`
* After discussing we all strongly favored consistency in the code.
* `data-testid` makes it clear that this is a test id, implying that it is used for testing; whereas `data-cy` implies cypress, but we don't use cypress for all our tests and the suffix isn't clear to someone who isn't familiar with cypress.

## Pros and Cons of the Alternatives

### Do Nothing

* `+` Easy
* `-` Unclear which tag to use for tests
* `-` Unclear what `data-cy` is if not familiar with cypress

### Use only `data-testid`

* `+` Consistent
* `+` Name clearly indicates purpose
* `-` Have to change all references

### Use only `data-cy`

* `+` Consistent
* `-` Name is unclear as to purpose
* `-` Have to change all references
