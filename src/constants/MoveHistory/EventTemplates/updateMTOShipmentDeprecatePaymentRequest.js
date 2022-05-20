import o from 'constants/MoveHistory/UIDisplay/Operations';
import d from 'constants/MoveHistory/UIDisplay/DetailsTypes';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';
import { PAYMENT_REQUEST_STATUS_LABELS } from 'constants/paymentRequestStatus';

export default {
  action: a.UPDATE,
  eventName: o.updateMTOShipment,
  tableName: t.payment_requests,
  detailsType: d.STATUS,
  getEventNameDisplay: ({ oldValues }) => `Updated payment request ${oldValues?.payment_request_number}`,
  getStatusDetails: ({ changedValues }) => {
    const { status } = changedValues;
    return PAYMENT_REQUEST_STATUS_LABELS[status];
  },
};
