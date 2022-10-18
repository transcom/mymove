import React from 'react';

import o from 'constants/MoveHistory/UIDisplay/Operations';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';
import LabeledDetails from 'pages/Office/MoveHistory/LabeledDetails';

const formatChangedValues = (historyRecord) => {
  const { type, time_military: timeMilitary } = historyRecord.changedValues;
  const deliveryTimeOrder = type === 'FIRST' ? 'first_available_delivery_time' : 'second_available_delivery_time';
  const newChangedValues = {
    ...historyRecord.changedValues,
  };

  newChangedValues[deliveryTimeOrder] = timeMilitary;
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
