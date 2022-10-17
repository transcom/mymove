import React from 'react';

import { formatMoveHistoryFullAddress } from 'utils/formatters';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';
import LabeledDetails from 'pages/Office/MoveHistory/LabeledDetails';

const ADDRESS_LABEL = {
  backupMailingAddress: 'backup_address',
  destinationAddress: 'destination_address',
  residentialAddress: 'residential_address',
  pickupAddress: 'pickup_address',
};

const formatChangedValues = (historyRecord) => {
  const { context, changedValues } = historyRecord;
  const address = formatMoveHistoryFullAddress(changedValues);

  const addressType = context.filter((contextObject) => contextObject.address_type)[0].address_type;
  const addressLabel = ADDRESS_LABEL[addressType];

  const newChangedValues = {
    [addressLabel]: address,
    ...changedValues,
  };

  if (context[0]?.shipment_type) {
    newChangedValues.shipment_type = context[0].shipment_type;
    newChangedValues.shipment_id_display = context[0].shipment_id_abbr.toUpperCase();
  }

  return { ...historyRecord, changedValues: newChangedValues };
};

export default {
  action: a.INSERT,
  eventName: '*',
  tableName: t.addresses,
  getEventNameDisplay: () => 'Updated address',
  getDetails: (historyRecord) => <LabeledDetails historyRecord={formatChangedValues(historyRecord)} />,
};
