import { isNull, get } from 'lodash';
import {
  LoadMove,
  LoadOrders,
  LoadServiceMember,
  UpdateServiceMember,
  LoadBackupContacts,
  UpdateBackupContact,
  LoadPPMs,
  ApproveBasics,
  ApprovePPM,
  ApproveReimbursement,
} from './api.js';

import { UpdateOrders } from 'scenes/Orders/api.js';
import { getEntitlements } from 'shared/entitlements.js';
import * as ReduxHelpers from 'shared/ReduxHelpers';

// SINGLE RESOURCE ACTION TYPES
const loadMoveType = 'LOAD_MOVE';
const loadOrdersType = 'LOAD_ORDERS';
const updateOrdersType = 'UPDATE_ORDERS';
const loadServiceMemberType = 'LOAD_SERVICE_MEMBER';
const updateServiceMemberType = 'UPDATE_SERVICE_MEMBER';
const loadBackupContactType = 'LOAD_BACKUP_CONTACT';
const updateBackupContactType = 'UPDATE_BACKUP_CONTACT';
const loadPPMsType = 'LOAD_PPMS';
const approveBasicsType = 'APPROVE_BASICS';
const approvePPMType = 'APPROVE_PPM';
const approveReimbursementType = 'APPROVE_REIMBURSEMENT';

// MULTIPLE-RESOURCE ACTION TYPES
const updateBackupInfoType = 'UPDATE_BACKUP_INFO';
const updateOrdersInfoType = 'UPDATE_ORDERS_INFO';
const loadDependenciesType = 'LOAD_DEPENDENCIES';

// SINGLE RESOURCE ACTION TYPES

const LOAD_MOVE = ReduxHelpers.generateAsyncActionTypes(loadMoveType);

const LOAD_ORDERS = ReduxHelpers.generateAsyncActionTypes(loadOrdersType);

const UPDATE_ORDERS = ReduxHelpers.generateAsyncActionTypes(updateOrdersType);

const LOAD_SERVICE_MEMBER = ReduxHelpers.generateAsyncActionTypes(
  loadServiceMemberType,
);

const UPDATE_SERVICE_MEMBER = ReduxHelpers.generateAsyncActionTypes(
  updateServiceMemberType,
);

const LOAD_BACKUP_CONTACT = ReduxHelpers.generateAsyncActionTypes(
  loadBackupContactType,
);

const UPDATE_BACKUP_CONTACT = ReduxHelpers.generateAsyncActionTypes(
  updateBackupContactType,
);

const LOAD_PPMS = ReduxHelpers.generateAsyncActionTypes(loadPPMsType);

const APPROVE_BASICS = ReduxHelpers.generateAsyncActionTypes(approveBasicsType);

const APPROVE_PPM = ReduxHelpers.generateAsyncActionTypes(approvePPMType);

export const APPROVE_REIMBURSEMENT = ReduxHelpers.generateAsyncActionTypes(
  approveReimbursementType,
);

// MULTIPLE-RESOURCE ACTION TYPES

const UPDATE_BACKUP_INFO = ReduxHelpers.generateAsyncActionTypes(
  updateBackupInfoType,
);

const UPDATE_ORDERS_INFO = ReduxHelpers.generateAsyncActionTypes(
  updateOrdersInfoType,
);

const LOAD_DEPENDENCIES = ReduxHelpers.generateAsyncActionTypes(
  loadDependenciesType,
);

// SINGLE-RESOURCE ACTION CREATORS

export const loadMove = ReduxHelpers.generateAsyncActionCreator(
  loadMoveType,
  LoadMove,
);

export const loadOrders = ReduxHelpers.generateAsyncActionCreator(
  loadOrdersType,
  LoadOrders,
);

export const updateOrders = ReduxHelpers.generateAsyncActionCreator(
  updateOrdersType,
  UpdateOrders,
);

export const loadServiceMember = ReduxHelpers.generateAsyncActionCreator(
  loadServiceMemberType,
  LoadServiceMember,
);

export const updateServiceMember = ReduxHelpers.generateAsyncActionCreator(
  updateServiceMemberType,
  UpdateServiceMember,
);

export const loadBackupContacts = ReduxHelpers.generateAsyncActionCreator(
  loadBackupContactType,
  LoadBackupContacts,
);

export const updateBackupContact = ReduxHelpers.generateAsyncActionCreator(
  updateBackupContactType,
  UpdateBackupContact,
);

export const loadPPMs = ReduxHelpers.generateAsyncActionCreator(
  loadPPMsType,
  LoadPPMs,
);

export const approveBasics = ReduxHelpers.generateAsyncActionCreator(
  approveBasicsType,
  ApproveBasics,
);

export const approvePPM = ReduxHelpers.generateAsyncActionCreator(
  approvePPMType,
  ApprovePPM,
);

export const approveReimbursement = ReduxHelpers.generateAsyncActionCreator(
  approveReimbursementType,
  ApproveReimbursement,
);
// MULTIPLE-RESOURCE ACTION CREATORS
//
// These action types typically dispatch to other actions above to
// perform their work and exist to encapsulate when multiple requests
// need to be made in response to a user action.

export function updateBackupInfo(
  serviceMemberId,
  serviceMemberPayload,
  backupContactId,
  backupContact,
) {
  const actions = ReduxHelpers.generateAsyncActions(updateBackupInfoType);
  return async function(dispatch, getState) {
    dispatch(actions.start());
    try {
      // TODO: perform these requests concurrently
      await dispatch(
        updateServiceMember(serviceMemberId, serviceMemberPayload),
      );
      await dispatch(updateBackupContact(backupContactId, backupContact));
      return dispatch(actions.success());
    } catch (ex) {
      return dispatch(actions.error(ex));
    }
  };
}

export function updateOrdersInfo(
  ordersId,
  orders,
  serviceMemberId,
  serviceMember,
) {
  const actions = ReduxHelpers.generateAsyncActions(updateOrdersInfoType);
  return async function(dispatch, getState) {
    dispatch(actions.start());
    try {
      // TODO: perform these requests concurrently
      serviceMember.current_station_id = serviceMember.current_station.id;
      await dispatch(updateServiceMember(serviceMemberId, serviceMember));

      orders.new_duty_station_id = orders.new_duty_station.id;
      await dispatch(updateOrders(ordersId, orders));
      return dispatch(actions.success());
    } catch (ex) {
      return dispatch(actions.error(ex));
    }
  };
}

export function loadMoveDependencies(moveId) {
  const actions = ReduxHelpers.generateAsyncActions(loadDependenciesType);
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
      // TODO: load PPMs in parallel to move using moveId
      await dispatch(loadPPMs(moveId));
      return dispatch(actions.success());
    } catch (ex) {
      return dispatch(actions.error(ex));
    }
  };
}

// Selectors
export function loadEntitlements(state) {
  const hasDependents = get(state, 'office.officeOrders.has_dependents', null);
  const rank = get(state, 'office.officeServiceMember.rank', null);
  if (isNull(hasDependents) || isNull(rank)) {
    return null;
  }
  return getEntitlements(rank, hasDependents);
}

// Reducer
const initialState = {
  moveIsLoading: false,
  ordersAreLoading: false,
  ordersAreUpdating: false,
  serviceMemberIsLoading: false,
  backupContactsAreLoading: false,
  ppmsAreLoading: false,
  moveHasLoadError: null,
  moveHasLoadSuccess: false,
  ordersHaveLoadError: null,
  ordersHaveLoadSuccess: false,
  ordersHaveUploadError: null,
  ordersHaveUploadSuccess: false,
  serviceMemberHasLoadError: null,
  serviceMemberHasLoadSuccess: false,
  serviceMemberHasUpdateError: null,
  serviceMemberHasUpdateSuccess: false,
  backupContactsHaveLoadError: null,
  backupContactsHaveLoadSuccess: false,
  ppmsHaveLoadError: null,
  ppmsHaveLoadSuccess: false,
  loadDependenciesHasError: null,
  loadDependenciesHasSuccess: false,
};

export function officeReducer(state = initialState, action) {
  switch (action.type) {
    // SINGLE-RESOURCE ACTION TYPES

    // MOVES
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
    case UPDATE_ORDERS.start:
      return Object.assign({}, state, {
        ordersAreUpdating: true,
        ordersHaveUpdateSuccess: false,
      });
    case UPDATE_ORDERS.success:
      return Object.assign({}, state, {
        ordersAreUpdating: false,
        officeOrders: action.payload,
        ordersHaveUpdateSuccess: true,
        ordersHaveUpdateError: false,
      });
    case UPDATE_ORDERS.failure:
      return Object.assign({}, state, {
        ordersAreUpdating: false,
        ordersHaveUpdateSuccess: false,
        ordersHaveUpdateError: true,
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
        serviceMemberHasLoadSuccess: false,
        serviceMemberHasLoadError: true,
        error: action.error.message,
      });

    case UPDATE_SERVICE_MEMBER.start:
      return Object.assign({}, state, {
        serviceMemberIsUpdating: true,
        serviceMemberHasUpdateSuccess: false,
      });
    case UPDATE_SERVICE_MEMBER.success:
      return Object.assign({}, state, {
        serviceMemberIsUpdating: false,
        officeServiceMember: action.payload,
        serviceMemberHasUpdateSuccess: true,
        serviceMemberHasUpdateError: false,
      });
    case UPDATE_SERVICE_MEMBER.failure:
      return Object.assign({}, state, {
        serviceMemberIsUpdating: false,
        serviceMemberHasUpdateSuccess: false,
        serviceMemberHasUpdateError: true,
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

    case UPDATE_BACKUP_CONTACT.start:
      return Object.assign({}, state, {
        backupContactIsUpdating: true,
        backupContactHasUpdateSuccess: false,
      });
    case UPDATE_BACKUP_CONTACT.success:
      return Object.assign({}, state, {
        backupContactIsUpdating: false,
        officeBackupContacts: [action.payload], // there is only one
        backupContactHasUpdateSuccess: true,
        backupContactHasUpdateFailure: false,
      });
    case UPDATE_BACKUP_CONTACT.failure:
      return Object.assign({}, state, {
        backupContactIsUpdating: false,
        backupContactHasUpdateSuccess: false,
        backupContactHasUpdateFailure: true,
        error: action.error.message,
      });

    // PPMs
    case LOAD_PPMS.start:
      return Object.assign({}, state, {
        PPMsAreLoading: true,
        PPMsHaveLoadSuccess: false,
      });
    case LOAD_PPMS.success:
      return Object.assign({}, state, {
        PPMsAreLoading: false,
        officePPMs: action.payload,
        PPMsHaveLoadSuccess: true,
        PPMsHaveLoadError: false,
      });
    case LOAD_PPMS.failure:
      return Object.assign({}, state, {
        PPMsAreLoading: false,
        officePPMs: null,
        PPMsHaveLoadSuccess: false,
        PPMsHaveLoadError: true,
        error: action.error.message,
      });

    // MOVE STATUS
    case APPROVE_BASICS.start:
      return Object.assign({}, state, {
        basicsIsApproving: true,
      });
    case APPROVE_BASICS.success:
      return Object.assign({}, state, {
        basicsIsApproving: false,
        officeMove: action.payload,
      });
    case APPROVE_BASICS.failure:
      return Object.assign({}, state, {
        basicsIsApproving: false,
        error: action.error.message,
      });

    // PPM STATUS
    case APPROVE_PPM.start:
      return Object.assign({}, state, {
        ppmIsApproving: true,
      });
    case APPROVE_PPM.success:
      return Object.assign({}, state, {
        ppmIsApproving: false,
        officePPMs: [action.payload],
      });
    case APPROVE_PPM.failure:
      return Object.assign({}, state, {
        ppmIsApproving: false,
        error: action.error.message,
      });

    // REIMBURSEMENT STATUS
    case APPROVE_REIMBURSEMENT.start:
      return Object.assign({}, state, {
        reimbursementIsApproving: true,
      });
    case APPROVE_REIMBURSEMENT.success:
      // TODO: Remove once we have multiple ppms
      let officePPM = get(state, 'officePPMs[0]');
      let newPPM = Object.assign({}, officePPM, {
        advance: action.payload,
      });
      return Object.assign({}, state, {
        reimbursementIsApproving: false,
        officePPMs: [newPPM],
      });
    case APPROVE_REIMBURSEMENT.failure:
      return Object.assign({}, state, {
        reimbursementIsApproving: false,
        error: action.error.message,
      });

    // MULTIPLE-RESOURCE ACTION TYPES
    //
    // These action types typically dispatch to other actions above to
    // perform their work and exist to encapsulate when multiple requests
    // need to be made in response to a user action.

    // BACKUP INFO
    case UPDATE_BACKUP_INFO.start:
      return Object.assign({}, state, {
        updateBackupInfoHasSuccess: false,
        updateBackupInfoHasError: false,
      });
    case UPDATE_BACKUP_INFO.success:
      return Object.assign({}, state, {
        updateBackupInfoHasSuccess: true,
        updateBackupInfoHasError: false,
      });
    case UPDATE_BACKUP_INFO.failure:
      return Object.assign({}, state, {
        updateBackupInfoHasSuccess: false,
        updateBackupInfoHasError: true,
        error: action.error.message,
      });

    // ORDERS INFO
    case UPDATE_ORDERS_INFO.start:
      return Object.assign({}, state, {
        updateOrdersInfoHasSuccess: false,
        updateOrdersInfoHasError: false,
      });
    case UPDATE_ORDERS_INFO.success:
      return Object.assign({}, state, {
        updateOrdersInfoHasSuccess: true,
        updateOrdersInfoHasError: false,
      });
    case UPDATE_ORDERS_INFO.failure:
      return Object.assign({}, state, {
        updateOrdersInfoHasSuccess: false,
        updateOrdersInfoHasError: true,
        error: action.error.message,
      });

    // ALL DEPENDENCIES
    case LOAD_DEPENDENCIES.start:
      return Object.assign({}, state, {
        loadDependenciesHasSuccess: false,
        loadDependenciesHasError: false,
      });
    case LOAD_DEPENDENCIES.success:
      return Object.assign({}, state, {
        loadDependenciesHasSuccess: true,
        loadDependenciesHasError: false,
      });
    case LOAD_DEPENDENCIES.failure:
      return Object.assign({}, state, {
        loadDependenciesHasSuccess: false,
        loadDependenciesHasError: true,
        error: action.error.message,
      });
    default:
      return state;
  }
}
