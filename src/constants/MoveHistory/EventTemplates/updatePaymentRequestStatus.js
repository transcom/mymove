import o from 'constants/MoveHistory/UIDisplay/Operations';
import d from 'constants/MoveHistory/UIDisplay/DetailsTypes';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';

export default {
  action: a.UPDATE,
  eventName: o.updatePaymentRequestStatus,
  tableName: t.payment_requests,
  detailsType: d.PAYMENT,
  getEventNameDisplay: ({ oldValues, changedValues }) => {
    const paymentRequestNumber = oldValues?.payment_request_number ?? changedValues?.payment_request_number;
    return `Submitted payment request ${paymentRequestNumber}`;
  },
};
