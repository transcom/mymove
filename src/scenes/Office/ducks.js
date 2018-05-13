import {
  LoadAccountingAPI,
  UpdateAccountingAPI,
  LoadMove,
  LoadOrders,
  LoadServiceMember,
  LoadBackupContacts,
} from './api.js';
import * as ReduxHelpers from 'shared/ReduxHelpers';

// Types
const loadAccountingType = 'LOAD_ACCOUNTING';
const updateAccountingType = 'UPDATE_ACCOUNTING';
const loadMoveType = 'LOAD_MOVE';
const loadOrdersType = 'LOAD_ORDERS';
const loadServiceMemberType = 'LOAD_SERVICE_MEMBER';
const loadBackupContactType = 'LOAD_BACKUP_CONTACT';

const LOAD_ACCOUNTING = ReduxHelpers.generateAsyncActionTypes(
  loadAccountingType,
);

const UPDATE_ACCOUNTING = ReduxHelpers.generateAsyncActionTypes(
  updateAccountingType,
);

const LOAD_MOVE = ReduxHelpers.generateAsyncActionTypes(loadMoveType);

const LOAD_ORDERS = ReduxHelpers.generateAsyncActionTypes(loadOrdersType);

const LOAD_SERVICE_MEMBER = ReduxHelpers.generateAsyncActionTypes(
  loadServiceMemberType,
);

const LOAD_BACKUP_CONTACT = ReduxHelpers.generateAsyncActionTypes(
  loadBackupContactType,
);

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

export const loadServiceMember = ReduxHelpers.generateAsyncActionCreator(
  loadServiceMemberType,
  LoadServiceMember,
);

export const loadBackupContacts = ReduxHelpers.generateAsyncActionCreator(
  loadBackupContactType,
  LoadBackupContacts,
);

export function loadMoveDependencies(moveId) {
  const actions = ReduxHelpers.generateAsyncActions(
    'loadMoveType | loadOrders | loadServiceMember | loadBackupContacts',
  );
  return async function(dispatch, getState) {
    dispatch(actions.start());
    try {
      await dispatch(loadMove(moveId));
      const move = getState().office.officeMove;
      await dispatch(loadOrders(move.orders_id));
      const orders = getState().office.officeOrders;
      await dispatch(loadServiceMember(orders.service_member_id));
      const sm = getState().office.officeServiceMember;
      await dispatch(loadBackupContacts(sm.id));
      return dispatch(actions.success());
    } catch (ex) {
      return dispatch(actions.error(ex));
    }
  };
}

// Reducer
const initialState = {
  accountingIsLoading: false,
  accountingIsUpdating: false,
  moveIsLoading: false,
  ordersAreLoading: false,
  serviceMemberIsLoading: false,
  backupContactsAreLoading: false,
  accountingHasLoadError: false,
  accountingHasLoadSuccess: null,
  accountingHasUpdateError: false,
  accountingHasUpdateSuccess: null,
  moveHasLoadError: false,
  moveHasLoadSuccess: null,
  ordersHaveLoadError: false,
  ordersHaveLoadSuccess: null,
  serviceMemberHasLoadError: false,
  serviceMemberHasLoadSuccess: null,
  backupContactsHaveLoadError: false,
  backupContactsHaveLoadSuccess: null,
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

    // SERVICE_MEMBER
    case LOAD_SERVICE_MEMBER.start:
      return Object.assign({}, state, {
        serviceMemberIsLoading: true,
        serviceMemberHasLoadSuccess: false,
      });
    case LOAD_SERVICE_MEMBER.success:
      return Object.assign({}, state, {
        serviceMemberIsLoading: false,
        officeServiceMember: action.payload,
        serviceMemberHasLoadSuccess: true,
        serviceMemberHasLoadError: false,
      });
    case LOAD_SERVICE_MEMBER.failure:
      return Object.assign({}, state, {
        serviceMemberIsLoading: false,
        officeServiceMember: null,
        serviceMemberHasLoadSuccess: false,
        serviceMemberHasLoadError: true,
        error: action.error.message,
      });

    // BACKUP CONTACT
    case LOAD_BACKUP_CONTACT.start:
      return Object.assign({}, state, {
        backupContactsAreLoading: true,
        backupContactsHaveLoadSuccess: false,
      });
    case LOAD_BACKUP_CONTACT.success:
      return Object.assign({}, state, {
        backupContactsAreLoading: false,
        officeBackupContacts: action.payload,
        backupContactsHaveLoadSuccess: true,
        backupContactsHaveLoadError: false,
      });
    case LOAD_BACKUP_CONTACT.failure:
      return Object.assign({}, state, {
        backupContactsAreLoading: false,
        officeBackupContacts: null,
        backupContactsHaveLoadSuccess: false,
        backupContactsHaveLoadError: true,
        error: action.error.message,
      });
    default:
      return state;
  }
}
