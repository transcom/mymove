import React from 'react';

import o from 'constants/MoveHistory/UIDisplay/Operations';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';
import LabeledDetails from 'pages/Office/MoveHistory/LabeledDetails';
import { getMtoShipmentLabel } from 'utils/formatMtoShipment';

const formatChangedValues = (historyRecord) => {
  const { changedValues } = historyRecord;

  const newChangedValues = {
    ...changedValues,
    ...getMtoShipmentLabel(historyRecord),
    reweigh_weight: changedValues.weight,
  };

  return { ...historyRecord, changedValues: newChangedValues };
};

export default {
  action: a.UPDATE,
  eventName: o.updateReweigh,
  tableName: t.reweighs,
  getEventNameDisplay: () => `Updated shipment`,
  getDetails: (historyRecord) => <LabeledDetails historyRecord={formatChangedValues(historyRecord)} />,
};
