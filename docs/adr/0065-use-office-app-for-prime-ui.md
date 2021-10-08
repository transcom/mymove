# *Use Office application for Prime UI*

- **Epic Story:** *[MB-8515][jira-epic]*
  - **User Story:** *[MB-8575][jira-chore]*

[jira-chore]: https://dp3.atlassian.net/browse/MB-8575
[jira-epic]: https://dp3.atlassian.net/browse/MB-8515

## Background

Currently, neither the Government nor the contractor is able to fully test to
use of the external GHC Prime API without the help of engineers who are already
working on the project. This means that the Government is not able to do
acceptance of the work being delivered end-to-end and the non-technical member
of the contracting team are not able to test internally or with users.

### Leveraging the Office app

Leveraging the Office application gives Truss engineers one application less to
maintain. Users of MilMove applications will interact with the same visual user
interface (UI) that they are used to with other MilMove applications. It is
leveraged by using a `Prime Simulator Role` for the Office user. This is set via
the Admin UI application, covered in [Use Query Builder for Admin
Interface](035-us-query-builder.md), by assigning the role to Office users.
Onboarding to Prime UI use can be done by non-engineers. Accessing the Prime UI
will done without the aid or use of engineers.

## Considered Alternatives

- *Leverage the Office Application to have a Prime Simulator Role*

## Decision Outcome

- Chosen Alternative: *Leverage the Office Application to have a Prime Simulator Role*

### *Leverage the Office Application to have a Prime Simulator Role*

- `+` *Any user is able to able to validate the work that has been completed by the
    contractor (Truss), without requiring the use of engineers.*

- `+` *Any user is able to test the system, end-to-end, internally, without requiring
    the use of engineers.*

- `+` *Any user is able to demo the system, end-to-end, internally, without requiring
    the use of engineers.*

- `+` *Trussel engineers will be exposed more to Prime API functionality and share
    knowledge of the system, end-to-end, internally, without requiring the use
    of specialized Prime API engineers.*

- `-` *There are some risks to this approach which are covered under __Decision
    Risks__.*

## Decision Risks

This section of the ADR is new addition as this particular ADR does not have any
other considered alternatives. This section may appear in future ADRs that
follow a similar process. This decision has some specific risks involved related
to its technical implementation. Specific guardrails must exist otherwise this
decision's risks will become a problem for the maintenance of this application.

### Roles restricted to non-production environments

Due to the nature of the Prime API contract, the Prime UI is an entirely
internal application that will not be available in production environments. A
secure migration for the application must be run in the production database
which is a non-operational (NOOP) migration. This type of migration should be
used for any environments which may be production-like or where the `Prime
Simulator Role` must not exist.

### USWDS components are used where possible

In order to get the benefits of using the Office app, the Prime UI must leverage
React components of the Office application or the United States Web Design
System (USWDS) which MilMove builds as a foundation. Portions of the application
will have a similar look and feel as the rest of the MilMove Office application.
This helps users have a unified visual language when interacting with the Prime
UI. Engineers contributing to the Prime UI will have a shared understanding of
how to interact with the Prime API. This leads to having knowledge of the Prime
API and the Office application shared across different engineering teams and
practices.

#### CODEOWNERS design reviews will be optional

Due to the reuse of React components, visual design of the application will not
be necessary as components will be reused from their approved designs.
Therefore, the designers will not need to be required to review any changes made
to the Prime UI application. Collaboration between design and engineering for
the Prime UI is encouraged but not required.

#### Handling Business Logic

Future decisions will need to be made around how to handle Business Logic that
the Users may expect to do in a single action that our RESTful API endpoints
are not capable of doing. This means that either the Support API or Client-side
data manipulation will have additional work to update or tie data behind the
scenes.
