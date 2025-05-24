import React from 'react';

import o from 'constants/MoveHistory/UIDisplay/Operations';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';
import LabeledDetails from 'pages/Office/MoveHistory/LabeledDetails';
import { formatMoveHistoryMaxBillableWeight, formatMoveHistoryGunSafe } from 'utils/formatters';

export const formatChangedValues = (historyRecord) => {
  const formattedRecord = formatMoveHistoryGunSafe(historyRecord);
  return formatMoveHistoryMaxBillableWeight(formattedRecord);
};

export default {
  action: a.UPDATE,
  eventName: o.updateMTOShipment,
  tableName: t.entitlements,
  getEventNameDisplay: () => 'Updated shipment',
  getDetails: (historyRecord) => <LabeledDetails historyRecord={formatChangedValues(historyRecord)} />,
};
