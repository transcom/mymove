import React from 'react';

import t from 'constants/MoveHistory/Database/Tables';
import a from 'constants/MoveHistory/Database/Actions';
import LabeledDetails from 'pages/Office/MoveHistory/LabeledDetails';

export default {
  action: a.INSERT,
  eventName: '*', // Needs wild card to handle both createOrders and createOrder
  tableName: t.moves,
  getEventNameDisplay: () => 'Created move',
  getDetails: (historyRecord) => {
    return <LabeledDetails historyRecord={historyRecord} />;
  },
};
