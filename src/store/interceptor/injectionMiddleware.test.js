import { INTERCEPT_RESPONSE } from './actions';
import { interceptorInjectionMiddleware, interceptInjection } from './injectionMiddleware';

describe('injection middleware', () => {
  it('safely performs a no-op if called before a dispatch reference is found', () => {
    expect(() => {
      interceptInjection({
        type: INTERCEPT_RESPONSE,
        hasError: true,
      });
    }).not.toThrow();
  });

  it('acts as true middleware and just forwards the action to the next listener', () => {
    const dispatch = jest.fn();
    const next = jest.fn();

    interceptorInjectionMiddleware({ dispatch })(next)({
      type: 'test_event',
    });
    expect(next).toHaveBeenCalledWith({ type: 'test_event' });
  });

  it('maintains a closed over reference to the dispatch function', () => {
    const dispatch = jest.fn();
    const next = jest.fn();

    interceptorInjectionMiddleware({ dispatch })(next)({
      type: 'test_event',
    });

    interceptInjection({
      type: INTERCEPT_RESPONSE,
      hasError: false,
    });

    expect(dispatch).toHaveBeenCalledWith({
      type: INTERCEPT_RESPONSE,
      hasError: false,
    });
  });
});
