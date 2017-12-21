# Personal Property Prototype Web Client

## Development

This project was bootstrapped with [create-react-app](https://github.com/facebookincubator/create-react-app).

Prerequisites:

* [pre-commit](http://pre-commit.com/) for running git pre-commit checks.
  * Install on MacOS with: `brew install pre-commit`
  * Install the shell linter with `brew install shellcheck`
* We use Prettier to auto-format Javascript with a pre-commit hook. Make sure [your editor](https://prettier.io/docs/en/editors.html) is configured to use it!

Getting started:

* Run `pre-commit install` to install git pre-commit checks.
* Enter the `client/` directory
* Run `yarn install` to install the dependenices.
* Run `yarn start` to run a local development server with live-refresh.
* Run `yarn test` to run the test suite with live-refresh.
