# Use snapshot to cleanup load testing

**User Story:** [Jira Story](https://dp3.atlassian.net/browse/MB-3071)

We use the tool Locust to perform load testing. We will be running this tool locally and it will hit endpoints on the experimental deployment of Milmove.

We need a solution to cleanup the database after load testing. Load testing creates a large number of moves and associated records such as orders, users, shipments and payment requests.

Once they are in the database, there's no obvious way to differentiate them from other valid records and clean them up.

We need a strategy for cleanup.

## Considered Alternatives

* Create and restore snapshot of the database.
* Create a batch delete endpoint in Support API. Delete script run by load testing machine.
* Upload delete UUIDs to AWS, create an AWS scheduled task to delete all.

## Decision Outcome

We will snapshot the database on experimental prior to load testing and restore it after.

As we are currently only planning on running against experimental, this is the cleaner, simpler way to solve for this. The other options require too much initial code + maintenance, and would only become necessary if we decided to run against staging.

## Pros and Cons of the Alternatives

### Option 1: Create and Restore Snapshot

Currently we are only targeting load testing against experimental. Therefore, we are able to snapshot the database and then restore it at the end of the load test. If possible, we will create a script that helps with the creation and restoration of this snapshot for ease of use by developer.

* `+` This will cleanly remove all load testing artifacts and is not subject to us coding anything to specifically find and delete the records.
* `+` This mechanism is currently used by infra and is a known process
* `-` Need to get permissions for app-eng to be able to complete this.
* `-` **Can only be used on experimental**. We can only do this in experimental or another private testing environment. Not suitable for staging where multiple people may be using it at the same time.

### Option 2: Delete batch endpoint in Support API. Delete script run by load testing machine

We could create a delete endpoint in the Support Api.
This could take batch lists of user UUIDs as the parents, and delete all users, orders, moves etc. associated.

We can't clean up the objects as we load test because that would defeat the purpose of load testing.

* `+` Could be used on staging and experimental.
* `-` We would have to log all created user UUIDs locally while we are running testing and after completion delete all of them. This opens up the possibility that the logs get deleted or corrupted, or the developer forgets to run the cleanup script. Once that information is lost, we can't reconstruct it.
* `-` Extra work needed to create endpoint and keep it updated as structure of nested objects changes.

### Option 3: Upload delete UUIDs to AWS, AWS scheduled task to delete all

We could configure load testing to periodically upload the list of user UUIDs to delete to AWS. We would need an endpoint to do this. In addition, we could use an AWS scheduled task that processes them all at end of day directly deleting from the DB.

We can't clean up the objects as we load test because that would defeat the purpose of load testing.

* `+` Could be used on staging and experimental.
* `-` Can’t run load testing overnight if task will cleanup objects on a set schedule.
* `-` Extra work needed to upload ids, create endpoint for uploads, create task to process, and keep it updated as the structure of nested objects changes.

