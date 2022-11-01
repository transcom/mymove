import React from 'react';

import o from 'constants/MoveHistory/UIDisplay/Operations';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';
import LabeledDetails from 'pages/Office/MoveHistory/LabeledDetails';

const formatChangedValues = (historyRecord) => {
  const newChangedValues = { ...historyRecord.changedValues, status: historyRecord.oldValues.status };
  return {
    ...historyRecord,
    changedValues: newChangedValues,
  };
};

export default {
  action: a.UPDATE,
  eventName: o.createMTOServiceItem,
  tableName: t.moves,
  getEventNameDisplay: () => 'Updated move',
  getDetails: (historyRecord) => <LabeledDetails historyRecord={formatChangedValues(historyRecord)} />,
};
