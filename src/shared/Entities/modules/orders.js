import { orders } from '../schema';
import { ADD_ENTITIES } from '../actions';
import { denormalize } from 'normalizr';
import { swaggerRequest } from 'shared/Swagger/request';
import { formatDateForSwagger } from 'shared/dates';
import { getClient } from 'shared/Swagger/api';
import { get, filter, isEmpty, isNull, head } from 'lodash';
import { fetchActive } from 'shared/utils';
import { selectServiceMemberFromLoggedInUser } from './serviceMembers';

export const STATE_KEY = 'orders';
export const loadOrdersLabel = 'Orders.loadOrders';
const updateOrdersLabel = 'Orders.updateOrders';
const createOrdersLabel = 'Orders.createOrders';
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

export function createOrders(orders, label = createOrdersLabel) {
  const swaggerTag = 'orders.createOrders';
  orders.report_by_date = formatDateForSwagger(orders.report_by_date);
  orders.issue_date = formatDateForSwagger(orders.issue_date);
  return swaggerRequest(getClient, swaggerTag, { createOrders: orders }, { label });
}

// Selectors
export const selectUpload = (state, id) => {
  return denormalize([id], orders, state.entities)[0];
};

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

export function selectUploadsForActiveOrders(state) {
  const orders = selectActiveOrLatestOrders(state);
  const uploadedOrders = get(state, `entities.documents.${orders.uploaded_orders}`);
  if (uploadedOrders) {
    return uploadedOrders.uploads
      .map((uploadId) => get(state, `entities.uploads.${uploadId}`))
      .filter((upload) => {
        if (upload === undefined) {
          console.warn('Upload not found in entities uploads');
          return false;
        }
        return true;
      });
  } else {
    return [];
  }
}

export function selectOrdersForServiceMemberId(state, serviceMemberId) {
  const orders = Object.values(state.entities.orders);
  filter(orders, (order) => order.service_member_id === serviceMemberId);
  return orders || [];
}

export function selectActiveOrLatestOrders(state) {
  // temp until full redux refactor: gets active or latest orders from entities if exist. If not, gets from orders.currentOrders.
  const serviceMember = get(state, 'user.userInfo.service_member', {});
  if (isNull(serviceMember)) {
    return {};
  }
  const orders = selectOrdersForServiceMemberId(state, serviceMember.id);
  let activeOrLatestOrders = fetchActive(orders) || head(orders);
  if (isEmpty(activeOrLatestOrders)) {
    const orders = get(state, 'user.userInfo.service_member.orders', {});
    activeOrLatestOrders = fetchActive(orders) || head(orders);
  }
  return activeOrLatestOrders || {};
}

// use this for redux refactored parts where we've loaded orders into entities
export function selectActiveOrLatestOrdersFromEntities(state) {
  const serviceMember = selectServiceMemberFromLoggedInUser(state);
  if (isNull(serviceMember)) {
    return {};
  }
  const orders = selectOrdersForServiceMemberId(state, serviceMember.id);
  let activeOrLatestOrders = fetchActive(orders) || head(orders);
  return activeOrLatestOrders || {};
}
