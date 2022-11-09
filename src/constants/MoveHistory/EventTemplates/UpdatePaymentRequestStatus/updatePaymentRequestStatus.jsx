import React from 'react';

import o from 'constants/MoveHistory/UIDisplay/Operations';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';
import PaymentDetails from 'pages/Office/MoveHistory/PaymentDetails';

export default {
  action: a.UPDATE,
  eventName: o.updatePaymentRequestStatus,
  tableName: t.payment_requests,
  getEventNameDisplay: ({ oldValues, changedValues }) => {
    const paymentRequestNumber = oldValues?.payment_request_number ?? changedValues?.payment_request_number;

    return <> Submitted payment request {paymentRequestNumber} </>;
  },
  getDetails: ({ context }) => <PaymentDetails context={context} />,
};
