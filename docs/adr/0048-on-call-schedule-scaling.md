# Milmove has an on-call schedule that scales

As the Milmove team has grown, we need to address several concerns with our on-call schedule:

- The schedule is very long and people cannot anticipate/plan for when they will be on-call.
- Team members are often on-call at the same time.
- The long gaps between times on-call mean that people aren’t able to learn the skills needed or discern patterns of failures.

## Considered Alternatives

- Each team has an on-call rotation
- Schedule is organized to minimize impact on teams, by round robin.
- Schedule is organized to ensure individuals are on-call a few times without large gaps

## Decision Outcome

- Chosen Alternative: _[alternative 1]_
- _[justification. e.g., only alternative, which meets KO criterion decision driver | which resolves force force | ... | comes out best (see below)]_
- _[consequences. e.g., negative impact on quality attribute, follow-up decisions required, ...]_ <!-- optional -->

## Pros and Cons of the Alternatives

### Each team has an on-call rotation

- `+` as teams get ownership of areas of the application, they could be alerted for those specific sub-systems
- `+` Only one person on a team would be on-call at a given time
- `+` People would be on-call more often
- `+` Consistent impact of on-call would improve teams planning
- `+` The on-call team member can address bugs and other non-feature work when not handling issues
- `-` At present, it is not clear which team should get alerted for an issue
- `-` People would be on-call every 4 weeks
- `-`More people on-call means less people focused on feature work
- `-` Not sure we really need 4 people on-call when there is low usage of the system at the moment. (But may make sense when we go live)
- `-` Not sure if the lead should be on-call or not.

### Schedule is organized to minimize impact on teams, by round robin

Primary: a 1, b 1, c 1, d 1, a 2, b 2, c 2, d 2, a 3, b 3, c 3, d 3, a 4, b 4, c 4, d 4

Secondary: b 1, c 1, d 1, a 2, b 2, c 2, d 2, a 3, b 3, c 3, d 3, a 4, b 4, c 4, d 4, a 1

Tertiary: leads

- `+` Only one person on a team would be on-call at a given time
- `+` Teams could consistently plan for impact of on-call rotation
- `-` On-call every 16 weeks, so still hard for individuals to predict their schedules
- `-`As teams focus on specific part of the project, the on-call person is less likely to have the context for issue at hand.
- `-`As soon as someone needs to swap, the impact to teams changes. (This could be mitigated by trying to keep swaps within a team.)

### Schedule is organized to ensure individuals are on-call a few times without large gaps

- `+` Only one person on a team would be on-call at a given time
- `+` People would be on-call a few times in a 8 week period, followed by an 8 week period they aren’t on-call
- `-` The complexity here is far greater than any benefit.
- `-`As soon as someone needs to swap, the main benefit of this plan is diminished
