import o from 'constants/MoveHistory/UIDisplay/Operations';
import d from 'constants/MoveHistory/UIDisplay/DetailsTypes';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';

export default {
  action: a.INSERT,
  eventName: o.createPaymentRequest,
  tableName: t.payment_requests,
  detailsType: d.PLAIN_TEXT,
  getEventNameDisplay: ({ changedValues }) => `Submitted payment request ${changedValues?.payment_request_number}`,
  getDetailsPlainText: ({ context }) => {
    return context
      .reduce((serviceItemsString, contextItem) => {
        return `${serviceItemsString}, ${contextItem.name}`;
      }, '')
      .slice(2);
  },
};
