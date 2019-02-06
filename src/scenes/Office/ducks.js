import { isNull, get } from 'lodash';
import { LoadMove, LoadOrders, DownloadPPMAttachments, PatchShipment, SendHHGInvoice } from './api.js';

import { UpdateOrders } from 'scenes/Orders/api.js';
import { getEntitlements } from 'shared/entitlements.js';
import { loadPPMs } from 'shared/Entities/modules/ppms';

import {
  loadServiceMember,
  updateServiceMember,
  loadBackupContacts,
  updateBackupContact,
} from 'shared/Entities/modules/serviceMembers';
import * as ReduxHelpers from 'shared/ReduxHelpers';

// SINGLE RESOURCE ACTION TYPES
const loadMoveType = 'LOAD_MOVE';
const loadOrdersType = 'LOAD_ORDERS';
const updateOrdersType = 'UPDATE_ORDERS';
const patchShipmentType = 'PATCH_SHIPMENT';
const sendHHGInvoiceType = 'SEND_HHG_INVOICE';
const downloadPPMAttachmentsType = 'DOWNLOAD_ATTACHMENTS';
const REMOVE_BANNER = 'REMOVE_BANNER';
const SHOW_BANNER = 'SHOW_BANNER';
const RESET_MOVE = 'RESET_MOVE';
const DRAFT_HHG_INVOICE = 'DRAFT_INVOICE';
const RESET_HHG_INVOICE = 'RESET_INVOICE';

// MULTIPLE-RESOURCE ACTION TYPES
const updateBackupInfoType = 'UPDATE_BACKUP_INFO';
const updateOrdersInfoType = 'UPDATE_ORDERS_INFO';
const loadDependenciesType = 'LOAD_DEPENDENCIES';

// SINGLE RESOURCE ACTION TYPES

export const resetMove = () => ({
  type: RESET_MOVE,
});

export const draftInvoice = () => ({
  type: DRAFT_HHG_INVOICE,
});

export const resetInvoiceFlow = () => ({
  type: RESET_HHG_INVOICE,
});

const LOAD_MOVE = ReduxHelpers.generateAsyncActionTypes(loadMoveType);

const LOAD_ORDERS = ReduxHelpers.generateAsyncActionTypes(loadOrdersType);

const UPDATE_ORDERS = ReduxHelpers.generateAsyncActionTypes(updateOrdersType);

const PATCH_SHIPMENT = ReduxHelpers.generateAsyncActionTypes(patchShipmentType);

const SEND_HHG_INVOICE = ReduxHelpers.generateAsyncActionTypes(sendHHGInvoiceType);

export const DOWNLOAD_ATTACHMENTS = ReduxHelpers.generateAsyncActionTypes(downloadPPMAttachmentsType);

// MULTIPLE-RESOURCE ACTION TYPES

const UPDATE_BACKUP_INFO = ReduxHelpers.generateAsyncActionTypes(updateBackupInfoType);

const UPDATE_ORDERS_INFO = ReduxHelpers.generateAsyncActionTypes(updateOrdersInfoType);

const LOAD_DEPENDENCIES = ReduxHelpers.generateAsyncActionTypes(loadDependenciesType);

// SINGLE-RESOURCE ACTION CREATORS

export const loadMove = ReduxHelpers.generateAsyncActionCreator(loadMoveType, LoadMove);

export const loadOrders = ReduxHelpers.generateAsyncActionCreator(loadOrdersType, LoadOrders);

export const updateOrders = ReduxHelpers.generateAsyncActionCreator(updateOrdersType, UpdateOrders);

export const patchShipment = ReduxHelpers.generateAsyncActionCreator(patchShipmentType, PatchShipment);

export const sendHHGInvoice = ReduxHelpers.generateAsyncActionCreator(sendHHGInvoiceType, SendHHGInvoice);

export const downloadPPMAttachments = ReduxHelpers.generateAsyncActionCreator(
  downloadPPMAttachmentsType,
  DownloadPPMAttachments,
);

export const removeBanner = () => {
  return {
    type: REMOVE_BANNER,
  };
};

export const showBanner = () => {
  return {
    type: SHOW_BANNER,
  };
};
// MULTIPLE-RESOURCE ACTION CREATORS
//
// These action types typically dispatch to other actions above to
// perform their work and exist to encapsulate when multiple requests
// need to be made in response to a user action.

export function updateBackupInfo(serviceMemberId, serviceMemberPayload, backupContactId, backupContact) {
  const actions = ReduxHelpers.generateAsyncActions(updateBackupInfoType);
  return async function(dispatch, getState) {
    dispatch(actions.start());
    try {
      // TODO: perform these requests concurrently
      await dispatch(updateServiceMember(serviceMemberId, serviceMemberPayload));
      await dispatch(updateBackupContact(backupContactId, backupContact));
      return dispatch(actions.success());
    } catch (ex) {
      return dispatch(actions.error(ex));
    }
  };
}

export function updateOrdersInfo(ordersId, orders, serviceMemberId, serviceMember) {
  const actions = ReduxHelpers.generateAsyncActions(updateOrdersInfoType);
  return async function(dispatch, getState) {
    dispatch(actions.start());
    try {
      // TODO: perform these requests concurrently
      serviceMember.current_station_id = serviceMember.current_station.id;
      await dispatch(updateServiceMember(serviceMemberId, serviceMember));

      if (!orders.has_dependents) {
        orders.spouse_has_pro_gear = false;
      }

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
      const serviceMemberId = orders.service_member_id;
      await dispatch(loadServiceMember(serviceMemberId));
      await dispatch(loadBackupContacts(serviceMemberId));
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
  const spouseHasProGear = get(state, 'office.officeOrders.spouse_has_pro_gear', null);
  const rank = get(state, 'office.officeServiceMember.rank', null);
  if (isNull(hasDependents) || isNull(spouseHasProGear) || isNull(rank)) {
    return null;
  }
  return getEntitlements(rank, hasDependents, spouseHasProGear);
}

// Reducer
const initialState = {
  moveIsLoading: false,
  ordersAreLoading: false,
  ordersAreUpdating: false,
  ppmsAreLoading: false,
  ppmIsUpdating: false,
  moveHasLoadError: null,
  moveHasLoadSuccess: false,
  officeMove: {},
  ordersHaveLoadError: null,
  ordersHaveLoadSuccess: false,
  ordersHaveUploadError: null,
  ordersHaveUploadSuccess: false,
  downloadAttachmentsHasError: null,
  ppmsHaveLoadError: null,
  ppmsHaveLoadSuccess: false,
  ppmHasUpdateError: null,
  ppmHasUpdateSuccess: false,
  hhgInvoiceIsSending: false,
  hhgInvoiceHasSendSuccess: false,
  hhgInvoiceHasFailure: false,
  hhgInvoiceInDraft: false,
  loadDependenciesHasError: null,
  loadDependenciesHasSuccess: false,
  flashMessage: false,
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
        officeShipment: get(action.payload, 'shipments.0', null),
        moveHasLoadSuccess: true,
        moveHasLoadError: false,
      });
    case LOAD_MOVE.failure:
      return Object.assign({}, state, {
        moveIsLoading: false,
        officeMove: null,
        officeShipment: null,
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

    // SHIPMENT
    case PATCH_SHIPMENT.start:
      return Object.assign({}, state, {
        shipmentIsUpdating: true,
        shipmentPatchSuccess: false,
      });
    case PATCH_SHIPMENT.success:
      return Object.assign({}, state, {
        shipmentIsUpdating: false,
        officeShipment: action.payload,
        shipmentPatchSuccess: true,
        shipmentPatchError: false,
      });
    case PATCH_SHIPMENT.failure:
      return Object.assign({}, state, {
        shipmentIsUpdating: false,
        shipmentPatchSuccess: false,
        shipmentPatchError: true,
        error: action.error.message,
      });
    case SEND_HHG_INVOICE.start:
      return Object.assign({}, state, {
        hhgInvoiceIsSending: true,
        hhgInvoiceHasSendSuccess: false,
        hhgInvoiceHasFailure: false,
        hhgInvoiceInDraft: false,
      });
    case SEND_HHG_INVOICE.success:
      return Object.assign({}, state, {
        hhgInvoiceIsSending: false,
        hhgInvoiceHasSendSuccess: true,
        hhgInvoiceHasFailure: false,
        hhgInvoiceInDraft: false,
      });
    case SEND_HHG_INVOICE.failure:
      return Object.assign({}, state, {
        hhgInvoiceIsSending: false,
        hhgInvoiceHasFailure: true,
        hhgInvoiceInDraft: false,
        error: action.error.message,
      });
    case DRAFT_HHG_INVOICE:
      return Object.assign({}, state, {
        hhgInvoiceInDraft: true,
        hhgInvoiceIsSending: false,
        hhgInvoiceHasSendSuccess: false,
        hhgInvoiceHasFailure: false,
      });
    case RESET_HHG_INVOICE:
      return Object.assign({}, state, {
        hhgInvoiceInDraft: false,
        hhgInvoiceIsSending: false,
        hhgInvoiceHasSendSuccess: false,
        hhgInvoiceHasFailure: false,
      });

    case SHOW_BANNER:
      return Object.assign({}, state, {
        flashMessage: true,
      });
    case REMOVE_BANNER:
      return Object.assign({}, state, {
        flashMessage: false,
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

    // PPM ATTACHMENTS GENERATOR
    case DOWNLOAD_ATTACHMENTS.start:
      return Object.assign({}, state, {
        downloadAttachmentsHasError: null,
      });
    case DOWNLOAD_ATTACHMENTS.failure:
      return Object.assign({}, state, {
        downloadAttachmentsHasError: action.error,
      });
    case RESET_MOVE:
      return Object.assign({}, state, {
        officeShipment: {},
      });
    default:
      return state;
  }
}
