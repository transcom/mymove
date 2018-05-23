import { reject, cloneDeep, concat, includes, get, isNull } from 'lodash';
import {
  CreateOrders,
  UpdateOrders,
  GetOrders,
  ShowCurrentOrdersAPI,
} from './api.js';
import { DeleteUploads } from 'shared/api.js';
import { getEntitlements } from 'shared/entitlements.js';
import * as ReduxHelpers from 'shared/ReduxHelpers';

// Types
const getOrdersType = 'GET_ORDERS';
export const GET_ORDERS = ReduxHelpers.generateAsyncActionTypes(getOrdersType);

const addUploadsType = 'ADD_UPLOADS';
export const ADD_UPLOADS = ReduxHelpers.generateAsyncActionTypes(
  addUploadsType,
);

const createOrUpdateOrdersType = 'CREATE_OR_UPDATE_ORDERS';
export const CREATE_OR_UPDATE_ORDERS = ReduxHelpers.generateAsyncActionTypes(
  createOrUpdateOrdersType,
);

const showCurrentOrdersType = 'SHOW_CURRENT_ORDERS';
export const SHOW_CURRENT_ORDERS = ReduxHelpers.generateAsyncActionTypes(
  showCurrentOrdersType,
);

const deleteUploadType = 'DELETE_UPLOAD';
export const DELETE_UPLOAD = ReduxHelpers.generateAsyncActionTypes(
  deleteUploadType,
);

// Actions
export const showCurrentOrders = ReduxHelpers.generateAsyncActionCreator(
  showCurrentOrdersType,
  ShowCurrentOrdersAPI,
);

export const createOrders = ReduxHelpers.generateAsyncActionCreator(
  createOrUpdateOrdersType,
  CreateOrders,
);

export const updateOrders = ReduxHelpers.generateAsyncActionCreator(
  createOrUpdateOrdersType,
  UpdateOrders,
);

export const loadOrders = ReduxHelpers.generateAsyncActionCreator(
  getOrdersType,
  GetOrders,
);

// Deletes a single upload
export function deleteUpload(uploadId) {
  return function(dispatch, getState) {
    const action = ReduxHelpers.generateAsyncActions(deleteUploadType);
    const state = getState();
    const currentOrders = state.orders.currentOrders;
    if (currentOrders) {
      return DeleteUploads(uploadId)
        .then(() => {
          const uploads = currentOrders.uploaded_orders.uploads;
          currentOrders.uploaded_orders.uploads = reject(uploads, upload => {
            return uploadId === upload.id;
          });
          dispatch(action.success(currentOrders));
        })
        .catch(err => action.error(err));
    } else {
      return Promise.resolve();
    }
  };
}

// Deletes an array of uploads
export function deleteUploads(uploadIds) {
  return function(dispatch, getState) {
    const action = ReduxHelpers.generateAsyncActions(deleteUploadType);
    const state = getState();
    const currentOrders = cloneDeep(state.orders.currentOrders);
    if (currentOrders && uploadIds.length) {
      return DeleteUploads(uploadIds)
        .then(() => {
          const uploads = currentOrders.uploaded_orders.uploads;
          currentOrders.uploaded_orders.uploads = reject(uploads, upload => {
            return includes(uploadIds, upload.id);
          });
          dispatch(action.success(currentOrders));
        })
        .catch(err => action.error(err));
    } else {
      return Promise.resolve();
    }
  };
}

export function addUploads(uploads) {
  return function(dispatch, getState) {
    const action = ReduxHelpers.generateAsyncActions(addUploadsType);
    const state = getState();
    const currentOrders = state.orders.currentOrders;
    if (currentOrders) {
      currentOrders.uploaded_orders.uploads = concat(
        currentOrders.uploaded_orders.uploads,
        ...uploads,
      );
      dispatch(action.success(currentOrders));
    }
  };
}

// Selectors
export function loadEntitlements(state) {
  const hasDependents = get(
    state.loggedInUser,
    'loggedInUser.service_member.orders.0.has_dependents',
    null,
  );
  const rank = get(
    state.loggedInUser,
    'loggedInUser.service_member.rank',
    null,
  );
  if (isNull(hasDependents) || isNull(rank)) {
    return null;
  }
  return getEntitlements(rank, hasDependents);
}

// Reducer
const initialState = {
  currentOrders: null,
  hasSubmitError: false,
  hasSubmitSuccess: false,
  error: null,
};
export function ordersReducer(state = initialState, action) {
  switch (action.type) {
    case CREATE_OR_UPDATE_ORDERS.success:
      return Object.assign({}, state, {
        currentOrders: action.payload,
        pendingOrdersType: null,
        hasSubmitSuccess: true,
        hasSubmitError: false,
        error: null,
      });
    case CREATE_OR_UPDATE_ORDERS.failure:
      return Object.assign({}, state, {
        currentOrders: {},
        hasSubmitSuccess: false,
        hasSubmitError: true,
        error: action.error,
      });
    case GET_ORDERS.success:
      return Object.assign({}, state, {
        currentOrders: action.payload,
        hasSubmitSuccess: true,
        hasSubmitError: false,
        error: null,
      });
    case GET_ORDERS.failure:
      return Object.assign({}, state, {
        currentOrders: {},
        hasSubmitSuccess: false,
        hasSubmitError: true,
        error: action.error,
      });
    case SHOW_CURRENT_ORDERS.start:
      return Object.assign({}, state, {
        currentOrders: null,
        showCurrentOrdersSuccess: false,
      });
    case SHOW_CURRENT_ORDERS.success:
      return Object.assign({}, state, {
        currentOrders: action.payload,
        showCurrentOrdersSuccess: true,
        showCurrentOrdersError: false,
      });
    case SHOW_CURRENT_ORDERS.failure:
      const error = action.error.statusCode === 404 ? null : action.error;
      return Object.assign({}, state, {
        currentOrders: null,
        showCurrentOrdersError: true,
        error,
      });
    case DELETE_UPLOAD.success:
      return Object.assign({}, state, {
        currentOrders: action.payload,
        hasSubmitSuccess: true,
        hasSubmitError: false,
        error: null,
      });
    case DELETE_UPLOAD.failure:
      return Object.assign({}, state, {
        currentOrders: action.payload,
        hasSubmitSuccess: false,
        hasSubmitError: true,
        error: action.error,
      });
    case ADD_UPLOADS.success:
      return Object.assign({}, state, {
        currentOrders: action.payload,
        hasSubmitSuccess: true,
        hasSubmitError: false,
        error: null,
      });
    default:
      return state;
  }
}
