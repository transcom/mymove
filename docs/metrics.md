# MilMove Application Metrics

This document provides an overview of the metrics collected by the MilMove application.

## Dashboards

Various AWS metrics have been aggregated in CloudWatch dashboards.
These include data about network requests, container resources, and errors.

* [Prod Dashboard](https://us-west-2.console.aws.amazon.com/cloudwatch/home?region=us-west-2#dashboards:name=mil-prod)
* [Staging Dashboard](https://us-west-2.console.aws.amazon.com/cloudwatch/home?region=us-west-2#dashboards:name=mil-staging)
* [Experimental Dashboard](https://us-west-2.console.aws.amazon.com/cloudwatch/home?region=us-west-2#dashboards:name=mil-experimental)

### Updating Dashboards

Dashboards are maintained by Terraform in the [MilMove infra repo](https://github.com/transcom/ppp-infra/blob/master/modules/aws-app-environment/main.tf#L840).
When changes are made in the AWS CloudWatch UI,
they can be exported by clicking "Actions" > "View/edit source".
The Terraform file can be updated with this source,
replacing environment variables as shown in the existing file.
Feel free to ask the Infrastructure team for help,
as they'll have access to deploy changes across environments.
