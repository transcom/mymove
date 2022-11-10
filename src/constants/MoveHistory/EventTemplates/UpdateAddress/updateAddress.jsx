import React from 'react';

import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';
import o from 'constants/MoveHistory/UIDisplay/Operations';
import LabeledDetails from 'pages/Office/MoveHistory/LabeledDetails';
import { formatMoveHistoryFullAddress } from 'utils/formatters';
import ADDRESS_TYPE from 'constants/MoveHistory/Database/AddressTypes';
import { getMtoShipmentLabel } from 'utils/formatMtoShipment';

const formatChangedValues = (historyRecord) => {
  const { context, changedValues, oldValues } = historyRecord;
  // order is important here, please keep oldValues first in the addressValues object
  const addressValues = {
    ...oldValues,
    ...changedValues,
  };
  const address = formatMoveHistoryFullAddress(addressValues);

  const addressType = context.filter((contextObject) => contextObject.address_type)[0].address_type;
  const addressLabel = ADDRESS_TYPE[addressType];

  const newChangedValues = {
    street_address_1: oldValues.street_address_1,
    street_address_2: oldValues.street_address_2,
    city: oldValues.city,
    state: oldValues.state,
    postal_code: oldValues.postal_code,
    [addressLabel]: address,
    ...getMtoShipmentLabel(historyRecord),
    ...changedValues,
  };

  return { ...historyRecord, changedValues: newChangedValues };
};

export default {
  action: a.UPDATE,
  eventName: o.patchServiceMember,
  tableName: t.addresses,
  getEventNameDisplay: () => 'Updated address',
  getDetails: (historyRecord) => <LabeledDetails historyRecord={formatChangedValues(historyRecord)} />,
};
