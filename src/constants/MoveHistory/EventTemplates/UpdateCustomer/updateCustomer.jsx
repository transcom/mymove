import React from 'react';

import a from 'constants/MoveHistory/Database/Actions';
import o from 'constants/MoveHistory/UIDisplay/Operations';
import t from 'constants/MoveHistory/Database/Tables';
import LabeledDetails from 'pages/Office/MoveHistory/LabeledDetails';

export default {
  action: a.UPDATE,
  eventName: o.updateCustomer,
  tableName: t.service_members,
  getEventNameDisplay: () => 'Updated profile',
  getDetails: (historyRecord) => <LabeledDetails historyRecord={historyRecord} />,
};
