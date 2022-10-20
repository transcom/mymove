import React from 'react';

import o from 'constants/MoveHistory/UIDisplay/Operations';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';
import LabeledDetails from 'pages/Office/MoveHistory/LabeledDetails';

export default {
  action: a.UPDATE,
  eventName: o.updateAllowance,
  tableName: t.service_members,
  getEventNameDisplay: () => 'Updated service member',
  getDetails: (historyRecord) => {
    return <LabeledDetails historyRecord={historyRecord} />;
  },
};
