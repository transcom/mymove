import o from 'constants/MoveHistory/UIDisplay/Operations';
import d from 'constants/MoveHistory/UIDisplay/DetailsTypes';

export default {
  action: '*',
  eventName: o.submitMoveForApproval,
  tableName: '*',
  detailsType: d.PLAIN_TEXT,
  getEventNameDisplay: () => 'Submitted move',
  getDetailsPlainText: () => '-',
};
