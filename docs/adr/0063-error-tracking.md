# *Establish a system for front end error tracking*

*We should have a way to know if front end errors are happening. Many SAAS products
need to be FedRAMPed in order to conform with ATO standards. Here we've
investigated a few different SAAS products*

*Sentry and Rollbar started off as
error tracking tools but seem to have upped their offerings to include other types
of logging. Datadog and Dynatrace started as logging tools but have recently
increased offerings to include error tracking. It would be nice to use the same tracking/logging for FE and BE*

## A couple of questions

* *Is there room in the budget for a SAAS*
* *If there is space in the budget given the pricing, is there more interested in a fully built out monitoring system
or should we focus solely on error tracking*
* *Does our ATO allow us to use any previously FedRAMPed tools in our system*
* *Should we even consider any non-FedRAMPed cloud solutions?*


## Considered Alternatives

* *Sentry*
* *Rollbar*
* *Datadog*
* *Dynatrace*
* *Create our own endpoint*

## Decision Outcome

* Chosen Alternative: *TBD*

## Pros and Cons of the Alternatives <!-- optional -->

### *Sentry*

* `+` *Has a self-hosted option*
* `+` *Was used on a previous truss project, so we have some in house experience*
* `+` *Specialized tool for error tracking*
* `-` *It would take time/resources for infra to setup/maintain a self hosted version*
* `-` *[Monthly fee](https://sentry.io/pricing/)*

### *Rollbar*

* `+` *[May be willing](https://rollbar.com/blog/blog/introducing-hassle-free-compliant-saas-error-monitoring) to get additional security certifications they need to work with clients*
* `-` *[Monthly fee](https://rollbar.com/pricing/)*

### *Datadog*

* `+` *FedRAMP In Process*
* `+` *Full monitoring suite*
* `-` *[Pricing](https://rollbar.com/pricing/) based on data ingested - if we aren't careful may accidentally log too much*

### *Dynatrace*

* `+` *FedRAMP Authorized*
* `+` *Full monitoring suite*
* `-` *[Monthly fee](https://www.dynatrace.com/pricing/)*

### *Create our own endpoint*

* `+` *We'd have full control of all logged data*
* `+` *No monthly fee once it's up and running*
* `+` *Wouldn't have to deal with FedRAMP*
* `-` *Would need to figure out a friendly way to look at/analyze data*
* `-` *Not sure how it would scale*

