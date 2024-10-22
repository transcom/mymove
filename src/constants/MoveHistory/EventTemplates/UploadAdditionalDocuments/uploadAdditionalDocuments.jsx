import React from 'react';

import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';
import o from 'constants/MoveHistory/UIDisplay/Operations';

export default {
  action: a.UPDATE,
  eventName: o.uploadAdditionalDocuments,
  tableName: t.moves,
  getEventNameDisplay: () => 'Updated move',
  getDetails: () => <>Uploaded additional document</>,
};
