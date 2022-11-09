import React from 'react';

import o from 'constants/MoveHistory/UIDisplay/Operations';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';
import LabeledDetails from 'pages/Office/MoveHistory/LabeledDetails';

const formatChangedValues = (historyRecord) => {
  const { changedValues, context } = historyRecord;

  const newChangedValues = {
    ...changedValues,
    shipment_type: context[0]?.shipment_type,
    shipment_id_display: context[0]?.shipment_id_abbr.toUpperCase(),
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
