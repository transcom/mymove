import React from 'react';

import d from 'constants/MoveHistory/UIDisplay/DetailsTypes';
import o from 'constants/MoveHistory/UIDisplay/Operations';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';
import LabeledDetails from 'pages/Office/MoveHistory/LabeledDetails';

const formatChangedValues = (historyRecord) => {
  const { context, changedValues } = historyRecord;
  const newChangedValues = {
    ...changedValues,
  };

  if (context[0]?.shipment_type) {
    newChangedValues.shipment_type = context[0].shipment_type;
    newChangedValues.shipment_id_display = context[0].shipment_id_abbr.toUpperCase();
  }

  return { ...historyRecord, changedValues: newChangedValues };
};

export default {
  action: a.UPDATE,
  eventName: o.updateMTOShipment,
  tableName: t.mto_shipments,
  detailsType: d.LABELED,
  getEventNameDisplay: () => 'Updated shipment',
  getDetails: (historyRecord) => <LabeledDetails historyRecord={formatChangedValues(historyRecord)} />,
};
