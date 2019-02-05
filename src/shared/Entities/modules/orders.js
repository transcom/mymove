import { orders } from '../schema';
import { ADD_ENTITIES } from '../actions';
import { denormalize } from 'normalizr';
import { swaggerRequest } from 'shared/Swagger/request';
import { getClient } from 'shared/Swagger/api';

export const STATE_KEY = 'orders';
const loadOrdersLabel = 'Orders.loadOrders';

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

export const selectUpload = (state, id) => {
  return denormalize([id], orders, state.entities)[0];
};
