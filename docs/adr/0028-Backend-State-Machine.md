# Backend Data State Machine

**User Story**: [Plan out and author new state machine approach](https://www.pivotaltracker.com/story/show/159704281)

Our current implementation of state machine throughout the app is haphazard. Certain models have a relatively complete set of states, well-defined transitions, and flows; others have a few states but not a complete set, or do not have the ability for us to move between states. We are aiming to create a cohesive vision of how state machine is applied so that in the future it is clear where we expect to find them in use, and if we need to implement a new one, it is clear how to do so.
The decision here is twofold:

1. How much of the system do we want covered with a state machine?
2. What kind of state machine mechanism do we want to use?

**To address the first**:
We did consider how using state machines creates a level of rigidity within the system. We don't want to force folks to use the system a certain way not because it is the best way to get to their goal but because it is the only way to do something. However, we also don't want our system to allow for models that fall out of flows in unusual circumstances. Therefore, we felt that applying state machines to certain models but not others made the most sense.
Please refer to this [Real Time Board](https://realtimeboard.com/app/board/o9J_kzRjJ2k=/) to leave comments and access the current model of state throughout the system.

## Considered Alternatives for state machine mechanism

* *Keep using our own state machines*
* *Use a state machine library*
  * *[Loop Lab Finite State Machine for Go](https://github.com/looplab/fsm)*
  * *[QOR/transition](https://github.com/qor/transition)*
  * *[A simple finite state machine for Golang](https://github.com/ryanfaerman/fsm)*

## Decision Outcome

* Chosen Alternative: Maintain out own state machine mechanism.
* The main decision driver here was it that although folks seemed to favor relying on a pre-built
* *[consequences. e.g., negative impact on quality attribute, follow-up decisions required, ...]* <!-- optional -->

## Pros and Cons of the Alternatives

### *Use a State Machine Library*

#### Decided on [Loop Lab Finite State Machine for Go](https://github.com/looplab/fsm)--most robust, highly recommended

* `+` One place in each model where all implementation occurs.
* `+` Provides hooks for enter/exit states, regardless of transition
* `+` Provides hooks for code to run before/after transitions (before transition can block a transition)
* `-` Have to implement everywhere we are using current state machines (includes adding a new field with serialization/de-serialization functions to each model that uses).
* `-` Have to change currently existing Status field to be some form of private because we wouldn't want that to be writable outside of the FSM paradigm. Note, making a private field isn't exactly possible in this situation, so we'd actually probably implement through an awful name like "state_DO_NOT_EDIT".
* `-` Probably heavier than we really need, since our transitions tend to be quite straightforward
* `-` Additional arguments to events are of type interface{}

* **To determine the difficulty of implementation, [we stubbed out implementing on the Move model here](https://github.com/transcom/mymove/tree/rek_fsm-loop). Please use this to form your opinion of the size of this task.**

### *Keep using our own version*

* `+` Have already defined the majority of necessary transitions.
* `+` Only write what you need
* `+` Additional arguments to events are of a specific type since the code does not need to be generic.
* `-` Doesn't save when you fire transitions
* `-` Each transition/state function has to be written independently
* `-` Legal transitions are not defined declaratively, possibly making the state graph harder to reason about

Comparison of implementation

| Feature                                                                                          | Loop Lab                                                         | Current System                                                                                                                                                        |
|--------------------------------------------------------------------------------------------------|------------------------------------------------------------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| Definition of state machine                                                                      | Abstracts out which states exist,  which transitions are allowed | Defines functions for each transition, has transitions return errors if transition is not allowed                                                                     |
| How to handle code that runs before/after transitions (before transition can block a transition) | Provides Hooks                                                   | Has a transition function You can return errors for  whatever reason you want a transition to fail,  and can do whatever else you want in the body of the transition. |
| How to handle enter/exit states                                                                  | Provides Hooks                                                   | Doesn’t have this simply Could be added, might get clunky                                                                                                             |