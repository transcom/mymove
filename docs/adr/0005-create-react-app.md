# Using Create React App

After landing on our decision to use React for this project, we wanted a template to start from that would provide as many sane defaults as possible.

## Considered Alternatives

* Rolling our own React Setup

## Decision Outcome

Create React App (or CRA) is a well-supported option for application bootstrapping that works on MacOS, Windows, and Linux. Still actively maintained by Facebook, CRA spares us from having to set up Webpack and also includes Jest, a popular testing framework, with no additional setup. A React-approved and widely used file structure is immediately created when the app is created, and `yarn test` is immediately available for automated testing. CRA is well documented and widely used, providing valuable support as the team trains up on React. We can also `eject` the app if we later decide to change the defaults provided by CRA, so a later change in approach is relatively simple.

## Pros and Cons of the Alternatives

### Rolling Our Own React Setup

* `+` Infinitely customizable
* `-` Infinitely customizable
* `-` Have to make choices we may not be ready to make yet
