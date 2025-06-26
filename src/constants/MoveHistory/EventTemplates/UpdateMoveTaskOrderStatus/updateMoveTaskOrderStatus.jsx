import React from 'react';

import o from 'constants/MoveHistory/UIDisplay/Operations';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';

export default {
  action: a.UPDATE,
  eventName: o.updateMoveTaskOrderStatus,
  tableName: t.moves,
  getEventNameDisplay: () => 'Approved move',
  getDetails: ({ changedValues }) => {
    return (
      <>
        <div>Created Move Task Order (MTO)</div>
        {(changedValues?.too_assigned_id !== undefined || changedValues?.too_destination_assigned_id !== undefined) && (
          <div>Task ordering officer unassigned</div>
        )}
      </>
    );
  },
};
