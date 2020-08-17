# Consolidate moves and move task orders into one database table

*[Jira Epic](https://dp3.atlassian.net/browse/MB-3021)*

Currently, a `moves` record is distinct from a `move_task_orders` (MTO) record, and there
is no direct association between the two. Data entered by customers during move onboarding
must flow into the MTO correctly. If we don't get ahead of this now, the teams will
potentially implement different versions of the database, resulting in rework and
confusion later. Having a consistent and uniform data model will improve
collaboration and productivity, and will enable us to demonstrate end-to-end
capability.

## Considered Alternatives

* Keep `moves` and `move_task_orders` in separate database tables
* Consolidate `moves` and `move_task_orders` into a single table

## Decision Outcome

* Chosen Alternative: Consolidate `moves` and `move_task_orders` into a single table
* Keep only the `moves` table and add fields to it to enable accurate
representation of an MTO.
* Below is the final state of the `moves` table. The "new" label means it was ported over from
the `move_task_orders` table.  These new fields will need to be nullable for
backwards-compatibility with the current production process.
  * `id`
  * `created_at`
  * `updated_at`
  * `orders_id`
  * `selected_move_type`
  * `status`
  * `locator`
  * `cancel_reason`
  * `show`
  * `contractor_id` (new)
  * `available_to_prime_at` (new)
  * `ppm_type` (new)
  * `ppm_estimated_weight` (new)
  * `reference_id` (new)

### Definitions of fields

* `id`: A UUID primary key

* `created_at`, `updated_at`: The usual Pop timestamps

* `orders_id`: A foreign key to the associated `orders` record

* `selected_move_type`: Allowed values are HHG, PPM, UB, POV, NTS, HHG_PPM (but
only HHG and PPM appear to be used currently).

* `status`: Allowed values are DRAFT, SUBMITTED, APPROVED, CANCELED

* `locator`: This is a 6-digit alphanumeric value that is a sharable,
human-readable identifier for a move (so it could be disclosed to support staff,
for instance). The MTO's `reference_id` is similar in nature but is in a `dddd-dddd`
format. The `reference_id` does serve currently as the prefix for payment request
numbers. We likely don’t need both `locator` and `reference_id` if these tables
merge. See [Slack discussion](https://ustcdp3.slack.com/archives/CP6PTUPQF/p1595605700223400).

* `cancel_reason`: A string to explain why a move was canceled.

* `show`: A boolean that allows admin users to prevent a move from showing up
in the TxO queue. This came out of a HackerOne engagement where hundreds of fake
moves were created. This defaults to `true`.

* `contractor_id`: This was added to represent the prime contractor who will
handle the move. This makes it easy to point the move to a different contractor in case it changes.

* `available_to_prime_at`: a date and time field that indicates when the move is
available for the Prime to handle. The presence of this field can be used to
determine whether or not to display the move to the Prime.

* `ppm_type`: currently used values are `FULL` and `PARTIAL`. This appears to be
different from having the `selected_move_type` in the `moves` table be `PPM` vs
`HHG_PPM`. See [Slack discussion](https://ustcdp3.slack.com/archives/CP6PTUPQF/p1595617833232800).

* `ppm_estimated_weight`: this is being set by the Prime currently so we are
keeping it for now.

* `reference_id`: A unique identifier for an MTO (which also serves as the prefix
for payment request numbers) in `dddd-dddd` format. There is still an ongoing
discussion as to whether or not we need this `reference_id` in addition to the
unique `locator` identifier, so we are keeping `reference_id` for now.

There will be future work to reconcile the `ppm_type` and `ppm_estimated_weight`
fields when we reenter conversations with the prime.

### Fields that we are not moving from move_task_orders to moves

* `is_canceled`: used to determine if an MTO was canceled or not. The moves
table already has a `status` field with a `CANCELED` option, so we can get rid of
`is_canceled` and use `status` instead.


## Pros and Cons of the Alternatives

### Keep moves and move_task_orders separate

* `+` Allows for the possibility of multiple MTOs for a single move.
* `-` Increase the risk of code complexity and data duplication.
* `-` Makes it more difficult to represent the move from the point of view of
all parties: service member, TOO, Prime. For example, to find a move related to an MTO, you have to find the `order_id` that the MTO points to, then the move that points to that same `orders_id`
(unless we add a foreign key to a `moves` record from a `move_task_orders` record)

### Consolidate moves and move_task_orders

* `+` It keeps things simple in the codebase because an MTO is essentially a
move that is available to the Prime. All moves will require an MTO except in one
specific scenario: when the service member chooses to handle the move on their
own (PPM) AND they receive counseling from services and not from the Prime.
* `+` A move can only have one MTO, and the information the MTO refers to also
applies to a move, so it makes sense to only have one DB table.
* `-` Because of the differences between both tables, the consolidation will
require more effort and might cause breaking API changes.
* `-` By not having a separate DB table for MTOs, there is a risk we might not
be representing an MTO accurately. An MTO is a legal construct with specific
requirements.

## Definitions

**TOO - Task Ordering Officer**. They are responsible for generating Task Orders
and "ordering" shipments and service items such as crating, shuttle service and
SIT. They are also a check on lines of accounting to make sure the correct ones
are on the MTO so the Prime is paid from the appropriate bucket of money.

**PPM - Personally Procured Move**. When a service member chooses to handle the
move on their own.

**MTO - Move Task Order**. Is similar to an order for goods from a contractor.
In the case of MilMove, the TOO is ordering services from the Prime Contractor.
When the Prime Contractor completes those services, they can request payment for
those services. Every service the Prime undertakes must be "ordered." The
government does this via a Move Task Order. It is the record of everything that
is ordered (approved) for the Prime to do. The Move Task Order contains all the
information about shipments, including approved service items, estimated weights,
actuals, requested and scheduled move dates etc.

