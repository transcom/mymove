import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';
import o from 'constants/MoveHistory/UIDisplay/Operations';

export default {
  action: a.UPDATE,
  eventName: o.deleteUpload,
  tableName: t.user_uploads,
  getEventNameDisplay: (historyRecord) => {
    switch (historyRecord.context[0]?.upload_type) {
      // While it does look strange that both cases return the same event name,
      // It makes sense to leave it as a switch case to be flexible for anticipated
      // future user upload events
      case 'orders':
      case 'amendedOrders':
        return 'Updated orders';
      default:
        return 'Updated move';
    }
  },
  getDetails: (historyRecord) => {
    switch (historyRecord.context[0]?.upload_type) {
      case 'orders':
        return `Deleted orders document ${historyRecord.context[0]?.filename}`;
      case 'amendedOrders':
        return `Deleted amended orders document ${historyRecord.context[0]?.filename}`;
      default:
        return '-';
    }
  },
};
