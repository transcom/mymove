import d from 'constants/MoveHistory/UIDisplay/DetailsTypes';
import t from 'constants/MoveHistory/Database/Tables';
import a from 'constants/MoveHistory/Database/Actions';
import o from 'constants/MoveHistory/UIDisplay/Operations';

export default {
  action: a.INSERT,
  eventName: o.createOrders,
  tableName: t.orders,
  detailsType: d.PLAIN_TEXT,
  getEventNameDisplay: () => 'Submitted orders',
  getDetailsPlainText: () => '-',
};
