import React from 'react';

import o from 'constants/MoveHistory/UIDisplay/Operations';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';
import LabeledDetails from 'pages/Office/MoveHistory/LabeledDetails';
import { formatAssignedOfficeUserFromContext } from 'utils/formatters';

const formatChangedValues = (historyRecord) => {
  const newChangedValues = {
    ...formatAssignedOfficeUserFromContext(historyRecord),
  };

  return { ...historyRecord, changedValues: newChangedValues };
};

export default {
  action: a.UPDATE,
  eventName: o.saveBulkAssignmentData,
  tableName: t.moves,
  getEventNameDisplay: () => 'Move assignment updated',
  getDetails: (historyRecord) => <LabeledDetails historyRecord={formatChangedValues(historyRecord)} />,
};
