import * as helpers from 'shared/ReduxHelpers';
import { GetSpec } from './api';
import {last, startsWith } from 'lodash';

const resource = 'SWAGGER';

export const actionsTypes = helpers.generateAsyncActionTypes(resource);

export const loadSchema = helpers.generateAsyncActionCreator(resource, GetSpec);

export default helpers.generateAsyncReducer(resource, v => ({ spec: v }));

const initialState = {
  requests: {}
};

export function requestReducer(state = initialState, action) {
  if (startsWith(action.type, '@@swagger') {
    const parts = action.type.split('/');
    switch (last(parts)) {
      case 'START':
        return Object.assign({}, state, {
          moveIsLoading: true,
          moveHasLoadSuccess: false,
        });
      case 'SUCCESS':

      case 'FAILURE':

      default:
    }
  }
}
