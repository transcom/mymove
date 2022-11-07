import React from 'react';

import o from 'constants/MoveHistory/UIDisplay/Operations';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';

export default {
  action: a.UPDATE,
  eventName: o.updateMTOStatusServiceCounselingCompleted,
  tableName: t.moves,
  getEventNameDisplay: () => 'Updated shipment',
  getDetailsPlainText: () => <> Counseling Completed </>,
};
