import React from 'react';

import o from 'constants/MoveHistory/UIDisplay/Operations';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';
import { getMtoShipmentLabel } from 'utils/formatMtoShipment';
import LabeledDetails from 'pages/Office/MoveHistory/LabeledDetails';
import { formatMoveHistoryFullAddress } from 'utils/formatters';

const formatChangedValues = (historyRecord) => {
  const final = formatMoveHistoryFullAddress(JSON.parse(historyRecord.context[0].sit_destination_final_address));
  const initial = formatMoveHistoryFullAddress(JSON.parse(historyRecord.context[0].sit_destination_initial_address));

  const newChangedValues = {
    sit_destination_address_final: final,
    sit_destination_address_initial: initial,
    ...getMtoShipmentLabel(historyRecord),
    ...historyRecord.changedValues,
  };

  return { ...historyRecord, changedValues: newChangedValues };
};

export default {
  action: a.INSERT,
  eventName: o.createSITAddressUpdateRequest,
  tableName: t.sit_address_updates,
  getEventNameDisplay: () => {
    return <div>Updated service item</div>;
  },
  getDetails: (historyRecord) => {
    return <LabeledDetails historyRecord={formatChangedValues(historyRecord)} />;
  },
};
