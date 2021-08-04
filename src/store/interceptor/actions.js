export const INTERCEPT_RESPONSE = 'INTERCEPT_RESPONSE';

export const interceptResponse = (hasError, traceId = '') => ({
  type: INTERCEPT_RESPONSE,
  hasError,
  traceId,
});
