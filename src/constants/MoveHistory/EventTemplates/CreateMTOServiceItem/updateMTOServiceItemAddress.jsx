import React from 'react';

import { formatMoveHistoryFullAddress } from 'utils/formatters';
import a from 'constants/MoveHistory/Database/Actions';
import o from 'constants/MoveHistory/UIDisplay/Operations';
import t from 'constants/MoveHistory/Database/Tables';
import AddressTypes from 'constants/MoveHistory/Database/AddressTypes';
import LabeledDetails from 'pages/Office/MoveHistory/LabeledDetails';

const formatChangedValues = (historyRecord) => {
  const { oldValues, changedValues, context } = historyRecord;
  const address = formatMoveHistoryFullAddress({ ...oldValues, ...changedValues });

  const addressType = context.filter((contextObject) => contextObject.address_type)[0].address_type;
  const addressLabel = AddressTypes[addressType];

  const newChangedValues = {
    shipment_type: context[0]?.shipment_type,
    shipment_id_display: context[0]?.shipment_id_abbr.toUpperCase(),
    ...changedValues,
  };

  newChangedValues[addressLabel] = address;
  return { ...historyRecord, changedValues: newChangedValues };
};

export default {
  action: a.UPDATE,
  eventName: o.createMTOServiceItem,
  tableName: t.addresses,
  getEventNameDisplay: () => 'Updated service item request',
  getDetails: (historyRecord) => <LabeledDetails historyRecord={formatChangedValues(historyRecord)} />,
};
