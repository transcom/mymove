import { selectShipment } from 'shared/Entities/modules/shipments';
import { get } from 'lodash';

const initialState = {
  currentShipmentID: null,
};

const SET_CURRENT_SHIPMENT = 'SET_CURRENT_SHIPMENT';

export default function uiReducer(state = initialState, action) {
  switch (action.type) {
    case SET_CURRENT_SHIPMENT:
      return Object.assign({}, state, {
        currentShipmentID: action.shipment.id,
      });
    default:
      return state;
  }
}

export function setCurrentShipment(shipment) {
  return function(dispatch, getState) {
    dispatch({ type: SET_CURRENT_SHIPMENT, shipment });
  };
}

// Selectors
export function currentShipment(state) {
  return selectShipment(state, get(state, 'ui.currentShipmentID'));
}
