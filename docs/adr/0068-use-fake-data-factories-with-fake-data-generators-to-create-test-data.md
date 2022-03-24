# Use fake data factories with fake data generators to create test data

## Problem statement

We have complex data models that can vary a lot depending on the situation, e.g. an HHG shipment looks very different
from a PPM shipment. This can make writing tests for our code more difficult because even if we aren't doing something
that affects every part of an object, we might still need to know the correct way to set up the data before we can even
test the part that we are interacting with. The problems with the way we do things now are a bit different between the
frontend and backend, but the chosen solution would ideally help with both.

Continuing with the shipment model example, if we want to write a test for a services counselor editing a shipment to
add counselor notes, we need a shipment that is in a state that makes sense, e.g. has all the required fields for that
type of shipment filled out. Ideally, we should be able to have something that lets us quickly get started on the tests
without having to dig into what valid values are for the other parts that we aren't interacting with directly, e.g.
orders.

To help with this on the backend, we have things like our `testdatagen` functions that will set up models in a
semi-realistic state, but they have flaws and inconsistencies. For example, you can't easily call the `testdatagen`
function for creating an address more than once without passing in `assertions` (overrides), because otherwise you'll
create the same exact address created twice. To help with that, we actually have 4 different address creation functions
so that you don't need to always pass in assertions, but that's not a great solution. In a sense, these are already
factories, but could use some fixing up to generate realistic fake data without needing overrides every time you want
non-hard-coded data.

As for the frontend, we tend to create fake objects that kind of look like the data within each test file, which again
goes back to the issue of needing to understand what the data _should_ look like. We end up with many tests that have
fake objects that don't really look like the data will, e.g. many tests have IDs for objects that are just a few digits
rather than a UUID. There are plenty of times that it doesn't matter too much, but it can make it that much harder for
someone that is looking at existing tests as examples in order to write new tests, but the new tests need more realistic
data.

## Measuring success

Initial success would probably look like us having a decent number of core factories set up and team members feeling
comfortable using them in their tests.

Long-term, we would ideally have replaced our standalone fake objects in FE tests as well as improved our `testadagen`
functions that we use for BE and e2e tests.

### Observability

_How will this change be observed by other team members?_

I'm not 100% sure what this section means...like does this mean how would others keep in line with this ADR? If so, I
think that would look like people creating/updating factories when models are created/updated. Also updating factories
for different states of things, e.g. what a service member looks like a different stages of on-boarding, or what a
shipment looks like at different stages of the move process.

### Ownership

AppEng would own this since it's related to the way we test our application code.

## Considered Alternatives

1. Create factories on the FE (using a package), replace the factories on the BE with a package, and use fake data generators for both the BE and FE.
1. Create factories on the FE (using a package) and use fake data generators for both the BE and FE.
1. Start using fake data generators for the BE/FE, but don't implement factories for FE.
1. Create factorise on the FE, but don't use fake data generators for FE or BE.
1. Leave things as they are.

Each of the options that has us creating factories for the FE could also have an alternative of us creating the
factories from scratch rather than using an existing package, but if we're starting from a blank slate, we might as well
use a supported package that comes with features we'll want out of the box, rather than coding it ourselves and possibly
making mistakes or having feature gaps. I.e. the argument for using open source code in general.

## Decision Outcome

Chosen Alternative -> Option 2: Create factories on the FE and use fake data generators for both the BE and the FE.

This is one of the options that requires the most work up front, but it enables easier testing. With the factories that
use fake data generators, people will be able to more easily create fake data for their tests without having to worry
about all the pre-existing complex interactions of other fields/models. Factory packages also usually have the concept
of traits (term may vary across languages/packages), which are a way of toggling several attributes on a model, e.g. we
could have a trait that sets how far along on-boarding a service member is and set the fields accordingly. This would
also make it easier to test different states for the data.

## Pros and Cons of the Alternatives

### Option 1: Create factories on the FE (using a package), replace the factories on the BE with a package, and use fake data generators for both the BE and FE

* `+` If there was a good library out there for `go`, it would mean we wouldn't have to worry as much about having
    implemented the factories correctly and consistently.
  * See [problem statement](#problem-statement) for example with hard-coded data.
  * Another example of a problem a package could hopefully solve is the consistency of something like our `Stub`
      attribute which isn't properly passed to all other factories in our current implementation.
* `-` There don't really seem to be any big `go` packages for creating factories, most revolve around just generating
    fake data for a `struct`.
* `-` Most work to implement since we'd have to start over with the BE factories.

The rest of the pros/cons match the next option so see that one for more info.

### Option 2: Create factories on the FE (using a package) and use fake data generators for both the BE and FE

* `+` Creating factories for the FE would free up devs from having to know exactly what an object should look like and
    be able to focus only on the part of an object that they are testing.
* `+` Using fake data generators for the FE factories would make it easier to generate realistic fake data for fields
    without having hard-coded data or having to pass in overrides every time how we do on the BE.
* `+` Refactoring our existing FE tests to use the new factories would serve as a test of our newly minted factories and
    might even raise issues that we'd missed previously.
* `+` Using fake data generators for the BE `testdatagen` functions would make them easier to re-use as often as needed
    with less need for passing in overrides each time.
* `-` One argument against fake data generators is that if a test fails, you can't re-run the test with the same exact
    data.
  * This can be mitigated in part by ensuring our test results are recorded (which they are) and by setting our own
      seed value for the fake data generator (if the package we use allows that) to make it easier to get the same
      data on subsequent re-runs.
* `-` Would take work to set up the factories for the FE since we don't have any at all.
* `-` If we do want to replace the existing usages of fake objects in tests, it would take time to do these refactors.
    We might opt to go with the the same update pattern we're using for `react-testing-library` of only updating the
    test if you are editing it for something else. The biggest downside there is that there's a possibility we won't
    come back to it.

### Option 3: Start using fake data generators for the BE/FE, but don't implement factories for FE

* `+` Using fake data generators for the FE objects would make it easier to generate realistic fake data for fields
    without having hard-coded data.
* `+` Using fake data generators for the BE `testdatagen` functions would make them easier to re-use as often as needed
    with less need for passing in overrides each time.
* `-` Not using factories for the FE would still leave us with devs needing to have a deeper knowledge of how our
    objects relate to each other and what the minimum data needed is in order to test what they're trying to focus on.
* `-` Measuring success for the FE is more ambiguous because we won't be replacing all the hard-coded data that exists
    with the generated data since some of that hard-coded data is being used in the tests as is on purpose, so seeing
    what needs to be updated vs left as is would be harder to do at a glance.
* `-` One argument against fake data generators is that if a test fails, you can't re-run the test with the same exact
    data.
  * This can be mitigated in part by ensuring our test results are recorded (which they are) and by setting our own
      seed value for the fake data generator (if the package we use allows that) to make it easier to get the same
      data on subsequent re-runs.

### Option 4: Create factories on the FE, but don't use fake data generators for FE or BE

* `+` Creating factories for the FE would free up devs from having to know exactly what an object should look like and
    be able to focus only on the part of an object that they are testing.
* `-` We are left with factories on the BE that need overrides (assertions) to be passed in every time you want
    different data.
* `-` Having factories on the FE without a fake data generator leaves us with some of the same problems we have with the
    `testdatagen` functions on the backend right now. Namely that we would need to pass overrides in every time we want
    new data instead of the hard-coded defaults.

### Leave things as they are

* `+` No extra work is needed right now.
* `-` We are left with factories on the BE that need overrides (assertions) to be passed in every time you want different data.
* `-` FE tests will continue using fake data that isn't realistic.
* `-` Not using factories for the FE would still leave us with devs needing to have a deeper knowledge of how our
    objects relate to each other and what the minimum data needed is in order to test what they're trying to focus on.

## Resources

* Some possibilities for `go` fake data generators:
  * [jaswdr/faker](https://github.com/jaswdr/faker)
  * [bxcodec/faker](https://github.com/bxcodec/faker)
* Article talking about [why using a factory bot can be good](https://www.codewithjason.com/why-use-factory-bot/).
  * It's for `ruby`, but the idea is applicable in other languages.
* [Slack thread where we discussed faker and factories](https://ustcdp3.slack.com/archives/CTQQJD3G8/p1646079626405239)
* [Front-end check-in notes where we discussed faker and factories](https://dp3.atlassian.net/wiki/spaces/MT/pages/1663500318/2022-03-03+Front+End+Check-In)
* [Back-end check-in notes where we discussed this ADR a bit](https://dp3.atlassian.net/wiki/spaces/MT/pages/1697611790/2022-03-24+Meeting+notes)
