# Use CSS to highlight unfinished features

For the MilMove demo on 8/4/2020, the decision was made to use a feature flag to conceal unfinished work, in order to reduce confusion on the part of the client as to what work was complete vs. not yet finished. The MilMove project encountered a similar issue in 2018 and decided to use a low-lift solution of highlighting unfinished work with a yellow background via a CSS class applied to UI elements. This decision resulted in increasing transparency around feature status, lowering confusion, and helped improve overall client communication around project progress.

Since client users and other folks are being set up with access to the staging environment and being instructed to report bugs, there is another use case besides demos for highlighting incomplete features. Using a CSS class would provide a clear visual cue for which features should be expected to work fully, vs. which features are still being worked on. This would help avoid the additional overhead on the part of the client as well as Truss, of "bugs" being reported, verified, triaged, as well as all the ensuing communication surrounding the issues.

Decision drivers include:

- The amount of time and cost of implementation, maintenance, removal, and surrounding cognitive overhead
- Aiding in client understanding of work status to decrease potential confusion or false expectation

Note:

- all incomplete features are already behind feature flags specific to environment, so everything discussed here is based on the assumption that neither the features in question nor the CSS class is visible in production environments.
- by "incomplete" or "unfinished" work, I am referring to the current state of how work is defined and planned across all the teams of MilMove, which seems to be a mixture of [incremental and iterative development](https://agility.im/frequent-agile-question/difference-incremental-iterative-development/). Examples include UI form element input data not yet submitting to the API, displaying a minimum amount of fields for an object rather than ALL expected fields, yet-to-be determined features (e.g. handling entitlement calculation on HHG's) and their consequent display (e.g. not showing entitlement on the HHG Review page), etc. These are a couple of examples that come to mind, but I am confident that I can rustle up multiple examples across the project, if need be.

## Considered Alternatives

- Continue as is (i.e. not highlighting unfinished features)
- Use feature flags to conceal unfinished work from the client
- Use CSS to highlight unfinished features

## Decision Outcome

- Chosen Alternative: TBD

## Pros and Cons of the Alternatives <!-- optional -->

### Continue as is (i.e. not highlighting unfinished features)

- `+` Very low effort and cost
- `-` Incomplete work will not be clearly communicated to the client when interacting with the app, potentially leading to confusion and false expectations.
- `-` Incomplete work will not be clearly communicated to the client during demos, squelching the opportunity for clarifying questions or course-correcting false assumptions for road map discussions.
- `-` For users reporting bugs on staging, there is a higher risk that known "bugs" will be reported to the MilMove team, adding to extra work and communication overhead in sorting through what are actually viable bugs.

### \* Use feature flags to conceal unfinished work from the client

- `-` Higher effort/cost than applying a CSS class - adding additional feature flags takes roughly 1.5 engineering days to implement vs. 1 hour for applying a CSS class. Additional overhead for removing flags than CSS class applications.
- `-` More time/effort to communicate how to set feature flags to the entire team, as well as to the client when interacting with the staging environment.
- `-` A feature flag only hides features, which does not aid in clarifying what has been started and yet to be done. If the client cannot see which features have been started, there are no opportunities to ask clarifying question or help establish appropriate expectations.

### Use CSS to highlight unfinished features

- `+` Minimal effort/cost compared to using feature flags - usually single line additions.
- `+` Lower overhead in communication to the team and client - users don't need instructions for turning on/off flags, only that a visual cue highlights not-yet-complete work for added clarity.
- `+` Incomplete work will be clearly communicated to the client during demo and when interacting with the app, avoiding confusion and false expectations.
- `+` Allows for the opportunity for clarifying questions to be asked during demo around specific features, aiding in client understanding of work status.
- `+` For users reporting bugs on staging, there should be a lower amount of false "bugs" reported, given that it should be clear to users what is complete vs. not.
- `+` Avoids concealing work from the client. If less confusion and higher trust is the desired outcome, highlighting rather than concealing incomplete work seems to bring us closer to that goal.
- `-` Takes more time than doing nothing (i.e. continuing as is).
