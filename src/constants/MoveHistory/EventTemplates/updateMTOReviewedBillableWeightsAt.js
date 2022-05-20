import detailsTypes from 'constants/MoveHistory/UIDisplay/DetailsTypes';
import operations from 'constants/MoveHistory/UIDisplay/Operations';
import actions from 'constants/MoveHistory/Database/Actions';
import tables from 'constants/MoveHistory/Database/Tables';

export default {
  action: actions.UPDATE,
  eventName: operations.updateMTOReviewedBillableWeightsAt,
  tableName: tables.moves,
  detailsType: detailsTypes.PLAIN_TEXT,
  getEventNameDisplay: () => 'Updated move',
  getDetailsPlainText: () => 'Reviewed weights',
};
