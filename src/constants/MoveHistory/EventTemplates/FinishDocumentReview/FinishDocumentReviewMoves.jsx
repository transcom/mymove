import React from 'react';

import o from 'constants/MoveHistory/UIDisplay/Operations';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';

export default {
  action: a.UPDATE,
  eventName: o.finishDocumentReview,
  tableName: t.moves,
  getEventNameDisplay: () => 'Updated move',
  getDetails: ({ changedValues }) => (
    <>
      <div>PPM Closeout Complete</div>
      {changedValues?.sc_closeout_assigned_id !== undefined ? <div>Closeout Counselor Unassigned</div> : null}
    </>
  ),
};
