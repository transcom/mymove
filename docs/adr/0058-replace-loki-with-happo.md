# Use Happo for visual regression testing

## Background

Visual regression testing is a type of automated testing that is meant to catch unintended visual side effects in rendered UI. This is an area of risk due to the global nature of CSS -- when changing CSS code to modify one element, there is always a possibility of inadvertent changes to other elements without knowing anything happened until it is discovered by happenstance after the fact.

Visual regression testing occurs by first taking screenshots of selected portions of rendered UI to use as references, then comparing screenshots with code changes against those references, and finally alerting the team to differences between the images for review. Because any changes, whether they are intended or not, require explicit approval, building visual regression testing into our workflow can help ensure all UI changes have been reviewed by members of the design and product teams, as well as provide a visual of cross-browser implementations without needing manual QA.

## Problem

MilMove currently uses a tool called Loki for performing visual regression tests on all components rendered in Storybook. While Loki has been helpful with catching unintended visual changes, it also has several downsides:

- Loki is executed in local development & CI environments, so reference image assets must be checked into the code repo, and running/updating the tests locally is time- and resource-intensive.
- Because reference images are checked into the repo, the action of approving changes requires an additional commit, push, and CI build.
- Loki only runs tests in Chrome, so visual implementations in other browsers (Firefox, IE, Safari) are not tested.
- Review and approval of changes happens in local dev environments, so actual changes can be relatively opaque and may not be explicitly reviewed by designers.
- Loki only supports visual testing of Storybook components, so we would not be able to add visual tests to the application if we wanted.

## Considered Alternatives

- Continue using Loki (do nothing)
- Switch to Happo
- Switch to Chromatic
- Use SauceLabs visual testing
- Stop visual regression testing

## Decision Outcome

- **Chosen Alternative: Switch to Happo**
- `+` Tests run on Happo's hosted platform, so engineers and CircleCI don't need to spend time/resources running them
- `+` Visual tests can run against all of the browsers we support: Chrome, Firefox, IE11, Safari, iOS Safari, Edge (on all pricing plans, including free open source)
- `+` Happo provides a UI for users with access to view and approve or reject changes, and approving does not require committing changes
- `+` Happo has a plugin to use with Storybook components, but can also be used to screenshot the application itself (for example, as part of Cypress E2E tests)
- `~` Happo has an on-premise option available if security is a concern (but pricing is by request, and will require additional infra support to set up and maintain)
- `-` Happo runs tests on its hosted platform, so it requires giving Happo access to our Github repo (which is currently public)
- `-` Happo costs money for non-open source projects (pricing tiers are $125 / $250 / \$500 / month depending on usage)
- `-` Test reports for free open source projects are public to anyone who has the link

### Strategy

My proposed steps for migrating from Loki to Happo are:

1. Add the [Happo Github app](https://github.com/apps/happo) to this repo
1. Install the JS dependencies (`happo.io, happo-plugin-storybook`)
1. Add required configuration for Happo & Storybook
1. Test trigger a Happo run from local environment
1. Skip components that are currently skipped in Loki if needed
1. Replace Loki script with Happo in CI (and test)
1. Update relevant documentation around Storybook tests
1. Remove deprecated Loki scripts, files, documentation from the repo

Since Happo offers a 30 day free trial, the above _should be able to_ be completed independently of providing payment information for an account.

### References

- [Happo: Getting Started](https://docs.happo.io/docs/getting-started)
- [Happo: Storybook plugin](https://docs.happo.io/docs/storybook)
- [Cross-browser screenshot testing with Happo.io and Storybook](https://medium.com/happo-io/cross-browser-screenshot-testing-with-happo-io-and-storybook-bfb0b848a97a)

### Additional Questions

- **Do we have budget for this? Estimate # of runs, expected cost.**
  - Budget is TBD. If we determine no budget, we can investigate whether MilMove qualifies for Happo's free open source plan (this would require test reports be public to anyone with the link, though). Budget should be weighed against the current amount of time Loki is requiring from both engineers and CircleCI, which is not insubstantial. Running the tests in CircleCI is one of our longest-running jobs, and debugging or even approving intended changes requires additional engineering time.
- **Does it actually help us ensure that we are maintaining IE compatibility?**
  - Loki does no testing in IE, and Happo does, so: yes.
- **Are we losing anything that Loki is providing?**
  - No.
- **Estimated time cost to implement (very rough estimate fine, determined by reading Happo “getting started” docs, etc.)**
  - 1 day of engineering time to implement the strategy outlined above.
- **Does this feasibly integrate with our CI pipeline?**
  - [Yes](https://docs.happo.io/docs/continuous-integration#happo-ci-circleci)
- **Does this feasibly run in a local dev environment? Is this extra development/infra effort?**
  - Happo tests can be triggered from local dev environments, but the report and screenshots will be generated in Happo's cloud environment. If we don't want to use the on-premise version, no extra dev/infra effort is required.

## Pros and Cons of the Alternatives

### Continue using Loki (do nothing)

- `+` Requires no effort because no changes
- `+` Loki is free, and runs locally so it has minimal security impact
- `-` Continues to cost engineering & CI time having to run and approve changes locally
- `-` Visual changes continue to be opaque since there is no UI for other team members to review them

### Switch to Chromatic

- `+` Tests run on Chromatic's hosted platform, so engineers and CircleCI don't need to spend time/resources running them
- `+` Visual tests run against Chrome & Firefox on all pricing plans
- `+` Chromatic provides a UI for users with access to view and approve or reject changes, and approving does not require committing changes
- `+` Chromatic is made by the same team as Storybook, and has first-class support for design systems and additional collaboration & documentation tools for teams
- `-` Chromatic only tests Storybook components, so testing application screens is not an option
- `-` Chromatic costs money (pricing tiers are $149 / $349 / \$649 / month)
- `-` Testing on IE11 is only available in the Pro plan (most expensive)

### Use SauceLabs visual testing

- `+` We already have a SauceLabs account that we use for manual cross-browser testing, and consolidating services is good
- `+` SauceLabs has recently acquired [Screener.io](https://screener.io/) (another automated visual testing service) and is adding a [visual testing product](https://saucelabs.com/platform/visual-testing) to its services
- `-` As far as I can tell, the offering is not yet ready and there is no way right now to sign up as a new customer as of this writing
- `-` SauceLabs also offers automated E2E testing, but we are not using it (using Cypress instead)
- Ultimately I think down the road we will want to make some decisions around what services we use for what tests, and consolidate technologies where possible, but this is not that moment. This ADR may be re-assessed once SauceLabs launches their visual testing service and we decide to re-evaluate.

### Stop visual regression testing

- `+` We can stop using engineer & CI time/resources on automated visual testing
- `-` We lose insight into visual changes and have to find another way to do cross-browser testing
