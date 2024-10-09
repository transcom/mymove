import React from 'react';

import o from 'constants/MoveHistory/UIDisplay/Operations';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';
import LabeledDetails from 'pages/Office/MoveHistory/LabeledDetails';

export default {
  action: a.UPDATE,
  eventName: o.updateMTOShipment,
  tableName: t.moves,
  getEventNameDisplay: () => 'Updated move',
  getDetails: (historyRecord) => {
    if (historyRecord?.changedValues?.excess_weight_qualified_at) {
      return 'Flagged for excess weight, total estimated weight > 90% weight allowance';
    }
    return <LabeledDetails historyRecord={historyRecord} />;
  },
};
