import d from 'constants/MoveHistory/UIDisplay/DetailsTypes';
import t from 'constants/MoveHistory/Database/Tables';
import { PAYMENT_REQUEST_STATUS_LABELS } from 'constants/paymentRequestStatus';

export default {
  action: 'UPDATE',
  eventName: '',
  tableName: t.payment_requests,
  detailsType: d.STATUS,
  getEventNameDisplay: ({ oldValues }) => `Updated payment request ${oldValues?.payment_request_number}`,
  getStatusDetails: ({ changedValues }) => {
    const { status } = changedValues;
    switch (status) {
      case 'SENT_TO_GEX':
        return 'Sent to GEX';
      case 'RECEIVED_BY_GEX':
        return 'Received';
      default:
        return PAYMENT_REQUEST_STATUS_LABELS[status];
    }
  },
};
