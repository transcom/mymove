import React from 'react';

import o from 'constants/MoveHistory/UIDisplay/Operations';
import t from 'constants/MoveHistory/Database/Tables';
import LabeledDetails from 'pages/Office/MoveHistory/LabeledDetails';
import { getMtoShipmentLabel } from 'utils/formatMtoShipment';

const formatChangedValues = (historyRecord) => {
  const { changedValues } = historyRecord;
  const newChangedValues = {
    ...changedValues,
    ...getMtoShipmentLabel(historyRecord),
  };

  if (historyRecord.context[0].shipment_type === 'PPM') {
    newChangedValues.status = null; // Status will always be draft when shipment is created, not useful information.
  }
  return { ...historyRecord, changedValues: newChangedValues };
};

export default {
  action: '*',
  eventName: o.createMTOShipment,
  tableName: t.mto_shipments,
  getEventNameDisplay: (historyRecord) => {
    if (historyRecord.context[0].shipment_type === 'PPM') {
      return 'PPM Shipment Created';
    }
    return 'Created shipment';
  },
  getDetails: (historyRecord) => <LabeledDetails historyRecord={formatChangedValues(historyRecord)} />,
};
