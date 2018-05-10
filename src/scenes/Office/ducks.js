import { GetAccountingAPI, UpdateAccountingAPI } from './api.js';
import * as ReduxHelpers from 'shared/ReduxHelpers';

// Types
const getAccountingType = 'GET_ACCOUNTING';
const updateAccountingType = 'UPDATE_ACCOUNTING';

const GET_ACCOUNTING = ReduxHelpers.generateAsyncActionTypes(getAccountingType);

const UPDATE_ACCOUNTING = ReduxHelpers.generateAsyncActionTypes(
  updateAccountingType,
);

export const getAccounting = ReduxHelpers.generateAsyncActionCreator(
  getAccountingType,
  GetAccountingAPI,
);

export const updateAccounting = ReduxHelpers.generateAsyncActionCreator(
  updateAccountingType,
  UpdateAccountingAPI,
);

// Reducer
const initialState = {
  hasSubmitError: false,
  hasSubmitSuccess: false,
  updateAccountingSuccess: false,
};

export function accountingReducer(state = initialState, action) {
  switch (action.type) {
    case GET_ACCOUNTING.start:
      return Object.assign({}, state, {
        hasSubmitSuccess: false,
      });
    case GET_ACCOUNTING.success:
      return Object.assign({}, state, {
        accounting: action.payload,
        hasSubmitSuccess: true,
        hasSubmitError: false,
      });
    case GET_ACCOUNTING.failure:
      return Object.assign({}, state, {
        accounting: null,
        hasSubmitSuccess: false,
        hasSubmitError: true,
        error: action.error,
      });

    case UPDATE_ACCOUNTING.start:
      return Object.assign({}, state, {
        updateAccountingSuccess: false,
      });
    case UPDATE_ACCOUNTING.success:
      return Object.assign({}, state, {
        accounting: action.payload,
        updateAccountingSuccess: true,
        updateAccountingError: false,
      });
    case UPDATE_ACCOUNTING.failure:
      return Object.assign({}, state, {
        updateAccountingSuccess: false,
        updateAccountingError: true,
        error: action.error,
      });

    default:
      return state;
  }
}
