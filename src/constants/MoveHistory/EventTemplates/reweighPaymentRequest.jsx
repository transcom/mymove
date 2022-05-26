import React from 'react';

import d from 'constants/MoveHistory/UIDisplay/DetailsTypes';
import o from 'constants/MoveHistory/UIDisplay/Operations';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';
import { PAYMENT_REQUEST_STATUS_LABELS } from 'constants/paymentRequestStatus';

export default {
  action: a.UPDATE,
  eventName: o.updateReweigh,
  tableName: t.payment_requests,
  detailsType: d.CUSTOM,
  getEventNameDisplay: ({ oldValues }) => `Updated payment request ${oldValues?.payment_request_number}`,
  getStatusDetails: ({ changedValues }) => {
    let status = '';
    if (changedValues.recalculation_of_payment_request_id) {
      status = 'Recalculated payment requestssss';
    } else if (changedValues.status) {
      status = PAYMENT_REQUEST_STATUS_LABELS[changedValues.status];
    } else {
      status = 'Undefined status';
    }
    return status;
  },
  getCustomDetails: ({ changedValues }) => {
    let customTemplate = <div>Undefined</div>;
    if (changedValues.recalculation_of_payment_request_id) {
      customTemplate = <div>Recalculated payment requestss</div>;
    } else if (changedValues.status) {
      customTemplate = (
        <div>
          <b>Status: </b>
          {PAYMENT_REQUEST_STATUS_LABELS[changedValues.status]}
        </div>
      );
    }
    return customTemplate;
  },
};
