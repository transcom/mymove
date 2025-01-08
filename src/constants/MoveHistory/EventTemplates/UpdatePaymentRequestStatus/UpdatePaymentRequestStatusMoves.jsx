import React from 'react';

import o from 'constants/MoveHistory/UIDisplay/Operations';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';

export default {
  action: a.UPDATE,
  eventName: o.updatePaymentRequestStatus,
  tableName: t.moves,
  getEventNameDisplay: () => 'Updated move',
  getDetails: ({ changedValues }) => {
    if (changedValues?.tio_assigned_id !== undefined) return <> Task Invoicing Officer Unassigned </>;
    return <> - </>;
  },
};
