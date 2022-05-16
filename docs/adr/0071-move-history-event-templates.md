# Introduce Move History Event Templates

🔒 **User Story:** [_MB-8115_](https://dp3.atlassian.net/browse/MB-8115)

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

Some notes to noodle on for more background information

> Investigate if there is a better way to generate Event Templates and possibly
> load them dynamically based on what folder they're in.

