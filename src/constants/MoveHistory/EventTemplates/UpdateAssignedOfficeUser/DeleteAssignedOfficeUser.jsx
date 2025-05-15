import React from 'react';

import o from 'constants/MoveHistory/UIDisplay/Operations';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';
import { MOVE_STATUSES } from 'shared/constants';

export default {
  action: a.UPDATE,
  eventName: o.deleteAssignedOfficeUser,
  tableName: t.moves,
  getEventNameDisplay: () => 'Move assignment updated',
  getDetails: ({ changedValues, oldValues }) => {
    if (changedValues.sc_assigned_id === null && oldValues?.status === MOVE_STATUSES.NEEDS_SERVICE_COUNSELING)
      return <>Counselor unassigned</>;
    if (changedValues.sc_assigned_id === null && oldValues?.status !== MOVE_STATUSES.NEEDS_SERVICE_COUNSELING)
      return <>Closeout counselor unassigned</>;
    if (changedValues.too_assigned_id === null || changedValues.too_destination_assigned_id === null)
      return <>Task ordering officer unassigned</>;
    if (changedValues.tio_assigned_id === null) return <>Task invoicing officer unassigned</>;
    return <>Unassigned</>;
  },
};
