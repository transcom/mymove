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
  isLoading: false,
  isUpdating: false,
  hasLoadError: false,
  hasLoadSuccess: null,
  hasUpdateError: false,
  hasUpdateSuccess: null,
};

export function officeAccountingReducer(state = initialState, action) {
  switch (action.type) {
    case LOAD_ACCOUNTING.start:
      return Object.assign({}, state, {
        isLoading: true,
        hasLoadSuccess: false,
      });
    case LOAD_ACCOUNTING.success:
      return Object.assign({}, state, {
        isLoading: false,
        accounting: action.payload,
        hasLoadSuccess: true,
        hasLoadError: false,
      });
    case LOAD_ACCOUNTING.failure:
      return Object.assign({}, state, {
        isLoading: false,
        accounting: null,
        hasLoadSuccess: false,
        hasLoadError: true,
        error: action.error.message,
      });

    case UPDATE_ACCOUNTING.start:
      return Object.assign({}, state, {
        isUpdating: true,
        hasUpdateSuccess: false,
      });
    case UPDATE_ACCOUNTING.success:
      return Object.assign({}, state, {
        isUpdating: false,
        accounting: action.payload,
        hasUpdateSuccess: true,
        hasUpdateError: false,
      });
    case UPDATE_ACCOUNTING.failure:
      return Object.assign({}, state, {
        isUpdating: false,
        hasUpdateSuccess: false,
        hasUpdateError: true,
        error: action.error.message,
      });

    default:
      return state;
  }
}
