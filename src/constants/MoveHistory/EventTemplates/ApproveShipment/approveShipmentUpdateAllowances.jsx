import React from 'react';

import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';
import LabeledDetails from 'pages/Office/MoveHistory/LabeledDetails';
import { formatMoveHistoryMaxBillableWeight } from 'utils/formatters';

export default {
  action: a.UPDATE,
  eventName: '*', // both approveShipment and approveShipments events can render this template
  tableName: t.entitlements,
  getEventNameDisplay: () => 'Updated shipment',
  getDetails: (historyRecord) => <LabeledDetails historyRecord={formatMoveHistoryMaxBillableWeight(historyRecord)} />,
};
