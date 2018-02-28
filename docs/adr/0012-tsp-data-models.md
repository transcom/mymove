# The TSP Data Models

**User Story:** [155524224](https://www.pivotaltracker.com/story/show/155524224)

The TSP Award Queue needs to be able to access information from the database efficiently enough to award shipments to TSP based on TDL, Performance Period, Best Value Score, Quality Band, and number of shipments already awarded.
Decision drivers included the anticipated format of frequent queries, query speed, and size of data.

## Considered Alternatives

* Separate `Quality Band Assignment`, `Performance Period`, and `Best Value Score` tables.
* `Transportation Service Provider Performance` table excluding `award_count` field. This table has all relevant information from `Best Value Score`, `Quality Band Assignment`, and `Performance Period` tables. Those tables do not exist. Join `Shipment Awards` table on `TDL` to determine the number of shipments already awarded to TSP (known as award count).
* `Transportation Service Provider Performance` table including `award_count` field. This is the same `Transportation Service Provider Performance` table as described above, with the addition of an `award_count` field, obviating the need for a join between the `Shipment Awards` and TSP performance tables to determine TSP award counts. Ends up denormalizing data from `Shipment Awards` table.

## Decision Outcome

### Chosen Alternative: *`Transportation Service Provider Performance` table including `award_count` field*

* Justification: Originally, the proposal included the three separate tables detailed in alternative 1. After realizing we'd be unable to pursue this route for the reasons described below, we concluded that a table that included the BVS information and quality band information was logical, since the latter depends on the former. After making this decision, it was clear the `BVS` table no longer contained information specific to BVS. So we began considering it a table to represent `Transportation Service Provider Performance`, which is what determines the order and number of shipments awarded to each TSP.

* Consequences: Denormalization means that we have to be vigilant that the `award_count` does not get out of sync with our source of truth, the `Shipment Awards` table. We also are still not totally clear how to index our new `Transportation Service Provider Performance` table so that we can reap the benefits of indices without bloating our table with them.

## Pros and Cons of the Alternatives

### *`Quality Band Assignment`, `Performance Period`, `Best Value Score` tables*

This was the original proposal. It included the following 3 tables:

* `Best Value Scores` table with an `ID`, `TDL` and `TSP` foreign keys, and a `BVS`.
* `Quality Band Assignment` table with `ID`, `band_number`, `performance_period`, number of `shipments_per_band`, and `TSPs` and `TDL` foreign keys.
* `Performance Period` table with `ID`, `start_date`, and `end_date`.

We quickly discovered blockers we could not ignore.

* `+` Separate interests
* `+` Conceptually easy to understand
* `-` DEALBREAKER: Made impossible assumptions (for example, included a `transportation_service_provider_ids` field, which would have called for a 1:many. It is not be possible to point from the 1 quality band to many TSPs, which is the direction 1:many relationships are oriented).
* `-` Would need to JOIN quality band assignment and performance period tables first to determine quality bands and then to determine which TSPs are next in line to receive shipment awards.

### *`Transportation Service Provider Performance` table excluding `award_count` field*

* `+` Same pros as chosen alternative
* `+` Does not denormalize data by having `award_count` derived from `shipment_awards`
* `-` Requires joining `shipment_awards` to `transportation_service_provider_performance` and iterating through all the TSPs in the relevant TDL to determine the next TSP to award a shipment to. This join was thought to be prohibitively slow for a query that would be made fairly frequently.
