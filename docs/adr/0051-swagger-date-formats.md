# Use only Swagger supported formats for dates

In our Swagger yaml files we should only be using date formats that are supported by Swagger.

## Considered Alternatives

* Leave everything as is. Continue to use unsupported date formats, namely, `datetime`

## Decision Outcome

* Chosen Alternative: **Use Swagger supported date formats, `date-time` or `date`, depending on whether we need to store an exact timestamp of the event.**

This option will assure that we are using Swagger supported data types

## Pros and Cons of the Alternatives

### Leave everything as is. Continue to use unsupported date formats, namely `datetime`

* `+` No changes needed to be done.
* `-` Leads to inconsistent data type usage for what should be similar data.
* `-` Using a data type format that is not supported by Swagger.

### Use Swagger supported date formats, `date-time` or `date`, depending on whether we need to store an exact timestamp of the event

* `+` Makes correct use of Swagger data types.
* `+` Maintains consistency in how we format dates in our APIs.
* `-` Requires changes to yaml files.
