import { get } from 'lodash';

import { ADD_ENTITIES } from '../actions';

import { swaggerRequest } from 'shared/Swagger/request';
import { formatDateForSwagger } from 'shared/dates';
import { getClient } from 'shared/Swagger/api';

export const STATE_KEY = 'orders';
export const loadOrdersLabel = 'Orders.loadOrders';
const updateOrdersLabel = 'Orders.updateOrders';
export const getLatestOrdersLabel = 'Orders.showServiceMemberOrders';

export default function reducer(state = {}, action) {
  switch (action.type) {
    case ADD_ENTITIES:
      return {
        ...state,
        ...action.payload.orders,
      };

    default:
      return state;
  }
}

export function fetchLatestOrders(serviceMemberId, label = getLatestOrdersLabel) {
  const swaggerTag = 'service_members.showServiceMemberOrders';
  return swaggerRequest(getClient, swaggerTag, { serviceMemberId }, { label });
}

export function loadOrders(ordersId, label = loadOrdersLabel) {
  const swaggerTag = 'orders.showOrders';
  return swaggerRequest(getClient, swaggerTag, { ordersId }, { label });
}

export function updateOrders(ordersId, orders, label = updateOrdersLabel) {
  const swaggerTag = 'orders.updateOrders';
  orders.report_by_date = formatDateForSwagger(orders.report_by_date);
  orders.issue_date = formatDateForSwagger(orders.issue_date);
  return swaggerRequest(getClient, swaggerTag, { ordersId, updateOrders: orders }, { label });
}

// Selectors
export function selectOrders(state, ordersId) {
  return get(state, `entities.orders.${ordersId}`) || {};
}

export function selectOrdersForMove(state, moveId) {
  const ordersId = get(state, `entities.moves.${moveId}.orders_id`);
  if (ordersId) {
    return selectOrders(state, ordersId);
  } else {
    return {};
  }
}
