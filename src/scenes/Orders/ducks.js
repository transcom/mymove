import { CreateOrders, UpdateOrders, GetOrders } from './api.js';

// Types
export const SET_PENDING_ORDERS_TYPE = 'SET_PENDING_ORDERS_TYPE';
export const CREATE_ORDERS = 'CREATE_ORDERS';
export const UPDATE_ORDERS = 'UPDATE_ORDERS';
export const CREATE_OR_UPDATE_ORDERS_SUCCESS =
  'CREATE_OR_UPDATE_ORDERS_SUCCESS';
export const CREATE_OR_UPDATE_ORDERS_FAILURE =
  'CREATE_OR_UPDATE_ORDERS_FAILURE';
export const GET_ORDERS = 'GET_ORDERS';
export const GET_ORDERS_SUCCESS = 'GET_ORDERS_SUCCESS';
export const GET_ORDERS_FAILURE = 'GET_ORDERS_FAILURE';

export const createOrdersRequest = () => ({
  type: CREATE_ORDERS,
});

export const updateOrdersRequest = () => ({
  type: UPDATE_ORDERS,
});

export const createOrUpdateOrdersSuccess = item => ({
  type: CREATE_OR_UPDATE_ORDERS_SUCCESS,
  item,
});

export const createOrUpdateOrdersFailure = error => ({
  type: CREATE_OR_UPDATE_ORDERS_FAILURE,
  error,
});

const getOrdersRequest = () => ({
  type: GET_ORDERS,
});

export const getOrdersSuccess = item => ({
  type: GET_ORDERS_SUCCESS,
  item,
  // item: items.length > 0 ? items[0] : null,
});

export const getOrdersFailure = error => ({
  type: GET_ORDERS_FAILURE,
  error,
});

// // Action creation
// export function setPendingOrdersType(value) {
//   return { type: SET_PENDING_ORDERS_TYPE, payload: value };
// }

export function createOrders(orderPayload) {
  return function(dispatch) {
    dispatch(createOrdersRequest());
    CreateOrders(orderPayload)
      .then(item => dispatch(createOrUpdateOrdersSuccess(item)))
      .catch(error => dispatch(createOrUpdateOrdersFailure(error)));
  };
}

export function updateOrders(orderId, orderPayload) {
  return function(dispatch) {
    dispatch(updateOrdersRequest());
    UpdateOrders(orderId, orderPayload)
      .then(item => dispatch(createOrUpdateOrdersSuccess(item)))
      .catch(error => dispatch(createOrUpdateOrdersFailure(error)));
  };
}

export function loadOrders(orderId) {
  return function(dispatch, getState) {
    const state = getState();
    const currentOrders = state.orders.currentOrders;
    if (!currentOrders) {
      dispatch(getOrdersRequest());
      GetOrders(orderId)
        .then(item => dispatch(getOrdersSuccess(item)))
        .catch(error => dispatch(getOrdersFailure(error)));
    }
  };
}

// Reducer
const initialState = {
  currentOrders: null,
  // pendingOrdersType: null,
  hasSubmitError: false,
  hasSubmitSuccess: false,
  error: null,
};
export function ordersReducer(state = initialState, action) {
  switch (action.type) {
    case UPDATE_ORDERS:
      return Object.assign({}, state, {
        hasSubmitSuccess: false,
      });
    case CREATE_OR_UPDATE_ORDERS_SUCCESS:
      return Object.assign({}, state, {
        currentOrders: action.item,
        pendingOrdersType: null,
        hasSubmitSuccess: true,
        hasSubmitError: false,
        error: null,
      });
    case CREATE_OR_UPDATE_ORDERS_FAILURE:
      return Object.assign({}, state, {
        currentOrders: {},
        hasSubmitSuccess: false,
        hasSubmitError: true,
        error: action.error,
      });
    case GET_ORDERS_SUCCESS:
      return Object.assign({}, state, {
        currentOrders: action.item,
        hasSubmitSuccess: true,
        hasSubmitError: false,
        error: null,
      });
    case GET_ORDERS_FAILURE:
      return Object.assign({}, state, {
        currentOrders: {},
        hasSubmitSuccess: false,
        hasSubmitError: true,
        error: action.error,
      });
    default:
      return state;
  }
}
