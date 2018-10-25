import { get } from 'lodash';

import { fetchActive } from 'shared/utils';
import { GET_LOGGED_IN_USER } from 'shared/User/ducks';
import { selectShipment } from 'shared/Entities/modules/shipments';

const initialState = {
  currentShipmentID: null,
  notifications: [],
};

const SET_CURRENT_SHIPMENT_ID = 'SET_CURRENT_SHIPMENT_ID';
const ADD_NOTIFICATION = 'ADD_NOTIFICATION';

// Levels from USDS Web Design System
export const NOTIFICATION_SEVERITY = {
  success: 'success',
  warning: 'warning',
  error: 'error',
  info: 'info',
};

function createNotification(action) {
  return {
    title: action.title || 'Attention',
    message: action.message || '',
    severity: action.severity || 'info',
    createdAt: new Date(),
  };
}

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
        console.error(e);
        return state;
      }
    case SET_CURRENT_SHIPMENT_ID:
      return {
        ...state,
        currentShipmentID: action.shipmentID,
      };
    case ADD_NOTIFICATION:
      return {
        ...state,
        notifications: [...state.notifications, createNotification(action)],
      };
    default:
      return state;
  }
}

export function addNotification({ title, message, severity }) {
  return function(dispatch) {
    return dispatch({ type: ADD_NOTIFICATION, title, message, severity });
  };
}

export function setCurrentShipmentID(shipmentID) {
  return function(dispatch, getState) {
    return dispatch({ type: SET_CURRENT_SHIPMENT_ID, shipmentID });
  };
}

// Selectors

export function getNotifications(state) {
  return get(state, 'ui.notifications');
}

export function getCurrentShipmentID(state) {
  return get(state, 'ui.currentShipmentID');
}

export function getCurrentShipment(state) {
  return selectShipment(state, get(state, 'ui.currentShipmentID'));
}
