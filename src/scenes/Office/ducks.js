import { LoadAccountingAPI, UpdateAccountingAPI } from './api.js';
import * as ReduxHelpers from 'shared/ReduxHelpers';

// Types
const loadAccountingType = 'LOAD_ACCOUNTING';
const updateAccountingType = 'UPDATE_ACCOUNTING';

const LOAD_ACCOUNTING = ReduxHelpers.generateAsyncActionTypes(
  loadAccountingType,
);

const UPDATE_ACCOUNTING = ReduxHelpers.generateAsyncActionTypes(
  updateAccountingType,
);

export const loadAccounting = ReduxHelpers.generateAsyncActionCreator(
  loadAccountingType,
  LoadAccountingAPI,
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
    case LOAD_ACCOUNTING.start:
      return Object.assign({}, state, {
        hasSubmitSuccess: false,
      });
    case LOAD_ACCOUNTING.success:
      return Object.assign({}, state, {
        accounting: action.payload,
        hasSubmitSuccess: true,
        hasSubmitError: false,
      });
    case LOAD_ACCOUNTING.failure:
      return Object.assign({}, state, {
        accounting: null,
        hasSubmitSuccess: false,
        hasSubmitError: true,
        error: action.error.message,
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
