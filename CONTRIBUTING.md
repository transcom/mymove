# Contributing to MyMove

Anyone is welcome to contribute code changes and additions to this project. If you'd like your changes merged into the master branch, please read the following document before opening a [pull request][pulls].

There are several ways in which you can help improve this project:

1. Fix an existing [issue][issues] and submit a [pull request][pulls].
1. Review open [pull requests][pulls].
1. Report a new [issue][issues]. _Only do this after you've made sure the behavior or problem you're observing isn't already documented in an open issue._

## Table of Contents

- [Getting Started](#getting-started)
- [Making Changes](#making-changes)
- [Code Style](#code-style)
- [Legalese](#legalese)

## Getting Started

MyMove is a [React](https://reactjs.org/) (version ???) and [Go](https://golang.org/) (version ???) application with a [PostgreSQL](https://www.postgresql.org) database (version ???). Development dependencies are managed using ???.

If you're using macOS, you can use [Postgres.app](https://postgresapp.com) to quickly get PostgreSQL installed and running. You may also install PostgreSQL using [Homebrew](https://brew.sh):

```sh
brew install postgresql
```

If you're using Homebrew, consider installing [Homebrew Services](https://github.com/Homebrew/homebrew-services) to easily start and stop PostgreSQL:

```sh
brew tap homebrew/services
brew services run postgresql
```

This project uses React (version ???) which can be installed using ??? like []().

```sh

```

This project uses Go (version ???) which can be installed using ??? like []().

Once you've installed React and Go, install the dependency thing:

```sh

```

With your PostgreSQL server started, run the following setup script from the root of the project:

```sh
bin/setup
```

This script will install the dependencies specified in the project's thing, set development environment configuration variables, and create and migrate the application's databases.

Lastly, start the application by `typing something` and opening [http://localhost:3000](http://localhost:3000) in your Web browser of choice.

### Environment Variables

The `bin/setup` script will generate (or update) a `.env` file in the root of the project and create several default environment variables.

To obtain a Google Maps API key, visit [Google's API Console](https://console.developers.google.com) and select "Credentials" in the left-hand navigation. Click "Create Credentials" and choose "API Key" from the list of options. Copy the generated API key and update the `.env` file, replacing `<your api key here>` with the generated API key:

```
export GOOGLE_MAPS_API_KEY=<your api key here>
```
Enable the Google Maps Geocoding API by going [here](https://console.developers.google.com/apis/api/geocoding_backend).

To work with encrypted seed data, you need to know the secret `SEEDS_ENC_KEY` variable and set it in the `.env` file. This was originally generated once during the `bin/setup` process. Once set, run `bin/rails db:setup` to load the encrypted data.

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

We use a number of tools to evaluate the quality and security of this project's code:

- The test suite uses [RSpec](http://rspec.info) (`bin/rspec`).
- Static code analysis uses [RuboCop](https://github.com/bbatsov/rubocop) (`bin/rubocop`).
- Static vulnerability scans use [Brakeman](http://brakemanscanner.org) (`bin/brakeman`).

Before submitting a [pull request][pulls], use the above tools to verify your changes.

## Code Style

Code formatting conventions are defined in the `.editorconfig` file which uses the [EditorConfig](http://editorconfig.org) syntax. There are [plugins for a variety of editors](http://editorconfig.org/#download) that utilize the settings in the `.editorconfig` file. It is recommended that you install the EditorConfig plugin for your editor of choice.

Your bug fix or feature addition won't be rejected if it runs afoul of any (or all) of these guidelines, but following the guidelines will definitely make everyone's lives a little easier.

## Legalese

Before submitting a pull request to this repository for the first time, you'll need to sign a [Developer Certificate of Origin](https://developercertificate.org) (DCO). To read and agree to the DCO, you'll add your name and email address to [CONTRIBUTORS.md][contributors]. At a high level, this tells us that you have the right to submit the work you're contributing in your pull request and says that you consent to us treating the contribution in a way consistent with the license associated with this software (as described in [LICENSE.md][license]) and its documentation ("Project").

You may submit contributions anonymously or under a pseudonym if you'd like, but we need to be able to reach you at the email address you provide when agreeing to the DCO. Contributions you make to this public Department of Defense repository are completely voluntary. When you submit a pull request, you're offering your contribution without expectation of payment and you expressly waive any future pay claims against the U.S. Federal Government related to your contribution.

[contributors]: https://github.com/deptofdefense/move.mil/blob/master/CONTRIBUTORS.md
[gemfile]: https://github.com/deptofdefense/move.mil/blob/master/Gemfile
[issues]: https://github.com/deptofdefense/move.mil/issues
[license]: https://github.com/deptofdefense/move.mil/blob/master/LICENSE.md
[pulls]: https://github.com/deptofdefense/move.mil/pulls
