# *Use Office application for Prime UI*

- **Epic Story:** *[MB-8515][jira-epic]*
  - **User Story:** *[MB-8575][jira-chore]*

[jira-chore]: https://dp3.atlassian.net/browse/MB-8575
[jira-epic]: https://dp3.atlassian.net/browse/MB-8515

## Background

Currently, neither the Government nor the contractor is able to fully test to
use of the external GHC Prime API without the help of engineers who are already
working on the project. This means that the Government is unable to do
acceptance of the work being delivered end-to-end and the non-technical members
of the contracting team are not able to test internally or with users.

### Leveraging the Office app

Leveraging the Office application gives Truss engineers one less application to
maintain. Users of MilMove applications will interact with the same visual user
interface (UI) that is used to with other MilMove applications. The Office
application is leveraged by using a `Prime Simulator Role` for the Office user.
This is set via the Admin UI application, covered in [Use Query Builder for
Admin Interface](035-us-query-builder.md), by assigning the role to Office
users.  Onboarding to Prime UI use can be done by non-engineers. Accessing the
Prime UI will done without the aid or use of engineering effort.

## Considered Alternatives

- *Leverage the Office Application to have a Prime Simulator Role*

## Decision Outcome

- Chosen Alternative: *Leverage the Office Application to have a Prime Simulator Role*

### *Leverage the Office Application to have a Prime Simulator Role*

- `+` *Any user is able to able to validate the work that has been completed by the
    contractor (Truss), without requiring the use of engineering effort.*

- `+` *Any user is able to test the system, end-to-end, internally, without requiring
    the use of engineering effort.*

- `+` *Any user is able to demo the system, end-to-end, internally, without requiring
    the use of engineering effort.*

- `+` *Trussel engineers will be exposed more to Prime API functionality and share
    knowledge of the system, end-to-end, internally, without requiring the use
    of specialized Prime API engineers.*

- `-` *There are some risks to this approach which are covered under **Decision
    Risks**.*

## Decision Risks

This section of the ADR is new addition as this particular ADR does not have any
other considered alternatives. This section may appear in future ADRs that
follow a similar process. This decision has some specific risks involved related
to its technical implementation. Specific guardrails must exist, otherwise this
decision's risks will become a problem for the maintenance of this application.

### Roles restricted to non-production environments

Due to the nature of the Prime API contract, the Prime UI is an entirely
internal application that will not be available in production environments. A
secure migration for the application must be run in the production database
which is a non-operational (NOOP) migration. This is achieved by having an empty
migration file with SQL comments. Please read the [documentation on Secure
Migrations][docusaurus] by searching the Docusaurus site. *A direct
link is not included here because the documentation site is going a
restructuring while this ADR is being written*. This type of migration should be
used for any environments which may be *production-like* or where the `Prime
Simulator Role` must not exist.

[docusaurus]: https://transcom.github.io/mymove-docs/

### USWDS components are used where possible

In order to get the benefits of using the Office app, the Prime UI must leverage
React components of the Office application or the United States Web Design
System (USWDS) which MilMove builds on as a foundation for the design system for
MilMove applications. Portions of the Prime UI application must have a similar
look and feel as the rest of the MilMove Office application. This helps users
have a unified visual language when interacting with the Prime UI. Engineers
contributing to the Prime UI will have a shared understanding of how to interact
with the Prime API. This leads to shared knowledge of the Prime API and Office
application across different engineering teams and practices.

For more clarity, Truss maintains the `React-USWDS` component library that is
used in MilMove applications. The USWDS is purely a CSS library and while we do
import it directly for some things, it's not the foundation of the application.

#### CODEOWNERS design reviews will be optional

Designers will not need to be required to review any changes made
to the Prime UI application. Collaboration between design and engineering for
the Prime UI is encouraged but not required. This is enforced by the CODEOWNERS
file having no reviewers for `src/*/PrimeUI/` directories.

#### Handling Business Logic

Future decisions will need to be made around how to handle Business Logic that
Users may expect to do in a single action that our RESTful API endpoints are not
capable of doing. This means that either the Support API or Client-side data
manipulation will have additional work to update or tie data behind the scenes.
