import React from 'react';

import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';
import o from 'constants/MoveHistory/UIDisplay/Operations';
import LabeledDetails from 'pages/Office/MoveHistory/LabeledDetails';
import { formatMoveHistoryMaxBillableWeight } from 'utils/formatters';

export default {
  action: a.UPDATE,
  eventName: o.approveShipments,
  tableName: t.entitlements,
  getEventNameDisplay: () => 'Updated shipment',
  getDetails: (historyRecord) => <LabeledDetails historyRecord={formatMoveHistoryMaxBillableWeight(historyRecord)} />,
};
