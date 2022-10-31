import React from 'react';

import a from 'constants/MoveHistory/Database/Actions';
import o from 'constants/MoveHistory/UIDisplay/Operations';
import t from 'constants/MoveHistory/Database/Tables';

export default {
  action: a.INSERT,
  eventName: o.updateMoveTaskOrderStatus,
  tableName: t.mto_service_items,
  getEventNameDisplay: () => 'Approved service item',
  getDetails: ({ context }) => <> {context[0]?.name} </>,
};
