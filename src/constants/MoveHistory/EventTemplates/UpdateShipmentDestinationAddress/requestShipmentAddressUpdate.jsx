import React from 'react';

import o from 'constants/MoveHistory/UIDisplay/Operations';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';
import LabeledDetails from 'pages/Office/MoveHistory/LabeledDetails';

const formatChangedValues = (historyRecord) => {
  const { changedValues } = historyRecord;
  const newChangedValues = {
    contractor_remarks: changedValues.contractor_remarks,
  };
  return { changedValues: newChangedValues };
};

export default {
  action: a.UPDATE,
  eventName: o.updateShipmentDestinationAddress,
  tableName: t.shipment_address_updates,
  getEventNameDisplay: () => 'Shipment destination address request',
  getDetails: (historyRecord) => {
    return <LabeledDetails historyRecord={formatChangedValues(historyRecord)} />;
  },
};
