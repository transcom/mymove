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
  const newOldValues = { ...historyRecord.oldValues };
  if (historyRecord.context[0]?.moving_expense_type)
    newOldValues.moving_expense_type = historyRecord.context[0].moving_expense_type;
  if (historyRecord.context[0]?.upload_type === 'proGearWeightTicket') newOldValues.belongs_to_self = true;
  else if (historyRecord.context[0]?.upload_type === 'spouseProGearWeightTicket') newOldValues.belongs_to_self = false;

  return { ...historyRecord, changedValues: newChangedValues, oldValues: newOldValues };
};

export default {
  action: a.INSERT,
  eventName: o.createPPMUpload,
  tableName: t.user_uploads,
  getEventNameDisplay: ({ context }) => {
    const eventLabel =
      context[0]?.upload_type === 'fullWeightTicket' || context[0]?.upload_type === 'emptyWeightTicket'
        ? 'Uploaded trip document'
        : 'Uploaded document';

    return <div>{eventLabel}</div>;
  },
  getDetails: (historyRecord) => {
    return <LabeledDetails historyRecord={formatChangedValues(historyRecord)} />;
  },
};
