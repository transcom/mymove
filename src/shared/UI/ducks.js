import { get } from 'lodash';

import { fetchFirst, fetchActive, fetchActiveShipment } from 'shared/utils';
import { GET_LOGGED_IN_USER } from 'shared/Data/users';
import { selectShipment } from 'shared/Entities/modules/shipments';
import { selectMove } from 'shared/Entities/modules/moves';

const initialState = {
  currentShipmentID: null,
};

const SET_CURRENT_SHIPMENT_ID = 'SET_CURRENT_SHIPMENT_ID';

export default function uiReducer(state = initialState, action) {
  switch (action.type) {
    case GET_LOGGED_IN_USER.success:
      try {
        const orders = get(action.payload, 'service_member.orders');
        const activeOrders = fetchActive(orders);
        const activeMove = fetchActive(get(activeOrders, 'moves'));
        const activeShipment = fetchActiveShipment(get(activeMove, 'shipments'));
        const latestOrder = fetchFirst(orders);
        const latestMove = fetchFirst(get(latestOrder, 'moves'));

        return {
          ...state,
          currentShipmentID: activeShipment ? activeShipment.id : null,
          currentMoveID: activeMove ? activeMove.id : null,
          latestMoveID: latestMove ? latestMove.id : null,
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
  return function(dispatch, getState) {
    return dispatch({ type: SET_CURRENT_SHIPMENT_ID, shipmentID });
  };
}

// Selectors
export function getCurrentShipmentID(state) {
  return get(state, 'ui.currentShipmentID');
}

export function getCurrentShipment(state) {
  return selectShipment(state, get(state, 'ui.currentShipmentID'));
}

export function getCurrentMoveID(state) {
  return get(state, 'ui.currentMoveID');
}

export function getCurrentMove(state) {
  return selectMove(state, getCurrentMoveID(state));
}

export function getLatestMoveID(state) {
  return get(state, 'ui.latestMoveID');
}

export function getLatestMove(state) {
  return selectMove(state, getLatestMoveID(state));
}
