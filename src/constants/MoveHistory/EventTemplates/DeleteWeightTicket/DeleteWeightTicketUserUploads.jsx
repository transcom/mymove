import React from 'react';

import o from 'constants/MoveHistory/UIDisplay/Operations';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';
import { getMtoShipmentLabel } from 'utils/formatMtoShipment';
import LabeledDetails from 'pages/Office/MoveHistory/LabeledDetails';
import { formatDataForPPM } from 'utils/formatPPMData';

const formatChangedValues = (historyRecord) => {
  const newChangedValues = {
    ...historyRecord.changedValues,
    ...getMtoShipmentLabel(historyRecord),
    ...formatDataForPPM(historyRecord),
  };

  return { ...historyRecord, changedValues: newChangedValues };
};

export default {
  action: a.UPDATE,
  eventName: o.deleteWeightTicket,
  tableName: t.user_uploads,
  getEventNameDisplay: ({ context }) => {
    const eventLabel =
      context[0]?.upload_type === 'fullWeightTicket' || context[0]?.upload_type === 'emptyWeightTicket'
        ? 'Deleted trip document'
        : 'Deleted document';

    return <div>{eventLabel}</div>;
  },
  getDetails: (historyRecord) => {
    return <LabeledDetails historyRecord={formatChangedValues(historyRecord)} />;
  },
};
