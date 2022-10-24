import React from 'react';

import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';
import o from 'constants/MoveHistory/UIDisplay/Operations';
import LabeledDetails from 'pages/Office/MoveHistory/LabeledDetails';

export default {
  action: a.UPDATE,
  eventName: o.updateServiceMemberBackupContact,
  tableName: t.backup_contacts,
  getEventNameDisplay: () => 'Updated profile',
  getDetails: (historyRecord) => <LabeledDetails historyRecord={historyRecord} />,
};
