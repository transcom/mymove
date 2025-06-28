import React from 'react';

import o from 'constants/MoveHistory/UIDisplay/Operations';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';
import LabeledDetails from 'pages/Office/MoveHistory/LabeledDetails';

const formatChangedValues = (historyRecord) => {
  const { changedValues } = historyRecord;
  const newChangedValues = {
    status: changedValues.status,
    office_remarks: changedValues.office_remarks,
  };
  return { changedValues: newChangedValues };
};

export default {
  action: a.UPDATE,
  eventName: o.reviewShipmentAddressUpdate,
  tableName: t.shipment_address_updates,
  getEventNameDisplay: () => 'Shipment destination address update',
  getDetails: (historyRecord) => {
    return <LabeledDetails historyRecord={formatChangedValues(historyRecord)} />;
  },
};
