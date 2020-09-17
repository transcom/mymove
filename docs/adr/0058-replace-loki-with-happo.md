# Use Happo for visual regression testing

## Background

Visual regression testing is a type of automated testing that is meant to catch unintended visual side effects in rendered UI. This is an area of risk due to the global nature of CSS -- when changing CSS code to modify one element, there is always a possibility of inadvertent changes to other elements without knowing anything happened until it is discovered by happenstance after the fact.

Visual regression testing occurs by first taking screenshots of selected portions of rendered UI to use as references, then comparing screenshots with code changes against those references, and finally alerting the team to differences between the images for review. Because any changes, whether they are intended or not, require explicit approval, building visual regression testing into our workflow can help ensure all UI changes have been reviewed by members of the design and product teams, as well as provide a visual of cross-browser implementations without needing manual QA.

## Problem

MilMove currently uses a tool called Loki for performing visual regression tests on all components rendered in Storybook. While Loki has been helpful with catching unintended visual changes, it also has several downsides:

- Loki is executed in local development & CI environments, so reference image assets must be checked into the code repo, and running/updating the tests locally is time- and resource-intensive.
- Loki only runs tests in Chrome, so visual implementations in other browsers (Firefox, IE, Safari) are not tested.
- Review and approval of changes happens in local dev environments, so actual changes can be relatively opaque and may not be explicitly reviewed by designers.
- Loki only supports visual testing of Storybook components, so we would not be able to add visual tests to the application if we wanted.

## Considered Alternatives

- Continue using Loki (do nothing)
- Switch to Happo
- Switch to Chromatic
- Stop visual regression testing

## Decision Outcome

- **Chosen Alternative: Switch to Happo**
- `+` Tests run on Happo's hosted platform, so engineers and CircleCI don't need to spend time/resources running them
- `+` Visual tests can run against all of the browsers we support: Chrome, Firefox, IE11, Safari, iOS Safari, Edge (on all pricing plans, including free open source)
- `+` Happo provides a UI for users with access to view and approve or reject changes
- `+` Happo has a plugin to use with Storybook components, but can also be used to screenshot the application itself (for example, as part of Cypress E2E tests)
- `~` Happo has an on-premise option available if security is a concern (but pricing is by request, and will require additional infra support to set up and maintain)
- `-` Happo runs tests on its hosted platform, so it requires giving Happo access to our Github repo (which is currently public)
- `-` Happo costs money for non-open source projects (pricing tiers are $125 / $250 / \$500 / month depending on usage)
- `-` Test reports for free open source projects are public to anyone who has the link

## Pros and Cons of the Alternatives

### Continue using Loki (do nothing)

- `+` Requires no effort because no changes
- `+` Loki is free, and runs locally so it has minimal security impact
- `-` Continues to cost engineering & CI time having to run and approve changes locally
- `-` Visual changes continue to be opaque since there is no UI for other team members to review them

### Switch to Chromatic

- `+` Tests run on Chromatic's hosted platform, so engineers and CircleCI don't need to spend time/resources running them
- `+` Visual tests run against Chrome & Firefox on all pricing plans
- `+` Chromatic provides a UI for users with access to view and approve or reject changes
- `+` Chromatic is made by the same team as Storybook, and has first-class support for design systems and additional collaboration & documentation tools for teams
- `-` Chromatic only tests Storybook components, so testing application screens is not an option
- `-` Chromatic costs money (pricing tiers are $149 / $349 / \$649 / month)
- `-` Testing on IE11 is only available in the Pro plan (most expensive)

### Stop visual regression testing

- `+` We can stop using engineer & CI time/resources on automated visual testing
- `-` We lose insight into visual changes and have to find another way to do cross-browser testing
