import React from 'react';

import a from 'constants/MoveHistory/Database/Actions';
import o from 'constants/MoveHistory/UIDisplay/Operations';
import t from 'constants/MoveHistory/Database/Tables';

export default {
  action: a.UPDATE,
  eventName: o.submitMoveForApproval,
  tableName: t.moves,
  getEventNameDisplay: () => 'Customer Signature',
  getDetails: () => <> Received customer signature </>,
};
