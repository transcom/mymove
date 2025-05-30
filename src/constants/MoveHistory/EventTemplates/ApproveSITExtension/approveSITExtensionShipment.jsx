import React from 'react';

import o from 'constants/MoveHistory/UIDisplay/Operations';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';
import LabeledDetails from 'pages/Office/MoveHistory/LabeledDetails';
import { getMtoShipmentLabel } from 'utils/formatMtoShipment';

const formatChangedValues = (historyRecord) => {
  const newChangedValues = {
    ...getMtoShipmentLabel(historyRecord),
    ...historyRecord.changedValues,
  };
  return { ...historyRecord, changedValues: newChangedValues };
};

export default {
  action: a.UPDATE,
  eventName: o.approveSITExtension,
  tableName: t.mto_shipments,
  getEventNameDisplay: () => 'SIT extension approved',
  getDetails: (historyRecord) => {
    return <LabeledDetails historyRecord={formatChangedValues(historyRecord)} />;
  },
};
