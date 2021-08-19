import { interceptResponse } from './actions';
import interceptorReducer, { initialState } from './reducer';

const SHORT_RESPONSE_GAP_MS = 100;
const LONG_RESPONSE_GAP_MS = 4000;

describe('interceptorReducer', () => {
  it('returns the initial state by default', () => {
    expect(interceptorReducer(undefined, undefined)).toEqual(initialState);
  });

  describe('INTERCEPT_RESPONSE', () => {
    it('updates state if an error was intercepted', () => {
      const oldState = {
        hasRecentError: false,
        timestamp: Date.now() - SHORT_RESPONSE_GAP_MS,
        traceId: '',
      };
      const traceId = 'some-trace-id';
      const newState = interceptorReducer(oldState, interceptResponse(true, traceId));
      expect(newState.hasRecentError).toEqual(true);
      expect(newState.timestamp - oldState.timestamp).toBeGreaterThanOrEqual(SHORT_RESPONSE_GAP_MS);
      expect(newState.traceId).toBe(traceId);
    });

    it('maintains hasRecentError value if a recent response was intercepted', () => {
      const oldState = {
        hasRecentError: true,
        timestamp: Date.now() - SHORT_RESPONSE_GAP_MS,
      };

      const newState = interceptorReducer(oldState, interceptResponse(false));
      expect(newState.hasRecentError).toEqual(true);
      expect(newState.timestamp - oldState.timestamp).toBeGreaterThanOrEqual(SHORT_RESPONSE_GAP_MS);
    });

    it('updates hasRecentError value if timestamp is stale and clears traceId', () => {
      const oldState = {
        hasRecentError: true,
        timestamp: Date.now() - LONG_RESPONSE_GAP_MS,
        traceId: 'some-trace-id',
      };

      const newState = interceptorReducer(oldState, interceptResponse(false));
      expect(newState.hasRecentError).toEqual(false);
      expect(newState.timestamp - oldState.timestamp).toBeGreaterThanOrEqual(LONG_RESPONSE_GAP_MS);
      expect(newState.traceId).toBe('');
    });

    it('maintains hasRecentError value if action value is false and state timestamp is stale', () => {
      const oldState = {
        hasRecentError: false,
        timestamp: Date.now() - LONG_RESPONSE_GAP_MS,
      };

      const newState = interceptorReducer(oldState, interceptResponse(false));
      expect(newState.hasRecentError).toEqual(false);
      expect(newState.timestamp - oldState.timestamp).toBeGreaterThanOrEqual(LONG_RESPONSE_GAP_MS);
    });
  });
});
