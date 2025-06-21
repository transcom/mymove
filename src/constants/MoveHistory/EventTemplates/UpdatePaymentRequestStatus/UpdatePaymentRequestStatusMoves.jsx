import React from 'react';

import o from 'constants/MoveHistory/UIDisplay/Operations';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';

export default {
  action: a.UPDATE,
  eventName: o.updatePaymentRequestStatus,
  tableName: t.moves,
  getEventNameDisplay: () => 'Updated move',
  getDetails: ({ changedValues }) => (
    <>
      <div>Payment Requests Addressed</div>
      {changedValues?.tio_payment_request_assigned_id !== undefined ? (
        <div>Task invoicing officer unassigned</div>
      ) : null}
    </>
  ),
};
