import { LoadAccountingAPI, UpdateAccountingAPI, LoadMove } from './api.js';
import * as ReduxHelpers from 'shared/ReduxHelpers';

// Types
const loadAccountingType = 'LOAD_ACCOUNTING';
const updateAccountingType = 'UPDATE_ACCOUNTING';
const loadMoveType = 'LOAD_MOVE';

const LOAD_ACCOUNTING = ReduxHelpers.generateAsyncActionTypes(
  loadAccountingType,
);

const UPDATE_ACCOUNTING = ReduxHelpers.generateAsyncActionTypes(
  updateAccountingType,
);

const LOAD_MOVE = ReduxHelpers.generateAsyncActionTypes(loadMoveType);

export const loadAccounting = ReduxHelpers.generateAsyncActionCreator(
  loadAccountingType,
  LoadAccountingAPI,
);

export const updateAccounting = ReduxHelpers.generateAsyncActionCreator(
  updateAccountingType,
  UpdateAccountingAPI,
);

export const loadMove = ReduxHelpers.generateAsyncActionCreator(
  loadMoveType,
  LoadMove,
);

// Reducer
const initialState = {
  accountingIsLoading: false,
  accountingIsUpdating: false,
  moveIsLoading: false,
  accountingHasLoadError: false,
  accountingHasLoadSuccess: null,
  accountingHasUpdateError: false,
  accountingHasUpdateSuccess: null,
  moveHasLoadError: false,
  moveHasLoadSuccess: null,
};

export function officeAccountingReducer(state = initialState, action) {
  switch (action.type) {
    case LOAD_ACCOUNTING.start:
      return Object.assign({}, state, {
        accountingIsLoading: true,
        accountingHasLoadSuccess: false,
      });
    case LOAD_ACCOUNTING.success:
      return Object.assign({}, state, {
        accountingIsLoading: false,
        accounting: action.payload,
        accountingHasLoadSuccess: true,
        accountingHasLoadError: false,
      });
    case LOAD_ACCOUNTING.failure:
      return Object.assign({}, state, {
        accountingIsLoading: false,
        accounting: null,
        accountingHasLoadSuccess: false,
        accountingHasLoadError: true,
        error: action.error.message,
      });

    case UPDATE_ACCOUNTING.start:
      return Object.assign({}, state, {
        accountingIsUpdating: true,
        accountingHasUpdateSuccess: false,
      });
    case UPDATE_ACCOUNTING.success:
      return Object.assign({}, state, {
        accountingIsUpdating: false,
        accounting: action.payload,
        accountingHasUpdateSuccess: true,
        accountingHasUpdateError: false,
      });
    case UPDATE_ACCOUNTING.failure:
      return Object.assign({}, state, {
        accountingIsUpdating: false,
        accountingHasUpdateSuccess: false,
        accountingHasUpdateError: true,
        error: action.error.message,
      });
    case LOAD_MOVE.start:
      return Object.assign({}, state, {
        moveIsLoading: true,
        moveHasLoadSuccess: false,
      });
    case LOAD_MOVE.success:
      return Object.assign({}, state, {
        moveIsLoading: false,
        move: action.payload,
        moveHasLoadSuccess: true,
        moveHasLoadError: false,
      });
    case LOAD_MOVE.failure:
      return Object.assign({}, state, {
        moveIsLoading: false,
        move: null,
        moveHasLoadSuccess: false,
        moveHasLoadError: true,
        error: action.error.message,
      });
    default:
      return state;
  }
}
