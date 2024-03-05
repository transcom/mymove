import React from 'react';

import o from 'constants/MoveHistory/UIDisplay/Operations';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';
import LabeledDetails from 'pages/Office/MoveHistory/LabeledDetails';
import { getMtoShipmentLabel } from 'utils/formatMtoShipment';
import { formatDataForPPM } from 'utils/formatPPMData';

const formatChangedValues = (historyRecord) => {
  const newChangedValues = {
    ...historyRecord.changedValues,
    ...getMtoShipmentLabel(historyRecord),
    ...formatDataForPPM(historyRecord),
  };

  return { ...historyRecord, changedValues: newChangedValues };
};

export default {
  action: a.UPDATE,
  eventName: o.updateMTOShipment,
  tableName: t.ppm_shipments,
  getEventNameDisplay: () => 'Updated shipment',
  getDetails: (historyRecord) => <LabeledDetails historyRecord={formatChangedValues(historyRecord)} />,
};
