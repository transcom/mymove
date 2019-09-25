# ecs-service-logs

## Description

**ecs-service-logs** is used to filter JSON-formatted log lines in CloudWatch.

## Usage

Easily filter JSON formatted application logs from an ECS Service or Task.  This tool compiles a chain of filters into a filter pattern in the format used by CloudWatch Logs.  You can filter application logs by ECS Cluster (--cluster), ECS Service (--service), and environment (--environment).  When filtering logs for a stopped task, use "--status STOPPED".  Trailing positional arguments are added to the query.  Equality (X=Y) and inverse equality (X!=Y) are supported.  Wildcards are also supported, e.g, "url!=health*".

[https://docs.aws.amazon.com/AmazonCloudWatch/latest/logs/FilterAndPatternSyntax.html](https://docs.aws.amazon.com/AmazonCloudWatch/latest/logs/FilterAndPatternSyntax.html).

```shell
Usage:
  ecs-service-logs show [flags] [msg=XYZ] [referer=XYZ]...

Flags:
      --aws-profile string               The aws-vault profile
      --aws-region string                The AWS Region (default "us-west-2")
      --aws-vault-keychain-name string   The aws-vault keychain name
  -c, --cluster string                   The cluster name
  -f, --ecs-task-def-family string       The ECS task definition family.
  -r, --ecs-task-def-revision string     The ECS task definition revision.
  -e, --environment string               The environment name
  -b, --git-branch string                The git branch
      --git-commit string                The git commit
  -h, --help                             help for show
  -l, --level string                     The log level: debug, info, warn, error, panic, fatal
  -n, --limit int                        If 1 or above, the maximum number of log events to print to stdout. (default -1)
  -p, --page-size int                    The page size or maximum number of log events to return during each API call.  The default is 10,000 log events. (default -1)
  -s, --service string                   The service name
      --status string                    The task status: RUNNING, STOPPED, ALL (default "ALL")
  -t, --tasks int                        If 1 or above, the maximum number of log streams (aka tasks) to print to stdout. (default 10)
  -v, --verbose                          Print section lines
```

## Examples

Search for a client IP Address.

```shell
ecs-service-logs show -s app -e staging x-forwarded-for=*1.2.3.4
```

Search for a client IP Address in only running tasks.

```shell
ecs-service-logs show --status RUNNING -c app-staging -s app x-forwarded-for=*1.2.3.4
```

Search for requests in a given environment, but not health checks (url is defined but does not start with /health).

```shell
ecs-service-logs show -s app -e staging url=* url!=/health*
```

Filter by url is defined and git commit.

```shell
ecs-service-logs show -s app -e experimental url=* git_commit=asdfnh98nwuefr9a8jf
```

Filter by url is defined and the number of headers is greater than 14.

```shell
ecs-service-logs show -s app -e experimental url=* "headers>14"
```

Search for requests with an event type, status, and specified time range

```shell
ecs-service-logs show -s app -e experimental event_type="create_office_user" --status=ALL --start-time=2019-09-16T23:43:00Z --end-time=2019-09-16T23:43:20Z
```

**Note:** `event_type` will follow the convention of
`{ACTION}_{SINGULAR_RECORD_TYPE}` where `ACTION` can be one of `create`, `update`, or `delete`
