import React from 'react';

import o from 'constants/MoveHistory/UIDisplay/Operations';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';
import { getMtoShipmentLabel } from 'utils/formatMtoShipment';
import LabeledDetailsWithToolTip from 'pages/Office/MoveHistory/LabeledDetailsWithToolTip';

const formatChangedValues = (historyRecord) => {
  const { changedValues } = historyRecord;
  const newChangedValues = {
    ...changedValues,
    ...getMtoShipmentLabel(historyRecord),
  };

  return { ...historyRecord, changedValues: newChangedValues };
};

const generateToolTipText = (historyRecord) => {
  return `Prev SIT entry date: ${historyRecord.oldValues.sit_entry_date}`;
};

export default {
  action: a.UPDATE,
  eventName: o.updateServiceItemSitEntryDate,
  tableName: t.mto_service_items,
  getEventNameDisplay: () => 'Updated SIT entry date',
  getDetails: (historyRecord) => {
    return (
      <div style={{ position: 'relative', display: 'inline-block' }}>
        <LabeledDetailsWithToolTip
          historyRecord={formatChangedValues(historyRecord)}
          toolTipText={generateToolTipText(historyRecord)}
        />
      </div>
    );
  },
};
