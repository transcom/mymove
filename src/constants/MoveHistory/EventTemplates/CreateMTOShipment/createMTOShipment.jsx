import React from 'react';

import a from 'constants/MoveHistory/Database/Actions';
import o from 'constants/MoveHistory/UIDisplay/Operations';
import d from 'constants/MoveHistory/UIDisplay/DetailsTypes';
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
  action: a.INSERT,
  eventName: o.createMTOShipment,
  tableName: t.mto_shipments,
  detailsType: d.LABELED,
  getEventNameDisplay: () => 'Created shipment',
  getDetails: (historyRecord) => <LabeledDetails historyRecord={formatChangedValues(historyRecord)} />,
};
