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

  return { ...historyRecord, changedValues: newChangedValues };
};

export default {
  action: a.UPDATE,
  eventName: o.updatePaymentServiceItemStatus,
  tableName: t.payment_service_items,
  getEventNameDisplay: (historyRecord) => {
    let actionPrefix = '';
    if (historyRecord.changedValues.status === 'DENIED') {
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
