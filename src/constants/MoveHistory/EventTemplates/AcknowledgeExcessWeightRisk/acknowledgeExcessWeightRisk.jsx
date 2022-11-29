import React from 'react';

import a from 'constants/MoveHistory/Database/Actions';
import o from 'constants/MoveHistory/UIDisplay/Operations';
import t from 'constants/MoveHistory/Database/Tables';

export default {
  action: a.UPDATE,
  eventName: o.acknowledgeExcessWeightRisk,
  tableName: t.moves,
  getEventNameDisplay: () => 'Updated move',
  getDetails: (historyRecord) => (
    <>{historyRecord?.changedValues?.excess_weight_acknowledged_at ? 'Dismissed excess weight alert' : '-'} </>
  ),
};
