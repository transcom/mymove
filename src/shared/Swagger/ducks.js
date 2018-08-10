import * as helpers from 'shared/ReduxHelpers';
import { GetSpec, GetPublicSpec } from './api';
import { filter, last, sortBy, startsWith } from 'lodash';

const resource = 'SWAGGER';

export const actionsTypes = helpers.generateAsyncActionTypes(resource);

export const loadSchema = helpers.generateAsyncActionCreator(resource, GetSpec);

export const loadPublicSchema = helpers.generateAsyncActionCreator(
  resource,
  GetPublicSpec,
);

export const swaggerReducer =  helpers.generateAsyncReducer(resource, v => ({ spec: v }));

const initialState = {
  byID: {},
  errored: {},
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
          }
        });
      default:
        return state;
    }
  }
  return state;
}

export function lastRequest(state, label) {
  const requests = filter(state.requests.byID, function(value, key) {
    return value.label === label;
  });
  const sorted = sortBy(requests, ['start']);
  return last(sorted);
}