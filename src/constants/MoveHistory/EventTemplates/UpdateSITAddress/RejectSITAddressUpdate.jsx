import React from 'react';

import o from 'constants/MoveHistory/UIDisplay/Operations';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';
import LabeledDetails from 'pages/Office/MoveHistory/LabeledDetails';
import { formatSITData } from 'utils/formatSITData';
import { getMtoShipmentLabel } from 'utils/formatMtoShipment';

// this allows office users to reject destination SIT updates requests made by the prime
const formatChangedValues = (historyRecord) => {
  const newChangedValues = {
    ...historyRecord.changedValues,
    ...formatSITData(historyRecord),
    ...getMtoShipmentLabel(historyRecord),
  };

  return { ...historyRecord, changedValues: newChangedValues };
};

export default {
  action: a.UPDATE,
  eventName: o.rejectSITAddressUpdate,
  tableName: t.sit_address_updates,
  getEventNameDisplay: () => {
    return <div>Updated service item</div>;
  },
  getDetails: (historyRecord) => {
    return <LabeledDetails historyRecord={formatChangedValues(historyRecord)} />;
  },
};
