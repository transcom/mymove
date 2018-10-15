import { get } from 'lodash';
import * as ReduxHelpers from 'shared/ReduxHelpers';
import { GET_LOGGED_IN_USER } from 'shared/User/ducks';
import { fetchActive } from 'shared/utils';
import { GetShipment } from './api';

// Types
export const GET_SHIPMENT = ReduxHelpers.generateAsyncActionTypes('GET_SHIPMENT');

// Action creation
export function fetchShipment(moveId, shipmentId) {
  const action = ReduxHelpers.generateAsyncActions('GET_SHIPMENT');
  return function(dispatch, getState) {
    dispatch(action.start);
    const state = getState();
    const currentShipment = state.shipments.currentShipment;
    if (!currentShipment) {
      return GetShipment(moveId, shipmentId)
        .then(item => dispatch(action.success(item)))
        .catch(error => dispatch(action.error(error)));
    }
    return Promise.resolve();
  };
}
// Selectors

// Reducer
const initialState = {
  currentHhg: null,
  hasSubmitError: false,
  hasSubmitSuccess: false,
  hasLoadSuccess: false,
  hasLoadError: false,
};
export function hhgReducer(state = initialState, action) {
  switch (action.type) {
    case GET_LOGGED_IN_USER.success:
      // Initialize state when we get the logged in user
      const activeOrders = fetchActive(get(action.payload, 'service_member.orders'));
      const activeMove = fetchActive(get(activeOrders, 'moves'));
      const activeHhg = fetchActive(get(activeMove, 'shipments'));
      return Object.assign({}, state, {
        currentHhg: activeHhg,
        hasLoadSuccess: true,
        hasLoadError: false,
      });
    case GET_SHIPMENT.start:
      return Object.assign({}, state, {
        hasLoadSuccess: false,
      });
    case GET_SHIPMENT.success:
      console.log('payload', action.payload);
      return Object.assign({}, state, {
        currentShipment: action.payload,
        hasLoadSuccess: true,
        hasLoadError: false,
      });
    case GET_SHIPMENT.failure:
      return Object.assign({}, state, {
        currentShipment: null,
        hasLoadSuccess: false,
        hasLoadError: true,
        error: action.error,
      });
    default:
      return state;
  }
}
