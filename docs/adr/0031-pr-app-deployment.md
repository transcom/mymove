# Pull Request App Deployment

**User Story:** [162125198](https://www.pivotaltracker.com/story/show/162125198)

We want a way for product and design to be able to explore the results of pull requests without requiring engineers to walk them through it. This will allow product and design to more firmly test pull requests before they get merged and deployed to staging and production without spending engineer time.

## Considered Alternatives

* Heroku Review Apps
* Re-purpose Terraform configs to create Per-Branch Environments
* No change

## Decision Outcome

* Chosen Alternative: No change at the moment.
* Design reviews have led to stronger communication and quicker turn-around of missed design implementations in pull requests. The setup time for this doesn't change as there is nothing to change in design, engineer or product process. This still gives us the option value to choose to build our own thing on ECS containers once we have the time to do it.
* Design and product still need to coordinate with engineers. Engineers can't merge in PRs as fast as they'd like due to busy design and product schedules.

## Pros and Cons of the Alternatives <!-- optional -->

### Heroku Review Apps

This was to use the Heroku pipelines with Github integration enabled with automatic review app deploy. We allotted 2 hours to work on this, and in that time, there was investigation into security risks, rather than actual deployment. If given a full week, I believe we could've gotten something up, but there are still the following cons as well.

* `+` Automatic Heroku pr review app deploys when PR is created.
* `+` Slick CLI and web based tooling to edit.
* `+` Unique url per review app that a design or product person can explore at their leisure.
* `-` Adds another third-party service with potential for security holes. It needs to have SOC2 certification and we need to get a copy of the report under a NDA.
* `-` Need time to educate engineers how to use.
* `-` Need time to educate product and design how to find and help themselves.
* `-` Since the repo has 3 apps in one, it would be tricky to figure out how to configure the Heroku pipeline to automatically deploy review apps for a given app while still using one Proc file (or figuring out how to deploy the `right` one that needs to be tested).
* `-` Since we use multiple hostnames, figuring out how to namespace the different hostnames to different custom domains could also be difficult.
* `-` Changing or adding environment variables would require changes in chamber as well as Heroku in order to make sure it builds, deploys, and runs properly.
* `-` Logs are obfuscated since it's deployed differently, so it's more difficult to get useful data about edge cases.
* `-` Unclear what the cost would be given variability of number of pull requests.
* `-` Unclear how we'd manage permissions to deploy (or shared users), which increases risks.

### Re-purpose Terraform configs to create Per-Branch Environments

* `+` It would be deployed in a similar fashion to our staging and prod environments (same hardware and deploy methods leads to better confidence in success to staging and prod).
* `+` We don't have to get new approval for a third party.
* `+` We can adapt this strategy to our needs as they arise.
* `+` All configuration can be done in one set of tooling.
* `+` We already have permissions mostly set up for AWS.
* `-` Infra has higher priorities on their backlog.
* `-` When things go wrong, we have to manage it.
* `-` Not clear how we'd manage ECS tasks and that would require investigation.
* `-` Would require building a system for per-branch DNS entries.
* `-` Would require a new RDS instance and creating per-branch databases.
* `-` Lifecycle management (when to destroy a deployed branch deploy) might be difficult.
* `-` Unclear what the cost would be.
