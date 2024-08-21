import React from 'react';

import o from 'constants/MoveHistory/UIDisplay/Operations';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';

export default {
  action: a.UPDATE,
  eventName: o.updateMTOPostCounselingInformation,
  tableName: t.moves,
  getEventNameDisplay: () => 'Updated move',
  getDetails: () => <> Prime Counseling Completed </>,
};

export const ppm = {
  action: a.UPDATE,
  eventName: o.updateMTOPostCounselingInformation,
  tableName: t.ppm_shipments,
  getEventNameDisplay: () => 'Updated Shipment',
  getDetails: () => <> Prime Counseling Completed for PPM Shipment </>,
};
