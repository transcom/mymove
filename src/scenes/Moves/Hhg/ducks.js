import { get } from 'lodash';
import * as ReduxHelpers from 'shared/ReduxHelpers';
import { GET_LOGGED_IN_USER } from 'shared/User/ducks';
import { fetchActive } from 'shared/utils';
import { GetShipment, GetMoveDatesSummary } from './api';

// Types
export const GET_SHIPMENT = ReduxHelpers.generateAsyncActionTypes(
  'GET_SHIPMENT',
);
export const GET_MOVE_DATES_SUMMARY = ReduxHelpers.generateAsyncActionTypes(
  'GET_MOVE_DATES_SUMMARY',
);

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

export function getMoveDatesSummary(moveId, moveDate) {
  const action = ReduxHelpers.generateAsyncActions('GET_MOVE_DATES_SUMMARY');
  return function(dispatch) {
    dispatch(action.start);
    return GetMoveDatesSummary(moveId, moveDate)
      .then(item => dispatch(action.success(item)))
      .catch(error => dispatch(action.error(error)));
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
  moveDates: null,
  isLoadingDates: false,
  hasDatesSuccess: false,
  hasDatesError: false,
};
export function hhgReducer(state = initialState, action) {
  switch (action.type) {
    case GET_LOGGED_IN_USER.success:
      // Initialize state when we get the logged in user
      const activeOrders = fetchActive(
        get(action.payload, 'service_member.orders'),
      );
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
    case GET_MOVE_DATES_SUMMARY.start:
      return Object.assign({}, state, {
        isLoadingDates: false,
        hasDatesSuccess: false,
        hasDatesError: false,
      });
    case GET_MOVE_DATES_SUMMARY.success:
      return Object.assign({}, state, {
        moveDates: action.payload,
        isLoadingDates: false,
        hasDatesSuccess: true,
        hasDatesError: false,
      });
    case GET_MOVE_DATES_SUMMARY.failure:
      return Object.assign({}, state, {
        isLoadingDates: false,
        hasDatesSuccess: false,
        hasDatesError: true,
        error: action.error,
      });
    default:
      return state;
  }
}
