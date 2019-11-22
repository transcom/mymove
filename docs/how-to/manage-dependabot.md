# How to Manage Dependabot

[Dependabot](https://dependabot.com) is used to monitor the repository dependencies and update them with automatic
pull requests against the `master` branch in the repo. The configuration is done via a file named
`.dependabot/config.yml`. Read more about [dependabot configuration](https://dependabot.com/docs/config-file/) in the
docs.

## Security

We use dependabot as part of our security measures. It ensures that the repository dependencies are up to date and
that any security vulnerabilities are caught as soon as new versions are published. Dependabot will even
add security release information in the text of the PR.

## Organization Level Settings

The settings for the Transcom organization can be found in the [Account Settings](https://app.dependabot.com/accounts/transcom/settings)
page. These manage settings for all repos under Transcom.

## Repo Management

Repo management should be done in the `.dependabot/config.yml` file. However, you can view and interact with
settings temporarily via the [repo management page](https://app.dependabot.com/accounts/transcom/repos/114694829).
This is a good place to try out new features without having to push a PR to the repository.
