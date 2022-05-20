import o from 'constants/MoveHistory/UIDisplay/Operations';
import d from 'constants/MoveHistory/UIDisplay/DetailsTypes';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';

export default {
  action: a.UPDATE,
  eventName: o.updateMoveTaskOrder,
  tableName: t.moves,
  detailsType: d.LABELED,
  getEventNameDisplay: () => 'Updated move',
};
