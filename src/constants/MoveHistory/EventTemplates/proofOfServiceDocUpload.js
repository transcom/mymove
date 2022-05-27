import a from 'constants/MoveHistory/Database/Actions';
import d from 'constants/MoveHistory/UIDisplay/DetailsTypes';
import t from 'constants/MoveHistory/Database/Tables';

export default {
  action: a.INSERT,
  eventName: 'createUpload',
  tableName: t.proof_of_service_docs,
  detailsType: d.PLAIN_TEXT,
  getEventNameDisplay: () => 'Uploaded document',
  getDetailsPlainText: ({ context }) => {
    return `Proof of service document uploaded for payment request ${context[0].payment_request_number}`;
  },
};
