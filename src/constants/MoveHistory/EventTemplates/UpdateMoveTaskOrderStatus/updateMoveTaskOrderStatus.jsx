import React from 'react';

import o from 'constants/MoveHistory/UIDisplay/Operations';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';

export default {
  action: a.UPDATE,
  eventName: o.updateMoveTaskOrderStatus,
  tableName: t.moves,
  getEventNameDisplay: ({ changedValues }) => {
    return changedValues?.available_to_prime_at ? <> Approved move </> : <> Move status updated </>;
  },
  getDetails: ({ changedValues }) => {
    return changedValues?.available_to_prime_at ? <> Created Move Task Order (MTO) </> : <> - </>;
  },
};
