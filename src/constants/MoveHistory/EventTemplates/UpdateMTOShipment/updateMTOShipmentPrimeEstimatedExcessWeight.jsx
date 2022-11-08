import React from 'react';

import o from 'constants/MoveHistory/UIDisplay/Operations';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';

export default {
  action: a.UPDATE,
  eventName: o.updateMTOShipment,
  tableName: t.moves,
  getEventNameDisplay: () => 'Updated move',
  getDetails: (historyRecord) => (
    <div>
      {historyRecord.changedValues.excess_weight_qualified_at
        ? 'Flagged for excess weight, total estimated weight > 90% weight allowance'
        : '-'}
    </div>
  ),
};
