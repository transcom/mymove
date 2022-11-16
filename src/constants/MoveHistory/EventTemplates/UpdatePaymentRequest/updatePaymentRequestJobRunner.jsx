import React from 'react';

import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';
import { PAYMENT_REQUEST_STATUS_LABELS as p } from 'constants/paymentRequestStatus';
import LabeledDetails from 'pages/Office/MoveHistory/LabeledDetails';

const formatChangedValues = (historyRecord) => {
  const { changedValues } = historyRecord;

  const newChangedValues = {
    ...changedValues,
    status: p[changedValues.status],
  };

  return { ...historyRecord, changedValues: newChangedValues };
};

export default {
  action: a.UPDATE,
  tableName: t.payment_requests,
  getEventNameDisplay: ({ oldValues }) => <> Updated payment request {oldValues?.payment_request_number} </>,
  getDetails: (historyRecord) => <LabeledDetails historyRecord={formatChangedValues(historyRecord)} />,
};
