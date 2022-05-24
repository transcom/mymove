import d from 'constants/MoveHistory/UIDisplay/DetailsTypes';
import o from 'constants/MoveHistory/UIDisplay/Operations';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';
import { PAYMENT_REQUEST_STATUS_LABELS } from 'constants/paymentRequestStatus';

export default {
  action: a.UPDATE,
  eventName: o.updateReweigh,
  tableName: t.payment_requests,
  detailsType: d.STATUS,
  getEventNameDisplay: ({ oldValues }) => `Updated payment request ${oldValues?.payment_request_number}`,
  getStatusDetails: ({ changedValues }) => {
    let status = '';
    if (changedValues.recalculation_of_payment_request_id) {
      status = 'Recalculated payment request';
    } else if (changedValues.status) {
      status = PAYMENT_REQUEST_STATUS_LABELS[changedValues.status];
    } else {
      status = 'Undefined status';
    }
    return status;
  },
};
