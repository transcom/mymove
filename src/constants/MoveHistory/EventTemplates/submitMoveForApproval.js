import a from 'constants/MoveHistory/Database/Actions';
import o from 'constants/MoveHistory/UIDisplay/Operations';
import t from 'constants/MoveHistory/Database/Tables';
import d from 'constants/MoveHistory/UIDisplay/DetailsTypes';

export default {
  action: a.UPDATE,
  eventName: o.submitMoveForApproval,
  tableName: t.moves,
  detailsType: d.PLAIN_TEXT,
  getEventNameDisplay: () => 'Submitted move',
  getDetails: () => 'Received customer signature',
};
