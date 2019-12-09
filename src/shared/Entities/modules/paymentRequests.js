import { swaggerRequest } from 'shared/Swagger/request';
import { getGHCClient } from 'shared/Swagger/api';

const getPaymentRequestLabel = 'PaymentRequests.getPaymentRequest';

export function getPaymentRequest(paymentRequestID, label = getPaymentRequestLabel) {
  const swaggerTag = 'paymentRequests.getPaymentRequest';
  return swaggerRequest(getGHCClient, swaggerTag, { paymentRequestID }, { label });
}
