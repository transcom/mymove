# How to query cloudwatch logs

Cloudwatch captures a large footprint of logs that orginate from sources like alb, rds, application logs etc. The verbose nature of logs makes finding actionable information hard. This document has some tips and tricks on how to identify and filter data effectively.

## Steps

* Select [cloudwatch](https://us-west-2.console.aws.amazon.com/cloudwatch/home?region=us-west-2) from aws console and then choose [insights](https://us-west-2.console.aws.amazon.com/cloudwatch/home?region=us-west-2#logs-insights:queryDetail=~(end~0~start~-3600~timeType~'RELATIVE~unit~'seconds~editorString~'fields*20*40timestamp*2c*20*40message*0a*7c*20sort*20*40timestamp*20desc*0a*7c*20limit*2020~isLiveTail~false~queryId~'*2faws*2flambda*2faws-health-notifier-prod~source~'*2faws*2flambda*2faws-health-notifier-prod))
* Select log stream you would like to query for example: select `/aws/rds/instance/app-staging/postgresql` if you want to query against staging sql logs
* Select the time period you want to query against (15 mins, 30 mins, 1 day etc)
* Click on "Run Query"

## Use cases

### [See last 100 sql queries in a time period for staging environment](https://us-west-2.console.aws.amazon.com/cloudwatch/home?region=us-west-2#logs-insights:queryDetail=~(end~'2019-04-06T03*3a59*3a59.999Z~start~'2019-04-02T04*3a00*3a00.000Z~timeType~'ABSOLUTE~tz~'Local~editorString~'fields*20*40message*0a*7c*20parse*20*22statement*3a*20*2a*22*20as*20statement*0a*7c*20filter*20statement*20like*20*22SELECT*20*22*0a*7c*20filter*20statement*20not*20like*20*22SELECT*201*3b*22*0a*7c*20sort*20*40timestamp*20desc*0a*7c*20limit*20100~isLiveTail~false~queryId~'*2faws*2frds*2finstance*2fapp-staging*2fpostgresql~source~'*2faws*2frds*2finstance*2fapp-staging*2fpostgresql))

```sql
fields @message
| parse "statement: *" as statement
| filter statement like "SELECT "
| filter statement not like "SELECT 1;"
| sort @timestamp desc
| limit 100
```

### [Find last two sql queries in staging that were executed aginst **moves** table](https://us-west-2.console.aws.amazon.com/cloudwatch/home?region=us-west-2#logs-insights:queryDetail=~(end~'2019-04-02T18*3a28*3a39.452Z~start~'2019-04-02T16*3a07*3a22.191Z~timeType~'ABSOLUTE~tz~'Local~editorString~'fields*20*40message*0a*7c*20parse*20*22statement*3a*20*2a*22*20as*20statement*0a*7c*20filter*20statement*20like*20*22FROM*20moves*22*0a*7c*20filter*20statement*20not*20like*20*22SELECT*201*3b*22*0a*7c*20sort*20*40timestamp*20desc*0a*7c*20limit*202~isLiveTail~false~queryId~'*2faws*2frds*2finstance*2fapp-staging*2fpostgresql~source~'*2faws*2frds*2finstance*2fapp-staging*2fpostgresql))

```sql
fields @message
| parse "statement: *" as statement
| filter statement like "FROM moves"
| filter statement not like "SELECT 1;"
| sort @timestamp desc
| limit 2
```