import { GetOrders, UpdateOrders } from './api.js';
import * as ReduxHelpers from 'shared/ReduxHelpers';

// Types
export const GET_ORDERS = ReduxHelpers.generateAsyncActionTypes('GET_ORDERS');

export const UPDATE_ORDERS = ReduxHelpers.generateAsyncActionTypes(
  'UPDATE_ORDERS',
);

// Action creation
export function updateOrders(orders) {
  const action = ReduxHelpers.generateAsyncActions('UPDATE_ORDERS');
  return function(dispatch, getState) {
    dispatch(action.start());
    const state = getState();
    const currentOrders = state.orders.currentOrders;
    if (currentOrders) {
      UpdateOrders(currentOrders.id, orders)
        .then(item => dispatch(action.success(item)))
        .catch(error => dispatch(action.error(error)));
    }
  };
}

export function loadOrders(serviceMemberId) {
  const action = ReduxHelpers.generateAsyncActions('GET_ORDERS');
  return function(dispatch, getState) {
    dispatch(action.start);
    const state = getState();
    const currentOrders = state.orders.currentOrders;
    if (!currentOrders) {
      console.log('hi');
      GetOrders(serviceMemberId)
        .then(item => dispatch(action.success(item)))
        .catch(error => dispatch(action.error(error)));
    }
  };
}

// Reducer
const initialState = {
  currentOrders: null,
  hasSubmitError: false,
  hasSubmitSuccess: false,
};
export function ordersReducer(state = initialState, action) {
  switch (action.type) {
    case GET_ORDERS.start:
      return Object.assign({}, state, {
        hasSubmitSuccess: false,
      });
    case GET_ORDERS.success:
      console.log(action.payload);
      return Object.assign({}, state, {
        currentOrders: action.payload,
        hasSubmitSuccess: true,
        hasSubmitError: false,
      });
    case GET_ORDERS.failure:
      return Object.assign({}, state, {
        currentOrders: null,
        hasSubmitSuccess: false,
        hasSubmitError: true,
        error: action.error,
      });
    default:
      return state;
  }
}
