import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';
import o from 'constants/MoveHistory/UIDisplay/Operations';

export default {
  action: a.INSERT,
  eventName: o.createUpload,
  tableName: t.user_uploads,
  getEventNameDisplay: () => 'Updated orders',
  getDetails: (historyRecord) => {
    return `Uploaded orders document ${historyRecord.context[0]?.filename}`;
  },
};
