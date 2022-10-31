import React from 'react';

import o from 'constants/MoveHistory/UIDisplay/Operations';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';
import LabeledDetails from 'pages/Office/MoveHistory/LabeledDetails';

export default {
  action: a.UPDATE,
  eventName: o.updateBillableWeight,
  tableName: t.entitlements,
  getEventNameDisplay: () => 'Updated move',
  getDetails: (historyRecord) => <LabeledDetails historyRecord={historyRecord} />,
};
