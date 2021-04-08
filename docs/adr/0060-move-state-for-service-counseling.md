# _Move statuses to support service counseling_

**User Story:** [Create new move "status" for Services Counseling](https://dp3.atlassian.net/browse/MB-7112)

_We need to define states for when moves go to service counseling and when service counseling is complete_
\*it would be nice to minimize the impact of this to downstream systems.

## Considered Alternatives

- _Create two states and one new timestamp_
  - When a customer submits a move, the submitted timestamp is always set.
  - If the move is routed to service counseling, the move state is set to `needs_service_counseling`. Otherwise it is set to `submitted`.
  - When the service counselor is done, the move state is set to `service_counseling_completed` and the `service_counseling_completed_at` is set to the current time.
  - The TOO queue excludes items in `draft`, `needs_service_counseling` and shows those in `service_counseling_completed` (along with other states it currently shows).
- _create one state and one new timestamp_
  - When a customer submits a move, the submitted timestamp is always set.
  - If the move is routed to service counseling, the move state is set to `needs_service_counseling`. Otherwise it is set to `submitted`.
  - When the service counselor is done, the move state is set to `submitted` and the `service_counseling_completed_at` is set to the current time.
  - The TOO queue excludes items in `draft`, `needs_service_counseling` and shows those in `submitted` (along with other states it currently shows).

## Decision Outcome

- Chosen Alternative: _Create two states and one new timestamp_
  - This more closely matches the mental model of the user and is not significantly more work than the alternative.
  - We don't get into a confusing situation where moving to submitted updates different timestamps.
