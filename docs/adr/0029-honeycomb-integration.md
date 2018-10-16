# HoneyComb Integration

* Status: accepted
* Deciders: @elliottdotgov, @hlieberman-gov, @ntwyman
* Date: 2018-10-01

## Context and Problem Statement

Currently, the best method for tracing a request through the MyMove site is done by inspecting the JSON formatted logs hosted in AWS CloudWatch Logs. The CloudWatch Logs service provides the ability to define point queries based on the fields in the JSON log entries.  For example, a query could be looking up all log entries with a severity level of “error”. Simple searches like looking up errors do provide some value, particularly when looking for known issues and if the volume of requests flowing through the site is low. However, CloudWatch Logs doesn’t provide engineers a good way to look into more abstract unknown problems that could be happening.

It’s important to be able to ask high level and ad-hoc questions like “what are the top 5 slowest API calls right now?” . More importantly, being able to ask follow up questions like “For those slow API calls, where is that time being spent?”.

This decision records is meant to pick a path forward to increase the visibility on the state of the MyMove site.

## Decision Drivers

### Operational Requirements

* **Outsource to a Software as a Service (SaaS) provider** - One way to solve these problems could be to build a better system in-house (e.g., deploying [Elasticsearch](https://github.com/elastic/elasticsearch ), [Logstash](https://github.com/elastic/logstash), and [Kibana](https://github.com/elastic/kibana)), but the operational overhead needed to maintain such a system is significant and isn’t a primary focus for MyMove.

### Usability Requirements

* **Be able to ask previously unknown questions** - Provides a clear interface for asking both granular and coarse question about the state of the MyMove site.
* **The service should be easy to use** - Doesn’t require learning a domain specific language(DSL) or deep understanding of regular expressions.

### Security Requirements

* **Encrypt all data sent from MyMove** - Service provides encryption at rest and in-transit.
* **Allow MyMove infrastructure to control access to encrypted data before being sent** - MyMove remains in control of how the data is encrypted or hashed before leaving the infrastructure.
* **Provide tools to control data how data is persisted** - Strict controls on how data is persisted and has controls for deleting columns or narrowed subsets.
* **Have good access control policies for accessing or deleting data** - Role based access for accessing and more importantly deleting data.

## Considered Options

* **Log aggregation service** - [Sumo Logic](https://www.sumologic.com/)
* **Application Performance Monitoring (APM) service** - [New Relic](https://newrelic.com/)
* **Event based observability service** - [HoneyComb](https://www.honeycomb.io/)

## Decision Outcome

Honeycomb straddles a difficult balance that provides a pretty good user interface for debugging complex datasets, but also provides a mechanism to encrypt any data we send before it leaves the MyMove infrastructure. Below is an example of using Honeycomb to see the 99th percentile request latencies (in milliseconds) for the slowest API calls. We’ve been testing the service in the MyMove staging environment to verify its usefulness before moving forward with a roll out to production. We did this to understand what data we will and won’t send to Honeycomb. Also it allowed us to understand what sorts of security controls are in place to protect and manage the data we do send.

## Pros and Cons of the Options

### [Sumo Logic](https://www.sumologic.com/)

* `+` Integrates with CloudWatch Logs
* `+` Encryption in-transit and at rest using keys maintained and rotated by Sumo Logic
* `+` Fine grained role based access controls.
* `-` Doesn’t support encrypting data before leaving MyMove infrastructure.
* `-` Requires managing search indexes and deciding which data is important to index ahead of time.
* `-` Querying JSON logs requires using JSONPath syntax, which isn’t much better than CloudWatch
* `-` Generating graphs around things like requests latencies requires defining dashboards and configuration ahead of time using complicated “Parse Expressions”.
* `-` Very Limited controls around data persistence. All deletion requests must go through support.

### New Relic

* `+` Provides good library for integrating with Go codebase.
* `+` In the process of being FedRAMP approved.
* `+` Provides defining both coarse and fine grained queries.
* `+` Fine grained role based access controls.
* `-` Data is encrypted in-transit, but not at rest.
* `-` Data cannot be deleted once written without an explicit support request.
* `-` Doesn’t support encrypting data before leaving MyMove infrastructure.
* `-` Difficult to extend beyond the quick start functionality.
* `-` Querying the New Relic Events services requires learning New Relic Query Language (NRQL).

### Honeycomb

* `+` Provides good library for integrating with Go codebase.
* `+` Data is encrypted in-transit and at rest.
* `+` Supports encrypting data before leaving MyMove infrastructure using Secure Tenancy.
* `+` Web based UI for querying against any number of fields in the Honeycomb dataset.
* `+` No complex query syntax to learn.
* `-` Limited role based access controls.
* `-` Secure Tenancy requires running new infrastructure in-house.
