import React from 'react';

import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';

export default {
  action: a.INSERT,
  eventName: '*',
  tableName: t.service_members,
  getEventNameDisplay: () => 'Created Profile',
  getDetails: () => <> - </>,
};
