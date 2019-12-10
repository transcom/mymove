import { swaggerRequest } from 'shared/Swagger/request';
import { getGHCClient } from 'shared/Swagger/api';
import { get } from 'lodash';

const getPaymentRequestLabel = 'PaymentRequests.getPaymentRequest';

export function getPaymentRequest(paymentRequestID, label = getPaymentRequestLabel) {
  const swaggerTag = 'paymentRequests.getPaymentRequest';
  return swaggerRequest(getGHCClient, swaggerTag, { paymentRequestID }, { label });
}

export function selectPaymentRequest(state, id) {
  return get(state, `entities.paymentRequests.${id}`) || {};
}
