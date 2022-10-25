import React from 'react';

import o from 'constants/MoveHistory/UIDisplay/Operations';
import d from 'constants/MoveHistory/UIDisplay/DetailsTypes';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';
import LabeledDetails from 'pages/Office/MoveHistory/LabeledDetails';

const formatChangedValues = (historyRecord) => {
  const { changedValues, context } = historyRecord;
  let newChangedValues = changedValues;

  if (historyRecord.context) {
    newChangedValues = {
      ...newChangedValues,
      shipment_type: context[0]?.shipment_type,
      shipment_id_display: context[0]?.shipment_id_abbr.toUpperCase(),
      service_item_name: context[0]?.name,
    };
  }

  return { ...historyRecord, changedValues: newChangedValues };
};

export default {
  action: a.UPDATE,
  eventName: o.updateMTOServiceItem,
  tableName: t.mto_service_items,
  detailsType: d.LABELED,
  getEventNameDisplay: () => 'Updated service item',
  getDetails: (historyRecord) => {
    return <LabeledDetails historyRecord={formatChangedValues(historyRecord)} />;
  },
};
