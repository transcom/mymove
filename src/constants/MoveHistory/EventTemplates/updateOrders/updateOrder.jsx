import React from 'react';

import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';
import LabeledDetails from 'pages/Office/MoveHistory/LabeledDetails';

const formatChangedValues = (historyRecord) => {
  let newChangedValues = historyRecord.changedValues;

  if (historyRecord.context) {
    newChangedValues = {
      ...newChangedValues,
      ...historyRecord.context[0],
    };
  }

  if ('has_dependents' in newChangedValues) {
    newChangedValues.has_dependents = newChangedValues.has_dependents ? 'Yes' : 'No';
  }

  return { ...historyRecord, changedValues: newChangedValues };
};

export default {
  action: a.UPDATE,
  eventName: '*',
  tableName: t.orders,
  getEventNameDisplay: () => 'Updated orders',
  getDetails: (historyRecord) => {
    return <LabeledDetails historyRecord={formatChangedValues(historyRecord)} />;
  },
};
