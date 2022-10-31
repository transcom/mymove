import o from 'constants/MoveHistory/UIDisplay/Operations';
import d from 'constants/MoveHistory/UIDisplay/DetailsTypes';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';

export default {
  action: a.UPDATE,
  eventName: o.updateMoveTaskOrderStatus,
  tableName: t.moves,
  detailsType: d.PLAIN_TEXT,
  getEventNameDisplay: (historyRecord) => {
    return historyRecord.changedValues?.available_to_prime_at ? 'Approved move' : 'Move status updated';
  },
  getDetailsPlainText: (historyRecord) => {
    return historyRecord.changedValues?.available_to_prime_at ? 'Created Move Task Order (MTO)' : '-';
  },
};
