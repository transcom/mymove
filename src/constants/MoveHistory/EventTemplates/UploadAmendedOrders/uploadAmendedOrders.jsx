import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';
import o from 'constants/MoveHistory/UIDisplay/Operations';

export default {
  action: a.INSERT,
  eventName: o.uploadAmendedOrders,
  tableName: t.user_uploads,
  getEventNameDisplay: () => 'Updated orders',
  getDetails: (historyRecord) => {
    return `Uploaded amended orders document ${historyRecord.context[0]?.filename}`;
  },
};
