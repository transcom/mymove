import { reject, pick, cloneDeep, concat, includes, get } from 'lodash';
import { CreateOrders, UpdateOrders, GetOrders, ShowServiceMemberOrders } from './api.js';
import { createOrUpdateMoveType } from 'scenes/Moves/ducks';
import { DeleteUploads } from 'shared/api';
import * as ReduxHelpers from 'shared/ReduxHelpers';
import { GET_LOGGED_IN_USER } from 'shared/Data/users';
import { fetchActive } from 'shared/utils';

// Types
const getOrdersType = 'GET_ORDERS';
export const GET_ORDERS = ReduxHelpers.generateAsyncActionTypes(getOrdersType);

const addUploadsType = 'ADD_UPLOADS';
export const ADD_UPLOADS = ReduxHelpers.generateAsyncActionTypes(addUploadsType);

const createOrUpdateOrdersType = 'CREATE_OR_UPDATE_ORDERS';
export const CREATE_OR_UPDATE_ORDERS = ReduxHelpers.generateAsyncActionTypes(createOrUpdateOrdersType);

const showCurrentOrdersType = 'SHOW_CURRENT_ORDERS';
export const SHOW_CURRENT_ORDERS = ReduxHelpers.generateAsyncActionTypes(showCurrentOrdersType);

const deleteUploadType = 'DELETE_UPLOAD';
export const DELETE_UPLOAD = ReduxHelpers.generateAsyncActionTypes(deleteUploadType);

// Actions
export const showServiceMemberOrders = ReduxHelpers.generateAsyncActionCreator(
  showCurrentOrdersType,
  ShowServiceMemberOrders,
);

export function createOrders(ordersPayload) {
  return function(dispatch, getState) {
    const action = ReduxHelpers.generateAsyncActions(createOrUpdateOrdersType);
    const moveAction = ReduxHelpers.generateAsyncActions(createOrUpdateMoveType);
    const state = getState();
    const currentOrders = state.orders.currentOrders;
    if (!currentOrders) {
      return CreateOrders(ordersPayload)
        .then(item => {
          const newMove = get(item, 'moves.0', null);
          dispatch(action.success(item));
          dispatch(moveAction.success(newMove));
        })
        .catch(error => dispatch(action.error(error)));
    } else {
      return Promise.reject();
    }
  };
}

export const updateOrders = ReduxHelpers.generateAsyncActionCreator(createOrUpdateOrdersType, UpdateOrders);

export const loadOrders = ReduxHelpers.generateAsyncActionCreator(getOrdersType, GetOrders);

// Deletes a single upload
export function deleteUpload(uploadId) {
  return function(dispatch, getState) {
    const action = ReduxHelpers.generateAsyncActions(deleteUploadType);
    const state = getState();
    if (state.orders.currentOrders) {
      return DeleteUploads(uploadId)
        .then(() => dispatch(action.success([uploadId])))
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
    if (state.orders.currentOrders && uploadIds.length) {
      return DeleteUploads(uploadIds)
        .then(() => dispatch(action.success(uploadIds)))
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
    if (state.orders.currentOrders) {
      dispatch(action.success(uploads));
    } else {
      dispatch(action.error(new Error("attempted to add uploads when orders don't exist")));
    }
  };
}

// Reducer
const initialState = {
  currentOrders: null,
  hasSubmitError: false,
  hasSubmitSuccess: false,
  hasLoadSuccess: false,
  hasLoadError: false,
  error: null,
};
function reshapeOrders(orders) {
  if (!orders) return null;
  return pick(orders, [
    'id',
    'has_dependents',
    'spouse_has_pro_gear',
    'issue_date',
    'new_duty_station',
    'orders_type',
    'report_by_date',
    'service_member_id',
    'uploaded_orders',
    'status',
  ]);
}
const removeUploads = (uploadIds, state) => {
  const newState = cloneDeep(state);
  newState.currentOrders.uploaded_orders.uploads = reject(state.currentOrders.uploaded_orders.uploads, upload => {
    return includes(uploadIds, upload.id);
  });
  return newState;
};
const insertUploads = (uploads, state) => {
  const newState = cloneDeep(state);
  newState.currentOrders.uploaded_orders.uploads = concat(state.currentOrders.uploaded_orders.uploads, ...uploads);
  return newState;
};
export function ordersReducer(state = initialState, action) {
  switch (action.type) {
    case GET_LOGGED_IN_USER.success:
      return Object.assign({}, state, {
        currentOrders: reshapeOrders(fetchActive(get(action.payload, 'service_member.orders'))),
        hasLoadError: false,
        hasLoadSuccess: true,
      });
    case CREATE_OR_UPDATE_ORDERS.success:
      return Object.assign({}, state, {
        currentOrders: reshapeOrders(action.payload),
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
        currentOrders: reshapeOrders(action.payload),
        hasLoadSuccess: true,
        hasLoadError: false,
        error: null,
      });
    case GET_ORDERS.failure:
      return Object.assign({}, state, {
        currentOrders: null,
        hasLoadSuccess: false,
        hasLoadError: true,
        error: action.error,
      });
    case SHOW_CURRENT_ORDERS.start:
      return Object.assign({}, state, {
        currentOrders: null,
        showCurrentOrdersSuccess: false,
      });
    case SHOW_CURRENT_ORDERS.success:
      return Object.assign({}, state, {
        currentOrders: reshapeOrders(action.payload),
        showCurrentOrdersSuccess: true,
        showCurrentOrdersError: false,
      });
    case SHOW_CURRENT_ORDERS.failure:
      return Object.assign({}, state, {
        currentOrders: null,
        showCurrentOrdersError: true,
        error: action.error,
      });
    case DELETE_UPLOAD.success:
      return {
        ...removeUploads(action.payload, state),
        hasSubmitSuccess: true,
        hasSubmitError: false,
        error: null,
      };
    case DELETE_UPLOAD.failure:
      return Object.assign({}, state, {
        hasSubmitSuccess: false,
        hasSubmitError: true,
        error: action.error,
      });
    case ADD_UPLOADS.success:
      return {
        ...insertUploads(action.payload, state),
        hasSubmitSuccess: true,
        hasSubmitError: false,
        error: null,
      };
    case ADD_UPLOADS.failure:
      return Object.assign({}, state, {
        hasSubmitSuccess: false,
        hasSubmitError: true,
        error: action.error,
      });
    default:
      return state;
  }
}
