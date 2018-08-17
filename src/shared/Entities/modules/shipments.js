import { denormalize, normalize } from 'normalizr';

import { shipments } from '../schema';
import { ADD_ENTITIES, addEntities } from '../actions';
import { getClient, checkResponse } from 'shared/api';

export const STATE_KEY = 'shipments';

export default function reducer(state = {}, action) {
  switch (action.type) {
    case ADD_ENTITIES:
      return {
        ...state,
        ...action.payload.shipments,
      };

    default:
      return state;
  }
}

export function createOrUpdateShipment(moveId, shipment, id) {
  if (id) {
    return updateShipment(id, shipment);
  } else {
    return createShipment(moveId, shipment);
  }
}

export function getShipment(moveId, shipmentId) {
  return async function(dispatch, getState, { schema }) {
    const client = await getClient();
    const response = await client.apis.shipments.getShipment({
      moveId,
      shipmentId,
    });
    checkResponse(response, 'failed to get shipment due to server error');
    const data = normalize(response.body, schema.shipment);
    dispatch(addEntities(data.entities));
    return response;
  };
}

export function createShipment(
  moveId,
  shipment /*shape: {pickup_address, requested_pickup_date, weight_estimate}*/,
) {
  return async function(dispatch, getState, { schema }) {
    const client = await getClient();
    const response = await client.apis.shipments.createShipment({
      moveId,
      shipment,
    });
    checkResponse(response, 'failed to create shipment due to server error');
    const data = normalize(response.body, schema.shipment);
    dispatch(addEntities(data.entities));
    return response;
  };
}

export function updateShipment(
  moveId,
  shipmentId,
  shipment /*shape: {pickup_address, requested_pickup_date, weight_estimate}*/,
) {
  return async function(dispatch, getState, { schema }) {
    const client = await getClient();
    const response = await client.apis.shipments.patchShipment({
      moveId,
      shipmentId,
      shipment,
    });
    checkResponse(response, 'failed to update shipment due to server error');
    const data = normalize(response.body, schema.shipment);
    dispatch(addEntities(data.entities));
    return response;
  };
}

export function selectShipment(state, id) {
  if (!id) {
    return null;
  }
  return denormalize([id], shipments, state.entities)[0];
}
