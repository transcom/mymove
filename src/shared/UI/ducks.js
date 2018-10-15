import { get } from 'lodash';

import { fetchActive } from 'shared/utils';
import { GET_LOGGED_IN_USER } from 'shared/User/ducks';
import { selectShipment } from 'shared/Entities/modules/shipments';

const initialState = {
  currentShipmentID: null,
};

const SET_CURRENT_SHIPMENT = 'SET_CURRENT_SHIPMENT';

export default function uiReducer(state = initialState, action) {
  switch (action.type) {
    case GET_LOGGED_IN_USER.success:
      try {
        const activeOrders = fetchActive(get(action.payload, 'service_member.orders'));
        const activeMove = fetchActive(get(activeOrders, 'moves'));
        const activeShipment = fetchActive(get(activeMove, 'shipments'));
        return {
          ...state,
          currentShipmentID: activeShipment ? activeShipment.id : null,
        };
      } catch (e) {
        console.debug(e);
        return state;
      }
    case SET_CURRENT_SHIPMENT:
      return {
        ...state,
        currentShipmentID: action.shipment.id,
      };
    default:
      return state;
  }
}

export function setCurrentShipment(shipment) {
  return function(dispatch, getState) {
    return dispatch({ type: SET_CURRENT_SHIPMENT, shipment });
  };
}

// Selectors
export function currentShipment(state) {
  return selectShipment(state, get(state, 'ui.currentShipmentID'));
}
