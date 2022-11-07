import React from 'react';

import operations from 'constants/MoveHistory/UIDisplay/Operations';
import actions from 'constants/MoveHistory/Database/Actions';
import tables from 'constants/MoveHistory/Database/Tables';

export default {
  action: actions.UPDATE,
  eventName: operations.updateMTOReviewedBillableWeightsAt,
  tableName: tables.moves,
  getEventNameDisplay: () => 'Updated move',
  getDetails: () => <> Reviewed weights </>,
};
