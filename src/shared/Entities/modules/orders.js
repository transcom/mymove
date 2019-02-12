import { orders } from '../schema';
import { ADD_ENTITIES } from '../actions';
import { denormalize } from 'normalizr';
import { swaggerRequest } from 'shared/Swagger/request';
import { getClient } from 'shared/Swagger/api';
import { get } from 'lodash';

export const STATE_KEY = 'orders';
const loadOrdersLabel = 'Orders.loadOrders';
const updateOrdersLabel = 'Orders.updateOrders';

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

export function loadOrders(ordersId) {
  const label = loadOrdersLabel;
  const swaggerTag = 'orders.showOrders';
  return swaggerRequest(getClient, swaggerTag, { ordersId }, { label });
}

export function updateOrders(ordersId, orders) {
  const label = updateOrdersLabel;
  const swaggerTag = 'orders.updateOrders';
  return swaggerRequest(getClient, swaggerTag, { ordersId, updateOrders: orders }, { label });
}

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

export function selectUplodsForOrders(state, ordersId) {
  const orders = selectOrders(state, ordersId);
  const uploadedOrders = get(state, `entities.documents.${orders.uploaded_orders}`);
  if (uploadedOrders) {
    return uploadedOrders.uploads.map(uploadId => get(state, `entities.uploads.${uploadId}`));
  } else {
    return [];
  }
}
