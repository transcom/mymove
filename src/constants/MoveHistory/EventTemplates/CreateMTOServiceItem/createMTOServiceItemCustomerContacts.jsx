import React from 'react';

import o from 'constants/MoveHistory/UIDisplay/Operations';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';
import LabeledDetails from 'pages/Office/MoveHistory/LabeledDetails';
import { getMtoShipmentLabel } from 'utils/formatMtoShipment';

const formatChangedValues = (historyRecord) => {
  const { changedValues } = historyRecord;
  const { type, time_military: timeMilitary } = historyRecord.changedValues;

  const newChangedValues = {
    ...getMtoShipmentLabel(historyRecord),
    ...changedValues,
  };

  if (type === 'FIRST') {
    newChangedValues.first_available_delivery_time = timeMilitary;
  } else {
    newChangedValues.second_available_delivery_time = timeMilitary;
    newChangedValues.second_available_delivery_date = changedValues.first_available_delivery_date;
    delete newChangedValues.first_available_delivery_date;
  }

  return { ...historyRecord, changedValues: newChangedValues };
};

export default {
  action: a.INSERT,
  eventName: o.createMTOServiceItem,
  tableName: t.mto_service_item_customer_contacts,
  getEventNameDisplay: () => 'Requested service item',
  getDetails: (historyRecord) => {
    return <LabeledDetails historyRecord={formatChangedValues(historyRecord)} />;
  },
};
