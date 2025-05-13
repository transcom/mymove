import React from 'react';

import o from 'constants/MoveHistory/UIDisplay/Operations';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';

export default {
  action: a.UPDATE,
  eventName: o.deleteAssignedOfficeUser,
  tableName: t.moves,
  getEventNameDisplay: () => 'Updated move',
  getDetails: ({ changedValues }) => {
    if (changedValues.sc_counseling_assigned_id === null) return <>Counselor unassigned</>;
    if (changedValues.sc_closeout_assigned_id === null) return <>Closeout counselor unassigned</>;
    if (changedValues.too_task_order_assigned_id === null) return <>Task ordering officer unassigned</>;
    if (changedValues.too_destination_assigned_id === null) return <>Task ordering officer unassigned</>;
    if (changedValues.tio_assigned_id === null) return <>Task invoicing officer unassigned</>;
    return <>Unassigned</>;
  },
};
