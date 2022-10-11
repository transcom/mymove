import React from 'react';

import d from 'constants/MoveHistory/UIDisplay/DetailsTypes';
import t from 'constants/MoveHistory/Database/Tables';
import a from 'constants/MoveHistory/Database/Actions';
import o from 'constants/MoveHistory/UIDisplay/Operations';
import LabeledDetails from 'pages/Office/MoveHistory/LabeledDetails';

export default {
  action: a.INSERT,
  eventName: o.createOrders,
  tableName: t.moves,
  detailsType: d.PLAIN_TEXT,
  getEventNameDisplay: () => 'Created move',
  getDetails: (historyRecord) => {
    return <LabeledDetails historyRecord={historyRecord} />;
  },
};
