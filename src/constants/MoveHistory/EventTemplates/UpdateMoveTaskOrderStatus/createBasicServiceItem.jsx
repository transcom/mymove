import a from 'constants/MoveHistory/Database/Actions';
import o from 'constants/MoveHistory/UIDisplay/Operations';
import d from 'constants/MoveHistory/UIDisplay/DetailsTypes';
import t from 'constants/MoveHistory/Database/Tables';

export default {
  action: a.INSERT,
  eventName: o.updateMoveTaskOrderStatus,
  tableName: t.mto_service_items,
  detailsType: d.PLAIN_TEXT,
  getEventNameDisplay: () => 'Approved service item',
  getDetailsPlainText: (historyRecord) => {
    return `${historyRecord.context[0]?.name}`;
  },
};
