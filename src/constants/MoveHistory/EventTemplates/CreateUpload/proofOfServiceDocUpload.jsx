import React from 'react';

import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';
import o from 'constants/MoveHistory/UIDisplay/Operations';

export default {
  action: a.INSERT,
  eventName: o.createUpload,
  tableName: t.proof_of_service_docs,
  getEventNameDisplay: () => 'Uploaded document',
  getDetails: ({ context }) => (
    <> Proof of service document uploaded for payment request {context[0].payment_request_number} </>
  ),
};
