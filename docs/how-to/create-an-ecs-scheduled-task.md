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
and make a similar target that matches your command in the section for `SCHEDULED TASK TARGETS` in the `Makefile.
It is important that you pass in only the environment variables (designated by `-e`) that are needed for your task
and no more. Be as specific as possible. Then when you invoke `make tasks_your_subcommand` you should see the
`Dockerfile.tasks_local` file built and your command run.

**WARNING:** Do not modify the `Dockerfile.tasks` or `Dockerifle.tasks_local` as they will already have your
code built into the `milmove-tasks` binary that is included. There is no need to modify these files.

## Provisioning with Terraform

Once you have confirmed that your code works locally it is time to update the AWS terraform code. Work with someone
on the infrastructure team to add [an ecs-scheduled-task module](https://github.com/transcom/ppp-infra/blob/6e57f84b937376a1a5f4556869304ed81f453ef4/modules/aws-app-environment/main.tf#L538-L557)
that references your subcommand.

**WARNING:** It is CRITICAL that the subcommand name you used earlier matches the `task_name` given to the module.
If the name doesn't match then you won't see the task work and the error messages you get will not be very helpful.

You will also want to write a cron job expression that matches how often you need to run the task. AWS has documentation
on [Cron Expressions](https://docs.aws.amazon.com/AmazonCloudWatch/latest/events/ScheduledEvents.html#CronExpressions) that
you should reference as it is slightly different than normal cron expressions on Unix.

Lastly the Infra team needs to give your scheduled task access to the [RDS IAM Role](https://github.com/transcom/ppp-infra/blob/6e57f84b937376a1a5f4556869304ed81f453ef4/modules/aws-app-environment/main.tf#L213-L216). This has to happen after the ECS task is provisioned due to some odd behavior
of Terraform v0.11.

## Updating CircleCI to Deploy


