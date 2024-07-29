import React from 'react';

import o from 'constants/MoveHistory/UIDisplay/Operations';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';
import LabeledDetails from 'pages/Office/MoveHistory/LabeledDetails';
import { getMtoShipmentLabel } from 'utils/formatMtoShipment';

const formatChangedValues = (historyRecord) => {
  const newChangedValues = {
    ...historyRecord.changedValues,
    ...getMtoShipmentLabel(historyRecord),
  };

  return { ...historyRecord, changedValues: newChangedValues };
};

export default {
  action: a.UPDATE,
  eventName: o.updateMTOPostCounselingInformation,
  tableName: t.ppm_shipments,
  getEventNameDisplay: () => 'Updated shipment',
  getDetails: (historyRecord) => <LabeledDetails historyRecord={formatChangedValues(historyRecord)} />,
};

// export default {
//   action: a.UPDATE,
//   eventName: o.updateMTOPostCounselingInformation,
//   tableName: t.ppm_shipments,
//   getEventNameDisplay: () => 'Updated Shipment',
//   getDetails: () => <> Prime Counseling Completed for PPM Shipment </>,
// };
