# Add a child table to mto_shipments for PPMs

**User Story:** [MB-11140](https://dp3.atlassian.net/browse/MB-11140)

Currently, there is a `personally_procurred_moves` table that is used for PPMs while there is a separate `mto_shipments` table that is used for HHG and other shipments types. Since the existing table doesn't follow the current conventions and we are working on this outcome it is a good time to revisit the DB structure to see if the original design continues to be the best choice moving forward.

## Decision Drivers

* Time and complexity
  * Do we have bandwidth to make these changes?
  * Do we have enough capacity handle the unknowns or deal with future issues this decision may cause?
  * Do we need more information?
    * Do we have enough future capacity to delay this?
* Flexibility
  * How likely is this solution to support new shipment types or even existing shipments?
* Consistency
  * How consistent are our design patterns?
    * Are they intuitive?

## Considered Alternatives

* [Use old table](#use-old-table)
* [Modify `mto_shipments`](#modify-mto_shipments)
* [Create a new table](#create-a-new-table)

## Decision Outcome

* Chosen Alternative: [Create a new table](#create-a-new-table)
* Positive Outcomes: We will have a pattern that is easier to understand and will potentially be helpful in the future when implementing new shipments. The tentative plan is to test this new pattern out with PPMs, and if things go smoothly, revisit this to be the default DB pattern for all new shipment types moving forward.
* Consequences: If this new pattern is not utilized elsewhere we may be adding more complexity by introducing another pattern that is only partially used.

## Pros and Cons of the Alternatives

### Use old table

Continue using the `personally_procurred_moves` table

* `+` Less work, most of the plumbing is already set up and can be leveraged to account for modifications made.
* `-` Additional work would be needed to support the office workflows and to update the code to current standards.
* `-` The PPM code will remain siloed
* `-` The nomenclature is not clear. We currently have moves which contain PPMs which can cause some confusions, as we normally consider the relationship to be moves contain shipments.
* `-` The PPM models and relationship will be unique so there will be additional overhead to learn this relationship when switching to PPMs from another area of the project

### Modify `mto_shipments`

Incorporate the columns on the `personally_procurred_moves` table into the `mto_shipments` table

* `+` This will make PPMs consistent with how we store and retrieve other shipment types
* `+` A `CHECK` constraint can be used to ensure only PPM related fields are used with a PPM shipment type
* `-` The `mto_shipments` table has a growing number of fields and it’s hard to know which fields are applicable for a given shipment type
* `-` As we continue to support more types of shipments, this is beginning to feel like a workaround solution rather than the best solution
* `-` This is a sizable amount of work, and there is some risk of the amount of work ballooning into more work that originally expecting.

### Create a new table[¹](#references)

Create a new table that is a child of the `mto_shipments` table that holds ppm specific info

* `++` This structure is more flexible and can be implemented retroactively for old shipment types and support new shipment types
* `+` This is semantically the most clear and intuitive
* `-` If the existing patterns for shipments are not modified, we will continue to have multiple patterns for different types of shipments
* `-` The most amount of work and the most unknowns (slightly more than [modifying mto_shipments](#modify-mto_shipments))

## References

1. [^](#create-a-new-table) [Class table inheritance](https://martinfowler.com/eaaCatalog/classTableInheritance.html)
2. [Slack thread about this conversation](https://ustcdp3.slack.com/archives/CP6PTUPQF/p1641937842037600)
3. [Backend Discovery](https://dp3.atlassian.net/wiki/spaces/MT/pages/1605238794/PPM+Bookings+Technical+Discovery#Database)
