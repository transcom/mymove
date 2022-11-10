import React from 'react';

import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';
import LabeledDetails from 'pages/Office/MoveHistory/LabeledDetails';
import { formatMoveHistoryFullAddress } from 'utils/formatters';
import ADDRESS_TYPE from 'constants/MoveHistory/Database/AddressTypes';
import { getMtoShipmentLabel } from 'utils/formatMtoShipment';

const formatChangedValues = (historyRecord) => {
  const { context, changedValues, oldValues } = historyRecord;
  const address = formatMoveHistoryFullAddress(changedValues);

  const addressType = context.filter((contextObject) => contextObject.address_type)[0].address_type;
  const addressLabel = ADDRESS_TYPE[addressType];

  const newChangedValues = {
    street_address_1: oldValues.street_address_1,
    street_address_2: oldValues.street_address_2,
    city: oldValues.city,
    state: oldValues.state,
    postal_code: oldValues.postal_code,
    [addressLabel]: address,
    ...changedValues,
    ...getMtoShipmentLabel(historyRecord),
  };

  return { ...historyRecord, changedValues: newChangedValues };
};

export default {
  action: a.UPDATE,
  eventName: '*',
  tableName: t.addresses,
  getEventNameDisplay: () => 'Updated address',
  getDetails: (historyRecord) => <LabeledDetails historyRecord={formatChangedValues(historyRecord)} />,
};
