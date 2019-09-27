# How To Create An ECS Scheduled Task

An ECS Scheduled Task is similar to a cron job. It runs a docker container inside our ECS Cluster in AWS on a
schedule defined by an CloudWatch Rule (typically a cron rule). These tasks can do things like update the DB on a
regular schedule or send out emails.

## Writing the Function

You will need to start out by writing a new function to do the task locally. The code for this function
should go into the `cmd/milmove-tasks/` directory and be named similar to the command that will kick off your
scheduled task. For instance `save_fuel_price_data.go` is the code for running the ECS Scheduled Task for
saving fuel price data. You integrate this code into the `cmd/milmove-tasks/main.go` file using the
[spf13/cobra](github.com/spf13/cobra) library as a sub command. Once the binary has been compiled with
`make bin/milmove-tasks` your command should be runnable by calling `bin/milmove-tasks save-fuel-price-data` or
similar based on the name you chose.  Since you're writing a sub-command the flags for that command will be
unique to your command alone and should be contained withing your golang file.

**WARNING:** It is very important that your sub-command name is less than or equal to 25 characters. This same name
is referenced by the CircleCI and AWS terraform code to provision your task. The name must be identical in all
cases and names longer than 25 characters will throw errors when provisioning.

## Test the code locally

You should be able to test the code locally either by running the code directly via `bin/milmove subcommand` or
via the `Makefile` targets you will need to create. In the file you should copy the target `tasks_save_fuel_price_data`
and make a similar target that matches your command in the section for `SCHEDULED TASK TARGETS` in the `Makefile`.
It is important that you pass in only the environment variables (designated by `-e`) that are needed for your task
and no more. Be as specific as possible. Then when you invoke `make tasks_your_subcommand` you should see the
`Dockerfile.tasks_local` file built and your command run.
that references your subcommand.

**WARNING:** It is CRITICAL that the subcommand name you used earlier matches the `task_name` given to the module.
If the name doesn't match then you won't see the task work and the error messages you get will not be very helpful.

You will also want to write a cron job expression that matches how often you need to run the task. AWS has documentation
on [Cron Expressions](https://docs.aws.amazon.com/AmazonCloudWatch/latest/events/ScheduledEvents.html#CronExpressions) that
you should reference as it is slightly different than normal cron expressions on Unix.

Lastly the Infra team needs to give your scheduled task access to the [RDS IAM Role](https://github.com/transcom/ppp-infra/blob/6e57f84b937376a1a5f4556869304ed81f453ef4/modules/aws-app-environment/main.tf#L213-L216). This has to happen after the ECS task is provisioned due to some odd behavior
of Terraform v0.11.

Your changes will now need to be deployed to all three environments: experimental, staging, and prod. This doesn't mean
your code will begin running as the tasks haven't been deployed by CircleCI at this point. That means its safe to provision
all of the environments in advance.

### AWS IAM Permissions

The IAM Permissions for the task are managed in a [Task Role Policy](https://github.com/transcom/ppp-infra/blob/6e57f84b937376a1a5f4556869304ed81f453ef4/modules/aws-app-environment/ecs-scheduled-task/main.tf#L149-L200).
You would modify the policy if there is an action that your task needs that isn't already available (like DB, S3, or
SES access). If you need to modify these to enable use of a new AWS resource please work with Infrastructure to
properly scope the access to be as narrow as possible for your task's use case (ie no "*" resources if possible).
This is because your changes will be available to all the other tasks and may inadvertently give access to a resource
that another task shouldn't have access to.  When in doubt please work with Infrastructure.

### AWS Resource Throttling

Many of the resources in AWS have built-in throttling. For example you can only send up to 80 emails per second with
SES. Your code needs to be aware of these throttling limits and take into account errors when throttling limits are
hit. It's also important that you add logging to your code so it's easier to determine what went wrong in your code
when working with AWS resources.

## Updating CircleCI to Deploy

To deploy with CircleCI you need to modify the [deploy_tasks_steps](https://github.com/transcom/mymove/blob/d676b217ea67dfd893d770a77bb9e2d898d0b891/.circleci/config.yml#L89-L98)
in the `.circleci/config.yml`. Each scheduled task needs a separate `deploy:` section in this step. The important
parts to modify in the `deploy` section `command` argument is the flag `--command` and `--command-args`. These are
used by the `bin/ecs-deploy-task-container` code to create a Task Definition in ECS that is similar to a Dockerfile
with `--command` mapping to `entrypoint` and `--command-args` mapping to `command` in a Dockerfile.

**Note:** You can deploy your code to the experimental environment to test things out at this point. It's suggested that
you update the `cron` rule in terraform to run once per minute (if that doesn't break throttling thresholds) so that you
can more quickly see if your code is doing what is expected.

Finally merge your code with the `master` branch and it will deploy in the regular manner of staging and then prod.

## Viewing logs

Similar to the application all logs for ECS Scheduled Tasks will end up in AWS CloudWatch Logs.  For example,
if you want to see the scheduled task logs in experimental you'd go to the `ecs-tasks-app-experimental` Log Group
and search for a log stream named after your subcommand. For instance, with `save-fuel-price-data` the log stream is
named `app-tasks/app-tasks-save-fuel-price-data-experimental/TASKID` (where `TASKID` is a random string that
corresponds to the task ID of the last run task). You can look at the `Last Event Time` column to find the log
stream you are looking for. Open up the log stream and look at your data.

**WARNING:** CloudWatch Insights is not super valuable for these kinds of logs. In general scheduled tasks run only
once per day or less and have very few log lines. If you use CloudWatch Insights you'll get logs across multiple
log streams and not necessarily in an order that makes sense for the execution of the task. Please avoid using
CloudWatch Insights when debugging your task if you wish to save yourself a lot of pain.
