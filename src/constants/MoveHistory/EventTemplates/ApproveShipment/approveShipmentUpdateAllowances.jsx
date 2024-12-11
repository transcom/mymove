import React from 'react';

import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';
import o from 'constants/MoveHistory/UIDisplay/Operations';
import LabeledDetails from 'pages/Office/MoveHistory/LabeledDetails';

// Use formatter from updateBillableWeight.jsx to keep consistency
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
  eventName: o.approveShipment,
  tableName: t.entitlements,
  getEventNameDisplay: () => 'Updated shipment',
  getDetails: (historyRecord) => <LabeledDetails historyRecord={formatChangedValues(historyRecord)} />,
};
