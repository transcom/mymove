import React from 'react';

import o from 'constants/MoveHistory/UIDisplay/Operations';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';
import LabeledDetails from 'pages/Office/MoveHistory/LabeledDetails';

const formatChangedValues = (historyRecord) => {
  const newChangedValues = {
    shipment_type: historyRecord.context[0]?.shipment_type,
    service_item_name: historyRecord.context[0]?.name,
    ...historyRecord.changedValues,
  };
  return { ...historyRecord, changedValues: newChangedValues };
};

export default {
  action: a.INSERT,
  eventName: o.createMTOServiceItem,
  tableName: t.mto_service_items,
  getEventNameDisplay: () => 'Requested service item',
  getDetails: (historyRecord) => {
    return <LabeledDetails historyRecord={formatChangedValues(historyRecord)} />;
  },
};
