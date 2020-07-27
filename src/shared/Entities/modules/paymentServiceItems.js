import { swaggerRequest } from 'shared/Swagger/request';
import { getGHCClient } from 'shared/Swagger/api';

const updatePaymentServiceItemOperation = 'paymentServiceItem.updatePaymentServiceItemStatus';
const paymentServiceItemSchemaKey = 'paymentServiceItem';
export function patchPaymentServiceItemStatus(
  moveTaskOrderID,
  paymentRequestID,
  paymentServiceItemStatus,
  ifMatchEtag,
  rejectionReason,
  label = updatePaymentServiceItemOperation,
  schemaKey = paymentServiceItemSchemaKey,
) {
  return swaggerRequest(
    getGHCClient,
    updatePaymentServiceItemOperation,
    {
      moveTaskOrderID,
      paymentRequestID,
      'If-match': ifMatchEtag,
      body: { status: paymentServiceItemStatus, rejectionReason },
    },
    { label, schemaKey },
  );
}
