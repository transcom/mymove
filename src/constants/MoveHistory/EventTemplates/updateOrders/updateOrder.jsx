import React from 'react';

import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';
import LabeledDetails from 'pages/Office/MoveHistory/LabeledDetails';

const formatChangedValues = (historyRecord) => {
  let newChangedValues;

  if (historyRecord.context) {
    newChangedValues = {
      ...historyRecord.changedValues,
      ...historyRecord.context[0],
    };
  } else {
    newChangedValues = historyRecord.changedValues;
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
