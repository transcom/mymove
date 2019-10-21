# How to Searching for Application Errors

Searching for errors from the application can be difficult.  Here is a playbook for finding specific errors.

## Tips

- Always check the timezone of the log messages. Messages from AWS infrastructure are in UTC whereas log messages
  from the app may be in PST.

## Playbook

### Tracking ALB Errors in S3

Watch [a video of searching ALB errors](https://zoom.us/recording/play/2DGo7KYqrvSAYvvEo2--zG_xx93K0ALUMLHJVXJl0X9WRo0lxkXXLAxp8IPHHDA8)
to see how this is done in real time.

**NOTE:** ALB logs are only stored in AWS S3, so the only method to search them is to download locally and search via the command line.

Often we'll see 500 errors from the application that trigger the `alb-staging-target-5xx-limit` alarm on the ALB.
The text is not super helpful:

```txt
Alarm Name
alb-staging-target-5xx-limit
Alarm Description
Target(s) behind ALB in staging are returning 5xx error codes.
Alarm reason
Threshold Crossed: 1 datapoint [2.0 (07/01/19 19:39:00)] was greater than or equal to the threshold (2.0).
Old State
OK
Current State
ALARM
Link to Alarm
https://console.aws.amazon.com/cloudwatch/home?region=us-west-2#alarm:alarmFilter=ANY;name=alb-staging-target-5xx-limit
```

You can download ALB logs for a specific date by using our `download-alb-logs` command line tool with:

```sh
cd transcom-ppp
../scripts/download-alb-logs tmp prod 2019/01/09,2019/01/10
```

Alternatively, you can run the `sync` commands manually as:

```sh
mkdir -p tmp/
aws s3 sync s3://transcom-ppp-aws-logs/alb/app-staging/AWSLogs/923914045601/elasticloadbalancing/us-west-2/2019/01/09/ ./tmp/
aws s3 sync s3://transcom-ppp-aws-logs/alb/app-staging/AWSLogs/923914045601/elasticloadbalancing/us-west-2/2019/01/10/ ./tmp/
```

Then search for the errors in the ALB logs using `big-cat`, `gunzip`, `read-alb-logs`, and `jq`.

```sh
make build_tools
export http_code=500
big-cat ./tmp/*.log.gz | gunzip | read-alb-logs | jq ". | select( .elbStatusCode | startswith(\"${http_code}\")) | {timestamp, clientPort, elbStatusCode, targetStatusCode, request, actionsExecuted}"
```

And you'll see events like this:

```sh
{
  "timestamp": "2019-06-04T16:49:53.728140Z",
  "clientPort": "193.232.106.88:47233",
  "elbStatusCode": "403",
  "targetStatusCode": "-",
  "request": "GET https://54.201.161.247:443/.env HTTP/1.1",
  "actionsExecuted": "waf"
}
{
  "timestamp": "2019-06-04T17:13:30.820491Z",
  "clientPort": "180.248.30.194:57707",
  "elbStatusCode": "403",
  "targetStatusCode": "-",
  "request": "GET http://54.201.161.247:80/ HTTP/1.1",
  "actionsExecuted": "waf"
}
```

You can use any `jq` filter you want on the output data like filtering on timestamp.

A simpler way to scan the logs is to use a helpful script we've built:

```sh
cd transcom-ppp
../scripts/scan-alb-logs staging 500 2019/01/09
```

#### Understanding ALB Log Entries

References:

- [AWS ALB Log Entries](https://docs.aws.amazon.com/elasticloadbalancing/latest/application/load-balancer-access-logs.html)

There are two forms of `elb_status_code target_status_code` that you'll see on a log line, `200 200` or `301 -`. The first
is the ELB status code, which is what the user sees.  The second is the Target status code, which is what the application
container is returning. If the Target status code is `-` that means that the ALB did not forward the traffic from the
user to the application container or the application container failed to respond.

The most common reason you're looking at these logs is to debug 5XX errors. Here is a short key for things you'll see:

| Code | Common meaning |
| --- | --- |
| 500 | The application is returning a 500 from an unhandled error. |
| 501 | The ALB returns a 501 when it receives an [unsupported Transfer-Encoding header](https://docs.aws.amazon.com/elasticloadbalancing/latest/application/load-balancer-troubleshooting.html#http-501-issues), e.g, `Transfer-Encoding: defalate`.   |
| 502 | The container was recycled because the ALB told ECS that the container was unhealthy. |

You can manually generate a 501 error with a deflate encoding: `curl -v -X POST --http1.1 -H "Transfer-Encoding: deflate" https://example.com`.

Its rare but sometimes the WAF has problems. You can check the `actions_executed` field in the ALB log to find out
what has happened. Normal traffic actions are `waf,forward`, which indicates the WAF analyzed the request and forwarded it to the container. If traffic was blocked by the WAF, you'll just see `waf`.  If the WAF itself had problems (which happens very rarely), you'll see `waf-failed`.

### Tracing Application Errors in CloudWatch

After finding the ALB error messages the next best thing to do is to open up CloudWatch Logs and look into either
`ecs-tasks-app-<environment>` or `ecs-tasks-app-client-tls-<environment>`.  As an example you can look into the logs
for Staging here:

[ecs-tasks-app-staging log streams](https://us-west-2.console.aws.amazon.com/cloudwatch/home?region=us-west-2#logStream:group=ecs-tasks-app-staging)

It's worth noting that each container running in Fargate creates a unique Log Stream. That means you'll have to
look through several log streams to find the events corresponding to the alarm.  The best way to do this is to
match up the `Last Event Time` with the time listed in the ALB alarm.

Open up the log stream and search for `{ $.resp-status > 499 }` in the search bar.  Here's an example:

[errors in log stream](https://us-west-2.console.aws.amazon.com/cloudwatch/home?region=us-west-2#logEventViewer:group=ecs-tasks-app-staging;stream=app/app-staging/c0fce04c-d248-4e90-b7c0-a4a5db187c68;filter=%7B%20$.resp-status%20%3E%20499%20%7D;start=2019-01-06T20:40:04Z)

You can glean a couple useful pieces of information from the log message.  Look especially at these items if they exist:

- git-commit
- host
- url
- user-id
- service-member-id
- office-user-id
- admin-user-id

Most logged error messages come after the log lines that have useful information.  The best way to find those lines
is to limit the log stream to a period of time that surrounds the error message you're looking at.

With the URL route and the user information you can usually track down the piece of code that is causing the issue.

### Tracking down users in Staging

If you have the `office-user-id` or `admin-user-id` you can run `make run_prod_migrations` to download the secret
prod migrations locally and match the users.
