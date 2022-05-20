import detailsTypes from 'constants/MoveHistory/UIDisplay/DetailsTypes';
import operations from 'constants/MoveHistory/UIDisplay/Operations';
import actions from 'constants/MoveHistory/Database/Actions';

export default {
  action: actions.UPDATE,
  eventName: operations.uploadAmendedOrders,
  tableName: 'orders',
  detailsType: detailsTypes.PLAIN_TEXT,
  getEventNameDisplay: () => 'Updated orders',
  getDetailsPlainText: () => '-',
};
