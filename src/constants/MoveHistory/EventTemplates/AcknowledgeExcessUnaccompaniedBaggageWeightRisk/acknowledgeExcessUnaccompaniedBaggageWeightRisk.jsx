import React from 'react';

import a from 'constants/MoveHistory/Database/Actions';
import o from 'constants/MoveHistory/UIDisplay/Operations';
import t from 'constants/MoveHistory/Database/Tables';
import LabeledDetails from 'pages/Office/MoveHistory/LabeledDetails';

export default {
  action: a.UPDATE,
  eventName: o.acknowledgeExcessUnaccompaniedBaggageWeightRisk,
  tableName: t.moves,
  getEventNameDisplay: () => 'Updated move',
  getDetails: (historyRecord) => {
    if (historyRecord?.changedValues?.excess_unaccompanied_baggage_weight_acknowledged_at) {
      return 'Dismissed excess unaccompanied baggage weight alert';
    }
    return <LabeledDetails historyRecord={historyRecord} />;
  },
};
