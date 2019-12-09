# Contributing to MyMove

Anyone is welcome to contribute code changes and additions to this project. If you'd like your changes merged into the master branch, please read the following document before opening a [pull request](https://github.com/transcom/mymove/pulls).

There are several ways in which you can help improve this project:

1. Fix an existing [issue](https://github.com/transcom/mymove/issues) and submit a [pull request](https://github.com/transcom/mymove/pulls).
1. Review open [pull requests](https://github.com/transcom/mymove/pulls).
1. Report a new [issue](https://github.com/transcom/mymove/issues). _Only do this after you've made sure the behavior or problem you're observing isn't already documented in an open issue._

## Table of Contents

- [Getting Started](#getting-started)
- [Making Changes](#making-changes)
- [Code Style](#code-style)
- [Legalese](#legalese)

## Getting Started

See our [Development Setup](https://github.com/transcom/mymove#development)

## Making Changes

1. Fork and clone the project's repo.
1. Install development dependencies as outlined above.
1. Create a feature branch for the code changes you're looking to make: `git checkout -b your-descriptive-branch-name origin/master`.
1. _Write some code!_
1. Run the application and verify that your changes function as intended: `something`.
1. If your changes would benefit from testing, add the necessary tests and verify everything passes by running `something`.
1. Commit your changes: `git commit -am 'Add some new feature or fix some issue'`. _(See [this excellent article](https://chris.beams.io/posts/git-commit) for tips on writing useful Git commit messages.)_
1. Push the branch to your fork: `git push -u origin your-descriptive-branch-name`.
1. Create a new pull request and we'll review your changes.

### Verifying Changes

We use a number of tools to evaluate the quality and security of this project's code. Before submitting a pull request, be sure that `make test` runs without error or test failure. Additionally, be sure that all of the `pre-commit checks pass`.

## Code Style

Please review our [front end](https://github.com/transcom/mymove/blob/master/docs/frontend.md) and [back end](https://github.com/transcom/mymove/blob/master/docs/backend.md) coding guidelines and do your best to follow the conventions and choices described therein.

Code formatting conventions are defined in the `.editorconfig` file which uses the [EditorConfig](http://editorconfig.org) syntax. There are [plugins for a variety of editors](http://editorconfig.org/#download) that utilize the settings in the `.editorconfig` file. It is recommended that you install the EditorConfig plugin for your editor of choice.

## Legalese

Before submitting a pull request to this repository for the first time, you'll need to sign a [Developer Certificate of Origin](https://developercertificate.org) (DCO). To read and agree to the DCO, you'll add your name and email address to [CONTRIBUTORS.md](https://github.com/transcom/mymove/blob/master/CONTRIBUTORS.md). At a high level, this tells us that you have the right to submit the work you're contributing in your pull request and says that you consent to us treating the contribution in a way consistent with the license associated with this software (as described in [LICENSE.md](https://github.com/transcom/mymove/blob/master/LICENSE.md)) and its documentation ("Project").

You may submit contributions anonymously or under a pseudonym if you'd like, but we need to be able to reach you at the email address you provide when agreeing to the DCO. Contributions you make to this public Department of Defense repository are completely voluntary. When you submit a pull request, you're offering your contribution without expectation of payment and you expressly waive any future pay claims against the U.S. Federal Government related to your contribution.

[contributors]: https://github.com/transcom/move.mil/blob/master/CONTRIBUTORS.md
[issues]: https://github.com/transcom/move.mil/issues
[license]: https://github.com/transcom/move.mil/blob/master/LICENSE.md
[pulls]: https://github.com/transcom/move.mil/pulls
