import React from 'react';

import o from 'constants/MoveHistory/UIDisplay/Operations';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';

export default {
  action: a.INSERT,
  eventName: o.updateMTOShipment,
  tableName: t.payment_requests,
  getEventNameDisplay: ({ changedValues }) => `Created payment request ${changedValues?.payment_request_number}`,
  getDetails: () => (
    <>
      <b>Status</b>: Pending
    </>
  ),
};
