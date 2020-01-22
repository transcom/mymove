import { swaggerRequest } from 'shared/Swagger/request';
import { getGHCClient } from 'shared/Swagger/api';
import { get } from 'lodash';

const getPaymentRequestLabel = 'PaymentRequests.getPaymentRequest';
const getPaymentRequestListLabel = 'PaymentRequests.getPaymentRequestList';
const updatePaymentRequestLabel = 'PaymentRequests.updatePaymentRequest';

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

export function selectPaymentRequests(state) {
  const paymentRequests = get(state, 'entities.paymentRequests') || {};
  return Object.values(paymentRequests);
}

export function updatePaymentRequest(
  { paymentRequestID, status, rejectionReason = '' },
  label = updatePaymentRequestLabel,
) {
  const swaggerTag = 'paymentRequests.updatePaymentRequestStatus';
  return swaggerRequest(getGHCClient, swaggerTag, { paymentRequestID, body: { status, rejectionReason } }, { label });
}
