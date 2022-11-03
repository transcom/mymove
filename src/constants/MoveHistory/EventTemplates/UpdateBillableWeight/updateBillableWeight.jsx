import React from 'react';

import o from 'constants/MoveHistory/UIDisplay/Operations';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';
import LabeledDetails from 'pages/Office/MoveHistory/LabeledDetails';

// To-do: Remove one max_billable_weight is its own value
const formatChangedValues = (historyRecord) => {
  const { changedValues } = historyRecord;
  const newChangedValues = { ...changedValues };
  if (changedValues.authorized_weight) {
    newChangedValues.max_billable_weight = changedValues.authorized_weight;
    delete newChangedValues.authorized_weight;
  }
  return { ...historyRecord, changedValues: newChangedValues };
};

export default {
  action: a.UPDATE,
  eventName: o.updateBillableWeight,
  tableName: t.entitlements,
  getEventNameDisplay: () => 'Updated move',
  getDetails: (historyRecord) => <LabeledDetails historyRecord={formatChangedValues(historyRecord)} />,
};
