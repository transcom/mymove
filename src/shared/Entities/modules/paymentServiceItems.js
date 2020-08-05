import { swaggerRequest } from 'shared/Swagger/request';
import { getGHCClient } from 'shared/Swagger/api';

const updatePaymentServiceItemOperation = 'paymentServiceItem.updatePaymentServiceItemStatus';
const paymentServiceItemSchemaKey = 'paymentServiceItem';
export function patchPaymentServiceItemStatus(
  moveTaskOrderID,
  paymentServiceItemID,
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
      paymentServiceItemID,
      'If-Match': ifMatchEtag,
      body: { status: paymentServiceItemStatus, rejectionReason },
    },
    { label, schemaKey },
  );
}

export default patchPaymentServiceItemStatus;
