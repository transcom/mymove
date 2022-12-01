import React from 'react';

import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';
import LabeledDetails from 'pages/Office/MoveHistory/LabeledDetails';

const formatChangedValues = (historyRecord) => {
  const { changedValues, oldValues } = historyRecord;
  let newChangedValues = { ...changedValues };

  const ordersType = newChangedValues.orders_type ?? oldValues.orders_type;
  if (changedValues.report_by_date) {
    if (ordersType === 'RETIREMENT') {
      newChangedValues.retirement_date = changedValues.report_by_date;
      delete newChangedValues.report_by_date;
    }
    if (ordersType === 'SEPARATION') {
      newChangedValues.separation_date = changedValues.report_by_date;
      delete newChangedValues.report_by_date;
    }
  }

  if (historyRecord.context) {
    newChangedValues = {
      ...newChangedValues,
      ...historyRecord.context[0],
    };
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
