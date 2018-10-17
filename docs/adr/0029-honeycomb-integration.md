# Honeycomb Integration

* Status: accepted
* Deciders: @elliottdotgov, @hlieberman-gov, @ntwyman
* Date: 2018-10-01

## Context and Problem Statement

Currently, the best method for tracing a request through the MyMove site is done by inspecting the JSON formatted logs hosted in AWS CloudWatch Logs. The CloudWatch Logs service provides the ability to define point queries based on the fields in the JSON log entries.  For example, a query could be looking up all log entries with a severity level of “error”. Simple searches like looking up errors do provide some value, particularly when looking for known issues and if the volume of requests flowing through the site is low. However, CloudWatch Logs doesn’t provide engineers a good way to look into more abstract unknown problems that could be happening.

It’s important to be able to ask high level and ad-hoc questions like “what are the top 5 slowest API calls right now?” . More importantly, being able to ask follow up questions like “For those slow API calls, where is that time being spent?”.

This decision record is meant to pick a path forward to increase the visibility on the state of the MyMove site.

## Decision Drivers

### Operational Requirements

* **Outsource to a Software as a Service (SaaS) provider** - One way to solve these problems could be to build a better system in-house (e.g., deploying [Elasticsearch](https://github.com/elastic/elasticsearch ), [Logstash](https://github.com/elastic/logstash), and [Kibana](https://github.com/elastic/kibana)), but the operational overhead needed to maintain such a system is significant and isn’t a primary focus for MyMove.

### Usability Requirements

* **Be able to ask previously unknown questions** - Provides a clear interface for asking both granular and coarse question about the state of the MyMove site.
* **The service should be easy to use** - Doesn’t require learning a domain specific language(DSL) or deep understanding of regular expressions.

### Security Requirements

* **Encrypt all data sent from MyMove** - Service provides encryption at rest and in-transit.
* **Allow MyMove infrastructure to control access to encrypted data before being sent** - MyMove remains in control of how the data is encrypted or hashed before leaving the infrastructure.
* **Provide controls to configure data retention and deletion** - Strict controls on how data is persisted and has controls for deleting columns or narrowed subsets.
* **Provide access control policies for reading, writing or deleting data** - Engineers accessing this service should only need to write and read access to the chosen service. The service should be able to restrict access.

## Considered Options

* **Log aggregation service** - [Sumo Logic](https://www.sumologic.com/)
* **Application Performance Monitoring (APM) service** - [New Relic](https://newrelic.com/)
* **Event based observability service** - [Honeycomb](https://www.honeycomb.io/)

## Decision Outcome

Honeycomb straddles a difficult balance that provides a pretty good user interface for debugging complex datasets, but also provides a mechanism to encrypt any data we send before it leaves the MyMove infrastructure. Below is an example of using Honeycomb to see the 99th percentile request latencies (in milliseconds) for the slowest API calls. We’ve been testing the service in the MyMove staging environment to verify its usefulness before moving forward with a roll out to production. We did this to understand what data we will and won’t send to Honeycomb. Also it allowed us to understand what sorts of security controls are in place to protect and manage the data we do send.

## Pros and Cons of the Options

### [Sumo Logic](https://www.sumologic.com/)

* `+` Integrates with CloudWatch Logs. [1]
* `+` Encryption in-transit and at rest using keys maintained and rotated by Sumo Logic. [2]
* `+` Fine grained role based access controls. [2]
* `-` Doesn’t support encrypting data before leaving MyMove infrastructure.
* `-` Requires managing search indexes and deciding which data is important to index ahead of time.
* `-` Querying JSON logs requires using JSONPath syntax, which isn’t much better than CloudWatch.
* `-` Generating graphs around things like requests latencies requires defining dashboards and configuration ahead of time using complicated "Parse Expressions". [3]
* `-` Very Limited controls around data persistence. All deletion requests must go through support. [4]

### [New Relic](https://newrelic.com/)

* `+` Provides good library for integrating with Go codebase. [5]
* `+` In the process of being FedRAMP approved. [6]
* `+` Provides defining both coarse and fine grained queries.
* `+` Fine grained role based access controls. [7]
* `-` Data is encrypted in-transit, but not at rest. [8]
* `-` Data cannot be deleted once written without an explicit support request. [9]
* `-` Doesn’t support encrypting data before leaving MyMove infrastructure.
* `-` Difficult to extend beyond the quick start functionality.
* `-` Querying the New Relic Events services requires learning New Relic Query Language (NRQL). [10]

### [Honeycomb](https://www.honeycomb.io/)

* `+` Provides good library for integrating with Go codebase. [11]
* `+` Data is encrypted in-transit and at rest. [12]
* `+` Supports encrypting data before leaving MyMove infrastructure using Secure Tenancy. [13]
* `+` Web based UI for querying against any number of fields in the Honeycomb dataset.
* `+` No complex query syntax to learn.
* `-` Limited role based access controls.
* `-` Secure Tenancy requires running new infrastructure in-house.
* `-` Datasets can be deleted through the web UI, but column based deletions need to be done through support. [14]

[1]: https://help.sumologic.com/Send-Data/Collect-from-Other-Data-Sources/Amazon-CloudWatch-Logs
[2]: https://www.sumologic.com/resource/white-paper/securing-the-sumo-logic-service/
[3]: https://www.sumologic.com/blog/it-operations/logs-to-metrics/
[4]: https://help.sumologic.com/Send-Data/Collector-FAQs/Delete-data-already-collected-to-Sumo-Logic
[5]: https://github.com/newrelic/go-agent
[6]: https://blog.newrelic.com/product-news/government-it-modernization/
[7]: https://blog.newrelic.com/product-news/role-based-access-control-rbac/
[8]: https://docs.newrelic.com/docs/using-new-relic/new-relic-security/security/security
[9]: https://docs.newrelic.com/docs/insights/use-insights-ui/manage-account-data/editing-deleting-insights-data
[10]: https://docs.newrelic.com/docs/insights/nrql-new-relic-query-language/nrql-resources/nrql-syntax-components-functions
[11]: https://github.com/honeycombio/beeline-go
[12]: https://www.honeycomb.io/security/
[13]: https://docs.honeycomb.io/authentication-and-security/secure-tenancy/
[14]: https://docs.honeycomb.io/getting-data-in/datasets/secure-manage/
