import React from 'react';

import { formatMoveHistoryFullAddress } from 'utils/formatters';
import a from 'constants/MoveHistory/Database/Actions';
import o from 'constants/MoveHistory/UIDisplay/Operations';
import t from 'constants/MoveHistory/Database/Tables';
import LabeledDetails from 'pages/Office/MoveHistory/LabeledDetails';
import ADDRESS_TYPE from 'constants/MoveHistory/Database/AddressTypes';
import { getMtoShipmentLabel } from 'utils/formatMtoShipment';

const formatChangedValues = (historyRecord) => {
  const { context, changedValues } = historyRecord;
  const address = formatMoveHistoryFullAddress(changedValues);

  const addressType = context.filter((contextObject) => contextObject.address_type)[0].address_type;
  const addressLabel = ADDRESS_TYPE[addressType];

  const newChangedValues = {
    [addressLabel]: address,
    ...changedValues,
    ...getMtoShipmentLabel(historyRecord),
  };

  return { ...historyRecord, changedValues: newChangedValues };
};

export default {
  action: a.INSERT,
  eventName: o.updateMTOShipment,
  tableName: t.addresses,
  getEventNameDisplay: () => 'Updated address',
  getDetails: (historyRecord) => <LabeledDetails historyRecord={formatChangedValues(historyRecord)} />,
};
