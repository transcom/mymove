# Move duplicated service member fields to orders

*[Jira Epic](https://dp3.atlassian.net/browse/MB-3375)*

When the `service_members` table was first created, it included a `rank` field. Not too long after, the `duty_station` field was added. Then about a year and a half later, the `move_orders` table was introduced and it included the `grade`, `origin_duty_station_id`, and `destination_duty_station_id` fields. The idea behind having those fields in the `move_orders` table (which has since been consolidated into the `orders` table) was to capture the service member's grade (aka rank) and where they moved from and to at the time of the order. Given that a service member can have multiple orders throughout their service, having that history is useful.

When a service member goes through onboarding, they enter their rank and origin duty station as part of the Profile portion of the flow, before the order is created. This stores the `rank` and `duty_station` on the `service_members` table.

When a TOO reviews a move, they need to see the service member's current rank and origin duty station. The TOO interface is currently fetching this information from the `orders` table. The way the `orders` table gets this information is by copying it from the `service_members` table into the
`orders` table.

The problem with storing that data in two tables is that we have to worry about keeping them in sync. For example, if a service member made a mistake when entering their origin duty station, they can update it (which will update it in the `service_members` table), but if the order was already created, it needs to be updated as well. We don't currently handle this scenario in the code.

To resolve data integrity issues, we propose removing those fields from the `service_members` table and storing them on the `orders` table only.

## Considered Alternatives

* Keep the fields in both `service_members` and `orders`
* Remove the fields from `service_members` and only store them in `orders`

## Decision Outcome

* Chosen Alternative: Remove the fields from `service_members` and only store them in `orders`

## Pros and Cons of the Alternatives

### Keep the fields in both `service_members` and `orders`

* `+` Does not require redesigning the onboarding flow
* `-` Results in data integrity issues
* `-` Increases code complexity when writing to two tables and worrying about how and when to keep them in sync

### Remove the fields from `service_members` and only store them in `orders`

* `+` Eliminates data integrity issues
* `+` Simplifies code because we are only writing to one table
* `-` Requires a redesign of the onboarding flow so that rank and origin duty station are captured during the order portion.
