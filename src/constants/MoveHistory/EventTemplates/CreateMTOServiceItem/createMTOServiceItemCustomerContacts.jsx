import React from 'react';

import o from 'constants/MoveHistory/UIDisplay/Operations';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';
import LabeledDetails from 'pages/Office/MoveHistory/LabeledDetails';
import { getMtoShipmentLabel } from 'utils/formatMtoShipment';

const formatChangedValues = (historyRecord) => {
  const { changedValues } = historyRecord;
  const { type, time_military: timeMilitary } = historyRecord.changedValues;

  const deliveryTimeOrder = type === 'FIRST' ? 'first_available_delivery_time' : 'second_available_delivery_time';
  const deliveryDateOrder = type === 'FIRST' ? 'first_available_delivery_date' : 'second_available_delivery_date';
  const newChangedValues = {
    ...getMtoShipmentLabel(historyRecord),
    ...changedValues,
  };

  if (type === 'FIRST') {
    newChangedValues[deliveryTimeOrder] = timeMilitary;
  } else {
    newChangedValues[deliveryTimeOrder] = timeMilitary;
    newChangedValues[deliveryDateOrder] = changedValues.first_available_delivery_date;
    delete newChangedValues.first_available_delivery_date;
  }
  // eslint-disable-next-line no-console
  console.log({ ...historyRecord, changedValues: newChangedValues });
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
