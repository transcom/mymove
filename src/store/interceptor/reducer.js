import { INTERCEPT_RESPONSE } from './actions';

const ACTION_RESPONSE_INTERVAL_MS = 3000;

export const initialState = {
  hasRecentError: false,
  timestamp: 0,
  traceId: '',
};

const interceptorReducer = (state = initialState, action = {}) => {
  switch (action?.type) {
    case INTERCEPT_RESPONSE: {
      const timestamp = Date.now();

      if (action.hasError) {
        return {
          ...state,
          hasRecentError: true,
          timestamp,
          traceId: action.traceId,
        };
      }

      if (timestamp > state.timestamp + ACTION_RESPONSE_INTERVAL_MS) {
        return {
          ...state,
          hasRecentError: false,
          timestamp,
          traceId: '',
        };
      }

      return {
        ...state,
        timestamp,
      };
    }

    default:
      return state;
  }
};

export default interceptorReducer;
