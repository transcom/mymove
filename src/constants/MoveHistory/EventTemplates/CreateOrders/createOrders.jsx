import React from 'react';

import t from 'constants/MoveHistory/Database/Tables';
import a from 'constants/MoveHistory/Database/Actions';
import o from 'constants/MoveHistory/UIDisplay/Operations';
import LabeledDetails from 'pages/Office/MoveHistory/LabeledDetails';

const formatChangedValues = (historyRecord) => {
  let newChangedValues = {
    ...historyRecord.changedValues,
  };

  if (historyRecord.context) {
    newChangedValues = {
      ...newChangedValues,
      ...historyRecord.context[0],
    };
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
