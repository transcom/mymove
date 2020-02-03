# *Handling time in the Prime API*

## Definitions

Date

* A time value that only includes year, month, and day
* `DATE` in postgres

Timestamp

* A more granular time value that also includes hours, minutes, seconds, and smaller units
* The smallest unit supported by timestamps varies by system; OS, database, and programming language can differ
* `TIMESTAMP` in postgres

## Background

Our system currently works with time values at two levels of resolution: date and timestamp. Date values typically result from an explicit day selection in the UI (“what day do you want your shipment picked up?”), while timestamps are typically generated implicitly by the system to record when a certain action occurred.

The problem arises when we need to perform operations that mix these two types. For example, try to answer this question: Is `2020-01-20T16:00:00` more than one day after `2020-01-19`? Answering this requires that we convert one of these values to the other’s resolution; we can either convert `2020-01-20T16:00:00` to a date (perhaps by truncating its hours, minutes, and seconds) or convert `2020-01-19` to a timestamp by “adding” precision and giving it values for the more granular subfields.

In order to ensure that we handle operations that involve multiple levels of time resolution consistently throughout the system, we need to have clear guidelines for how to do so.

An alternative to this is to prevent the system from storing temporal values in two different types. Doing this will just move the complexity around the system, though, since the backend only working with timestamps would then require the UI and API consumers to define and
implement their own conversions between the two date types.

## Proposed Solution

* Use timestamps when recording when an action has occurred. Examples include created_at, updated_at, etc. values

* Use dates when accepting a scheduled date from the Prime. Examples include scheduled move date, scheduled pickup date, etc. Dates are typically formatted as `YYYY-MM-DD` strings in payloads

### Convert dates to timestamps when comparing with timestamps

1. Interpret a date in Pacific time to convert it to a naive timestamp

2. Use beginning-of-day (12:00:00 am) or end-of-day (11:59:59 pm) depending on operation being performed. In general, use the most forgiving interpretation:
    * If you’re calculating if a date is more than 10 days in the future, use end-of-day.
    * If you’re calculating if a date is more than 2 days ago, use beginning-of-day.

### Serialization

Use [RFC3339](https://tools.ietf.org/html/rfc3339) unless you have a reason not to. It handles both dates and timestamps. `ISO8601` is roughly equivalent for the purposes of this document; however, it contains many optional features and it is recommended that projects move to `RFC3339` to avoid potential issues resulting from this complexity.

* Date example: `2007-11-13T00:00:00Z` (Note that this includes zeroed out hours, minutes, etc.)
* Timestamp example: `2007-11-13T09:13:00Z`

## Considered Alternatives

* Store all time values as timestamps
* Record a timezone in the database for all dates

## Decision Outcome

### Chosen Alternative: Use proposed rules to convert dates to timestamps as needed

* `+` This allows us to build on top of our existing database and its data, which is already mixing date and timestamps.
* `+` We can start with a simple approach of interpreting all dates in pacific time and layer on a more timezone-aware approach later as required by business rules.
* `-` Logic that uses values at both levels of resolution will need to follow the guidelines defined here, which due to its complexity is more likely to be lead to bugs. We can mitigate this by defining and using some well-tested and documented helpers.

## Pros and Cons of the Alternatives <!-- optional -->

### Store all time values as timestamps

* `+` This makes comparisons trivial as all time values are at the same level of granularity
* `+` Moving logic from the database to the application is trivial, since there is no involved logic to translate.
* `-` We already have a lot of date values we'd need to migrate to this approach
* `-` Would require the semantics of data to change, e.g. we would store a time that was the end of the selected day for requested pickup date. That value could then be used to determine if a pickup actually happened before that time.
* `-` We would need to be able to translate from day selected in the UI to timestamps in the database, probably using either the user's local time or a timezone determined from the geographic location of the relevant event.

### Record a timezone in the database for all dates

* `+` We would have all the data needed to properly interpret dates in the database, which would simplify some calculations we need to do.
* `-` This adds complexity to the Prime API for consumers, as they would need to provide a timezone or location for every date.
* `-` We already have a lot of date values we'd need to migrate to this approach