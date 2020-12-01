import { pick, get } from 'lodash';
import { GetOrders } from './api.js';
import * as ReduxHelpers from 'shared/ReduxHelpers';
import { GET_LOGGED_IN_USER } from 'shared/Data/users';
import { fetchActive } from 'shared/utils';

// Types
const getOrdersType = 'GET_ORDERS';
export const GET_ORDERS = ReduxHelpers.generateAsyncActionTypes(getOrdersType);

const showCurrentOrdersType = 'SHOW_CURRENT_ORDERS';
export const SHOW_CURRENT_ORDERS = ReduxHelpers.generateAsyncActionTypes(showCurrentOrdersType);

export const loadOrders = ReduxHelpers.generateAsyncActionCreator(getOrdersType, GetOrders);

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
export function ordersReducer(state = initialState, action) {
  switch (action.type) {
    case GET_LOGGED_IN_USER.success:
      return Object.assign({}, state, {
        currentOrders: reshapeOrders(fetchActive(get(action.payload, 'service_member.orders'))),
        hasLoadError: false,
        hasLoadSuccess: true,
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
    default:
      return state;
  }
}
