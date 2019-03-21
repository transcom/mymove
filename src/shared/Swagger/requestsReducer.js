import { last, omit, startsWith } from 'lodash';

const initialState = {
  byID: {},
  errored: {},
  lastErrors: {},
};

export function requestsReducer(state = initialState, action) {
  if (startsWith(action.type, '@@swagger')) {
    const parts = action.type.split('/');
    switch (last(parts)) {
      case 'START':
        return Object.assign({}, state, {
          byID: {
            ...state.byID,
            [action.request.id]: action.request,
          },
        });
      case 'SUCCESS':
        return Object.assign({}, state, {
          byID: {
            ...state.byID,
            [action.request.id]: action.request,
          },
        });
      case 'FAILURE':
        return Object.assign({}, state, {
          byID: {
            ...state.byID,
            [action.request.id]: action.request,
          },
          errored: {
            ...state.errored,
            [action.request.id]: action.request,
          },
          lastErrors: {
            ...state.lastErrors,
            [action.label]: action.request,
          },
        });
      case 'RESET':
        return Object.assign({}, state, {
          lastErrors: omit(state.lastErrors, [action.label]),
          byID: {},
        });
      default:
        return state;
    }
  }
  return state;
}
