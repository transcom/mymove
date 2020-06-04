# Use CircleCI to build only Pull Requests and master

Currently our CircleCI costs are very high, well over the initial expected budget. By default CircleCI builds every branch regardless of if it has a Pull Request (PR) or not. Given the size of the MilMove team and our typical pattern of creating branches in the main repo this results in a large number of builds. The builds are also triggered every time a branch changes, which can be quite frequent on active development branches.

## Considered Alternatives

* Do nothing, keep the default settings for CircleCI
* Switch on the CircleCI option to only build PRs and our default branch (master)
* Keep building all branches and review build pipeline

## Decision Outcome

*Chosen Alternative:* Switch on the CircleCI option to only build PRs and our default branch (master)A

Once accepted this switch `CircleCI -> Project Settings -> Advanced -> Only build pull requests` should be enabled.

### How do we undo this

We turn the `CircleCI -> Project Settings -> Advanced -> Only build pull requests` option off.

## Pros and Cons of the Alternatives

### Do nothing, keep the default settings for CircleCI

* `+` Easiest to do nothing
* `-` Costs of CircleCI continue to grow
* `-` Building all branches will continue building even branches that are not ready for review

### Switch on the CircleCI option to only build PRs and our default branch (master)

* `+` Easy to implement, just a project setting in CircleCI
* `+` Should reduce our CircleCI costs significantly by only building PR branches
* `+` CircleCI usage reduction will buy us time before having to review all stages of the pipeline
* `+` Encourage people to push more often not just when ready for review
* `-` Lose the ability to have our branches run in CircleCI without creating a PR

### Keep building all branches and review build pipeline

* `+` Review of the pipeline could lead to pipeline efficiencies
* `-` Very time consuming to review
* `-` Still builds every branch on every commit to each branch, even without a PR
