import React from 'react';

import o from 'constants/MoveHistory/UIDisplay/Operations';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';
import LabeledDetails from 'pages/Office/MoveHistory/LabeledDetails';

const formatChangedValues = (historyRecord) => {
  let newChangedValues = {
    ...historyRecord.changedValues,
  };

  if (historyRecord.context) {
    newChangedValues = {
      ...newChangedValues,
      ...historyRecord.context[0],
    };
  }

  return { ...historyRecord, changedValues: newChangedValues };
};

export default {
  action: a.UPDATE,
  eventName: o.createSITExtension,
  tableName: t.mto_shipments,
  getEventNameDisplay: () => 'Updated shipment',
  getDetails: (historyRecord) => {
    return <LabeledDetails historyRecord={formatChangedValues(historyRecord)} />;
  },
};
