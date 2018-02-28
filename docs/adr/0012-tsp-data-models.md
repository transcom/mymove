# The TSP Data Models

**User Story:** *[155524224](https://www.pivotaltracker.com/story/show/155524224)*

The TSP Award Queue needs to be able to access information from the database efficiently enough to award shipments to TSP based on TDL, Performance Period, Best Value Score, Quality Band, and number of shipments already awarded.
Decision drivers included the anticipated format of frequent queries, query speed, and size of data.

## Considered Alternatives

* Separate Quality Band Assignment, Performance Period, and Best Value Score tables.
* TSP Performance table excluding award count. This table has all relevant information from Best Value Score, Quality Band, and Performance Period tables. Those tables do not exist. Join Shipment Awards table on TDL to determine the number of shipments already awarded to TSP (known as award count).
* TSP Performance table including award count. This is the same TSP Performance table as described above, but also includes an award_count field, obviating the need for Shipment Awards table join to determine TSP award counts. Essentially, denormalizing data from Shipment Awards table.

## Decision Outcome

### Chosen Alternative: *TSP Performance table including award count*

* Justification: This option allows us to eliminate separate tables for performance period and quality bands that contain almost no unique information. Originally, the proposal included the three separate tables described in alternative 1. After realizing we'd be unable to pursue this route for the reasons detailed below, we thought that a table that included the BVS information and quality band information was logical, since the latter was dependent on the former. At that point, the BVS table no longer contained information specific to BVS, at which point we began considering it a table to represent TSP performance. Performance is relevant because that determines the order amount of shipments awarded.

* Consequences: Denormalization means that we have to be vigilant that the award_count does not get out of sync with our source of truth, the Shipment Awards table. We also are still not totally clear how to index our new TSP Performance table so that we can reap the benefits of indices without bloating our table with them.

## Pros and Cons of the Alternatives

### Quality Band Assignment, Performance Period, Best Value Score tables

This was the original proposal. It included the following 3 tables:

* Best Value Scores table with an id, tdl, bvs, and tsp
* Quality Band Assignment table with id, TSPs, TDL, band number, performance period, and number of shipments per band
* Performance Period table with id, start date, and end date.

* `+` Separate interests
* `+` Conceptually easy to understand
* `-` DEALBREAKER: Made impossible assumptions (for example, included a "transportation_service_provider_ids" field, which would have called for a 1:many. It is not be possible to point from the 1 quality band to many TSPs, which is the direction 1:many relationships are oriented).
* `-` Would need to JOIN quality band assignment and performance period tables first to determine quality bands and then to determine which TSPs are next in line to receive shipment awards.

### TSP Performance table without award count

* `+` Same pros as chosen alternative
* `+` Does not denormalize data by having award_count derived from shipment awards
* `-` Requires joining shipment_awards to transportation_service_provider_performance and iterating through all the TSPs in the relevant TDL in order to determine which would be the next TSP to award a shipment to. This join was thought to be prohibitive in terms of speed for a query that would be made fairly frequently.
