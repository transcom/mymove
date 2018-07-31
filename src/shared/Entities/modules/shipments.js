import { shipments } from '../schema';
import { ADD_ENTITIES } from '../actions';
import { denormalize } from 'normalizr';

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

export function createShipment(moveId, shipment) {
  return async function(dispatch, getState, { schema }) {
    const client = await getClient();
    const response = await client.apis.shipment.createShipment({
      moveId,
      shipment,
    });
    checkResponse(response, 'failed to create shipment due to server error');
    const data = normalize(response.body, schema.shipment);
    dispatch(addEntities(data.entities));
    return response;
  };
}

export const updateShipment = (moveId, shipmentId, shipment) => {
  return async function(dispatch, getState, { schema }) {
    const client = await getClient();
    const response = await client.apis.shipment.updateShipment({
      moveId,
      shipmentId,
      shipment,
    });
    checkResponse(response, 'failed to update shipment due to server error');
    const data = normalize(response.body, schema.shipment);
    dispatch(addEntities(data.entities));
    return response;
  };
};

export const selectShipment = (state, id) => {
  return denormalize([id], shipments, state.entities)[0];
};
