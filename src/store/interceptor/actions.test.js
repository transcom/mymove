import { INTERCEPT_RESPONSE, interceptResponse } from './actions';

describe('interceptor actions', () => {
  it('interceptResponse returns an action object', () => {
    expect(interceptResponse(true)).toEqual({
      type: INTERCEPT_RESPONSE,
      hasError: true,
      traceId: '',
    });
    expect(interceptResponse(true, 'some-id')).toEqual({
      type: INTERCEPT_RESPONSE,
      hasError: true,
      traceId: 'some-id',
    });
  });
});
