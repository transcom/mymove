import React from 'react';

import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';
import o from 'constants/MoveHistory/UIDisplay/Operations';
import LabeledDetails from 'pages/Office/MoveHistory/LabeledDetails';

export default {
  action: a.UPDATE,
  eventName: o.approveShipment,
  tableName: t.entitlements,
  getEventNameDisplay: () => 'Updated allowances',
  getDetails: (historyRecord) => <LabeledDetails historyRecord={historyRecord} />,
};
