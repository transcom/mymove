# Pull Request App Deployment

**User Story:** [162125198](https://www.pivotaltracker.com/story/show/162125198)

We want a way for product and design to be able to explore the results of pull requests without requiring engineers to walk them through it. This will allow product and design to more firmly test pull requests before they get merged and deployed to staging and production without spending engineer time.

## Considered Alternatives

* Heroku
* EC2 containers
* No change

## Decision Outcome

* Chosen Alternative: No change at the moment.
* Design reviews have led to stronger communication and quicker turn-around of missed design implementations in pull requests. The setup time for this doesn't change as there is nothing to change in design, engineer or product process.
* Design and product still need to coordinate with engineers. Engineers can't merge in PRs as fast as they'd like due to busy design and product schedules.

## Pros and Cons of the Alternatives <!-- optional -->

### Heroku

* `+` Automatic Heroku pr review app deploys when PR is created.
* `+` Slick CLI and web based tooling to edit.
* `+` Unique url per review app that a design or product person can explore at their leisure.
* `-` Adds another third-party service with potential for security holes.
* `-` Need time to educate engineers how to use.
* `-` Need time to educate product and design how to find themselves.
* `-` Since the repo has 3 apps in one, it would be tricky to figure out how to configure the Heroku pipeline to deploy review apps while still
* `-`

### EC2 Containers

* `+` It would be deployed in a similar fashion to our staging and prod environments (same hardware and deploy methods)
* `+` We don't have to get new approval for a third party.
* `+` We can adapt this strategy to our needs as they arise.
* `-` Infra has higher priorities on their backlog.

### No change

* `+` Nothing to change
* `+` *[argument 2 pro]*
* `-` *[argument 1 con]*
* *[...]* <!-- numbers of pros and cons can vary -->
