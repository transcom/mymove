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
  eventName: o.updateWeightTicket,
  tableName: t.mto_shipments,
  getEventNameDisplay: () => {
    return <div>Updated shipment</div>;
  },
  getDetails: (historyRecord) => {
    return <LabeledDetails historyRecord={formatChangedValues(historyRecord)} />;
  },
};
