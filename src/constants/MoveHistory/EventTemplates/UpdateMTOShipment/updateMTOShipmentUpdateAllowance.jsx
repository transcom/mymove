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

export const isUpdatedAllowances = (changedValues) => {
  return changedValues.gun_safe !== undefined || changedValues.gun_safe_weight !== undefined;
};

export default {
  action: a.UPDATE,
  eventName: o.updateMTOShipment,
  tableName: t.entitlements,
  getEventNameDisplay: ({ changedValues }) =>
    isUpdatedAllowances(changedValues) ? 'Updated allowances' : 'Updated shipment',
  getDetails: (historyRecord) => <LabeledDetails historyRecord={formatChangedValues(historyRecord)} />,
};
