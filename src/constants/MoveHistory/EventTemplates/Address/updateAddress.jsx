import React from 'react';

import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';
import LabeledDetails from 'pages/Office/MoveHistory/LabeledDetails';
import { formatMoveHistoryFullAddress } from 'utils/formatters';

const ADDRESS_LABEL = {
  backupMailingAddress: 'backup_address',
  destinationAddress: 'destination_address',
  residentialAddress: 'residential_address',
  pickupAddress: 'pickup_address',
};

const formatChangedValues = (historyRecord) => {
  const { context, changedValues, oldValues } = historyRecord;
  const address = formatMoveHistoryFullAddress(changedValues);

  const addressType = context.filter((contextObject) => contextObject.address_type)[0].address_type;
  const addressLabel = ADDRESS_LABEL[addressType];

  const newChangedValues = {
    street_address_1: oldValues.street_address_1,
    street_address_2: oldValues.street_address_2,
    city: oldValues.city,
    state: oldValues.state,
    postal_code: oldValues.postal_code,
    [addressLabel]: address,
    ...changedValues,
  };

  if (context[0]?.shipment_type) {
    newChangedValues.shipment_type = context[0].shipment_type;
  }

  return { ...historyRecord, changedValues: newChangedValues };
};

export default {
  action: a.UPDATE,
  eventName: '*',
  tableName: t.addresses,
  getEventNameDisplay: () => 'Updated address',
  getDetails: (historyRecord) => <LabeledDetails historyRecord={formatChangedValues(historyRecord)} />,
};
