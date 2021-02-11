import { get } from 'lodash';

import { fetchActive, fetchActiveShipment } from 'shared/utils';
import { GET_LOGGED_IN_USER_SUCCESS } from 'store/auth/actions';

const initialState = {
  currentShipmentID: null,
};

const SET_CURRENT_SHIPMENT_ID = 'SET_CURRENT_SHIPMENT_ID';

export default function uiReducer(state = initialState, action) {
  switch (action.type) {
    case GET_LOGGED_IN_USER_SUCCESS:
      try {
        const activeOrders = fetchActive(get(action.payload, 'service_member.orders'));
        const activeMove = fetchActive(get(activeOrders, 'moves'));
        const activeShipment = fetchActiveShipment(get(activeMove, 'shipments'));

        return {
          ...state,
          currentShipmentID: activeShipment ? activeShipment.id : null,
        };
      } catch (e) {
        console.error(e);
        return state;
      }
    case SET_CURRENT_SHIPMENT_ID:
      return {
        ...state,
        currentShipmentID: action.shipmentID,
      };
    default:
      return state;
  }
}

export function setCurrentShipmentID(shipmentID) {
  return function (dispatch, getState) {
    return dispatch({ type: SET_CURRENT_SHIPMENT_ID, shipmentID });
  };
}

// Selectors
export function getCurrentShipmentID(state) {
  return get(state, 'ui.currentShipmentID');
}
