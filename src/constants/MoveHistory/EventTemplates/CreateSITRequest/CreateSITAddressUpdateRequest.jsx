import React from 'react';

import o from 'constants/MoveHistory/UIDisplay/Operations';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';
import { getMtoShipmentLabel } from 'utils/formatMtoShipment';
import LabeledDetails from 'pages/Office/MoveHistory/LabeledDetails';
import { formatSITData } from 'utils/formatSITdata';

// this is for the prime to request an update to a destination SIT address
const formatChangedValues = (historyRecord) => {
  const newChangedValues = {
    ...formatSITData(historyRecord),
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
