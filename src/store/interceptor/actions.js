export const INTERCEPT_RESPONSE = 'INTERCEPT_RESPONSE';

export const interceptResponse = (hasError) => ({
  type: INTERCEPT_RESPONSE,
  hasError,
});
