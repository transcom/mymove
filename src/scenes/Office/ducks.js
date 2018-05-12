import {
  LoadAccountingAPI,
  UpdateAccountingAPI,
  LoadMove,
  LoadOrders,
} from './api.js';
import * as ReduxHelpers from 'shared/ReduxHelpers';

// Types
const loadAccountingType = 'LOAD_ACCOUNTING';
const updateAccountingType = 'UPDATE_ACCOUNTING';
const loadMoveType = 'LOAD_MOVE';
const loadOrdersType = 'LOAD_ORDERS';

const LOAD_ACCOUNTING = ReduxHelpers.generateAsyncActionTypes(
  loadAccountingType,
);

const UPDATE_ACCOUNTING = ReduxHelpers.generateAsyncActionTypes(
  updateAccountingType,
);

const LOAD_MOVE = ReduxHelpers.generateAsyncActionTypes(loadMoveType);

const LOAD_ORDERS = ReduxHelpers.generateAsyncActionTypes(loadOrdersType);

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

export const loadOrders = ReduxHelpers.generateAsyncActionCreator(
  loadOrdersType,
  LoadOrders,
);

// Reducer
const initialState = {
  accountingIsLoading: false,
  accountingIsUpdating: false,
  moveIsLoading: false,
  ordersAreLoading: false,
  accountingHasLoadError: false,
  accountingHasLoadSuccess: null,
  accountingHasUpdateError: false,
  accountingHasUpdateSuccess: null,
  moveHasLoadError: false,
  moveHasLoadSuccess: null,
  ordersHaveLoadError: false,
  ordersHaveLoadSuccess: null,
};

export function officeReducer(state = initialState, action) {
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

    // Moves
    case LOAD_MOVE.start:
      return Object.assign({}, state, {
        moveIsLoading: true,
        moveHasLoadSuccess: false,
      });
    case LOAD_MOVE.success:
      return Object.assign({}, state, {
        moveIsLoading: false,
        officeMove: action.payload,
        moveHasLoadSuccess: true,
        moveHasLoadError: false,
      });
    case LOAD_MOVE.failure:
      return Object.assign({}, state, {
        moveIsLoading: false,
        officeMove: null,
        moveHasLoadSuccess: false,
        moveHasLoadError: true,
        error: action.error.message,
      });

    // ORDERS
    case LOAD_ORDERS.start:
      return Object.assign({}, state, {
        ordersAreLoading: true,
        ordersHaveLoadSuccess: false,
      });
    case LOAD_ORDERS.success:
      return Object.assign({}, state, {
        ordersAreLoading: false,
        officeOrders: action.payload,
        ordersHaveLoadSuccess: true,
        ordersHaveLoadError: false,
      });
    case LOAD_ORDERS.failure:
      return Object.assign({}, state, {
        ordersAreLoading: false,
        officeOrders: null,
        ordersHaveLoadSuccess: false,
        ordersHaveLoadError: true,
        error: action.error.message,
      });
    default:
      return state;
  }
}
