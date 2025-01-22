import React from 'react';

import o from 'constants/MoveHistory/UIDisplay/Operations';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';
import { getMtoShipmentLabel } from 'utils/formatMtoShipment';
import LabeledDetails from 'pages/Office/MoveHistory/LabeledDetails';

const formatChangedValues = (historyRecord) => {
  const newChangedValues = {
    ...historyRecord.changedValues,
    ...getMtoShipmentLabel(historyRecord),
  };

  // Removed unneeded values to avoid clutter in audit log
  if (newChangedValues.status === 'APPROVED') {
    delete newChangedValues.rejection_reason;
  }

  delete newChangedValues.status;
  const newHistoryRecord = { ...historyRecord };
  delete newHistoryRecord.changedValues.status;
  return { ...newHistoryRecord, changedValues: newChangedValues };
};

export default {
  action: a.UPDATE,
  eventName: o.updatePaymentServiceItemStatus,
  tableName: t.payment_service_items,
  getEventNameDisplay: (historyRecord) => {
    let actionPrefix = '';
    if (historyRecord.changedValues.rejection_reason !== null || historyRecord.changedValues.status === 'REJECTED') {
      actionPrefix = 'Rejected';
    } else if (historyRecord.changedValues.status === 'APPROVED') {
      actionPrefix = 'Approved';
    } else {
      actionPrefix = 'Updated';
    }
    return <div>{actionPrefix} Payment Service Item</div>;
  },
  getDetails: (historyRecord) => {
    return <LabeledDetails historyRecord={formatChangedValues(historyRecord)} />;
  },
};
