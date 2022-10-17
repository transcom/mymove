import React from 'react';

import t from 'constants/MoveHistory/Database/Tables';
import a from 'constants/MoveHistory/Database/Actions';
import o from 'constants/MoveHistory/UIDisplay/Operations';
import LabeledDetails from 'pages/Office/MoveHistory/LabeledDetails';

const formatChangedValues = (historyRecord) => {
  let newChangedValues;

  if (historyRecord.context) {
    newChangedValues = {
      ...historyRecord.changedValues,
      ...historyRecord.context[0],
    };
  }

  if ('has_dependents' in newChangedValues) {
    newChangedValues.has_dependents = newChangedValues.has_dependents ? 'Yes' : 'No';
  }

  return { ...historyRecord, changedValues: newChangedValues };
};

export default {
  action: a.INSERT,
  eventName: o.createOrders,
  tableName: t.orders,
  getEventNameDisplay: () => 'Created orders',
  getDetails: (historyRecord) => {
    return <LabeledDetails historyRecord={formatChangedValues(historyRecord)} />;
  },
};
