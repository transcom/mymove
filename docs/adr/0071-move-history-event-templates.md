# Introduce Move History Events

- 🔒 **User Story:** [_MB-8115_](https://dp3.atlassian.net/browse/MB-8115)
- 🔒 **User Story:** [_MB-12515_](https://dp3.atlassian.net/browse/MB-12515)

There are a number of additions to the front-end codebase for Move History Log.
The changes discussed in this ADR are around the files found in `constants/`.
These files are currently flat within the `src/constants/` directory with names
that start with `moveHistory` or `historyLog`. These files all have exports that
are leveraged in a number of ways within the files themselves and also within
the Move History components located at `src/pages/Office/MoveHistory/`.

## Constant Objects

There are a number of JavaScript objects that are used to map responses from the
API into human-readable strings. These objects are all within a single file
currently. Any updates to these objects has the potential of causing merge
conflicts.

## Events

With the History Log outcome, the MilMove team needs to create event templates
which are JavaScript Objects that facilitate the customization, rendering, and
identifying of different events being used in the History Log. The structure of
these Objects are outlined below.

```javascript
{
  action: '',
  eventName: '',
  tableName: '',
  detailsType: '',
  getEventNameDisplay: (historyRecord) => '',
  // One of the following functions are used to populate the details column of
  // an event. The function that is called is based on the `detailsType`
  // property above.
  getDetailsPlainText: (historyRecord) => '',
  getDetailsLabeledDetails: (historyRecord) => {
    let newChangedValues = {};
    // add to newChangedValues Object.
    return newChangedValues
  },
  getStatusDetails: (historyRecord) => {
    let newChangedValues = {};
    // add to newChangedValues Object.
    return newChangedValues
  },
}
```

### Event Template

This ADR is introducing this Object above for Move History Event Templates. This
ADR is not an exhaustive list of the types of Details column will be displayed
but include three examples based on the `detailsType` property above. Currently
the maintenance and addition of these event types requires engineers to edit a
single file named `src/constants/moveHistory/moveHistoryEventTemplate.js`. This
has proven to lead to many merge conflicts as engineers are adding event
templates to the same file. This file structure makes it non-trivial to verify
which event names have been added to the project because there is only a single
file that ever gets updated. Tests written for testing the event names happen at
the event template file rather than near the actual event templates that are
being tested. This has proven to create scenarios where certain events are not
tested because they are mistakenly forgotten to be added. At the time of this
writing, the file above is over 550 lines of code supporting 32 distinct event
templates. At the time of this writing, it's estimated that we will be adding
another 21 distinct event templates.

### Proposal

Organizing these `src/constants/` files for Move History into individual
JavaScript modules is a great way to encapsulate any changes so multiple
engineers are able to work on the PO9 outcome at the same time without working
within the same file. This will lead to much less merge conflicts around similar
features. The added benefits of individual modules is that anyone with access to
the repository on GitHub can see what Move History Events are supported by
looking for the folders under `src/constants/MoveHistory/EventTemplates/`.

Below is an example of the proposed file structure for Move History constants.

```sh
src/constants/MoveHistory
├── BuildTemplate.js
├── GetTemplate.js
└── EventTemplates
   ├── updateMTOReviewedBillableWeightAt.js
   └── updateMTOReviewedBillableWeightAt.test.js
└── Database
   ├── Tables.js
   └── FieldToDisplayName.js
└── LabeledFields
   ├── OrdersOptions.js
└── UIDisplay
   └── HistoryLogRecordShapes.js
```

This example is a minimal set of files. A refactor for all the current exports
will include many more files under certain directories. The file and directory
names are capitalized and camel-case for readability. The event templates should be named after the user action to which they correspond (usually based on their event name) and their names should not end in "event." Any tests for these files
will be written alongside the files that they are testing. Each template test file would technically test the result of GetTemplate.js and use the event template as the expected result. For `EventTemplates`
tests, due to this refactor there will no longer be a test for the
`moveHistoryEventTemplate` file as this file will only be responsible for
getting templates. The main export of this file will be a new file called
`GetTemplate.js` and won't require tests as it is only a helper function.

Most of the functionality of `GetTemplate.js` will remain the same. This may
lead to merge conflicts, but they will be relatively painless as the main reason
of this file is to keep track and load specific `EventTemplates` files.

## Considered Alternatives

- _Do nothing_
- _Organize constants files into individual modules_ (**chosen**)

## Decision Outcome

### Chosen Alternative: _Organize constants files into individual modules_

#### Justification: This is the only decision which will help prevent merge conflicts

- `+` _The outcome is much easier to follow by looking at the files_
- `+` _Prevents merge conflicts as engineers are able to work separately on
  tests and features in separate files_
- `-` _Implementing this will block work for PO9 outcome while it's getting done_

## Pros and Cons of the Alternatives

### _Do nothing_

- `+` _No effort_
- `-` _Continues to cause merge conflicts as developers work on the same files_
- `-` _Understanding the amount of Events that are complete is opaque_
