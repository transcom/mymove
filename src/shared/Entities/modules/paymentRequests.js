import { swaggerRequest } from 'shared/Swagger/request';
import { getGHCClient } from 'shared/Swagger/api';
import { get } from 'lodash';

const getPaymentRequestLabel = 'PaymentRequests.getPaymentRequest';
const getPaymentRequestListLabel = 'PaymentRequests.getPaymentRequestList';

export function getPaymentRequest(paymentRequestID, label = getPaymentRequestLabel) {
  const swaggerTag = 'paymentRequests.getPaymentRequest';
  return swaggerRequest(getGHCClient, swaggerTag, { paymentRequestID }, { label });
}

export function getPaymentRequestList(label = getPaymentRequestListLabel) {
  const swaggerTag = 'paymentRequests.listPaymentRequests';
  return swaggerRequest(getGHCClient, swaggerTag, {}, { label });
}

export function selectPaymentRequest(state, id) {
  return get(state, `entities.paymentRequests.${id}`) || {};
}
