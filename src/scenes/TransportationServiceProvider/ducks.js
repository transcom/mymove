import { LoadShipment } from './api.js';

import * as ReduxHelpers from 'shared/ReduxHelpers';

// SINGLE RESOURCE ACTION TYPES
const loadShipmentType = 'LOAD_SHIPMENT';

// MULTIPLE-RESOURCE ACTION TYPES
const loadTspDependenciesType = 'LOAD_TSP_DEPENDENCIES';

// SINGLE RESOURCE ACTION TYPES

const LOAD_SHIPMENT = ReduxHelpers.generateAsyncActionTypes(loadShipmentType);

// MULTIPLE-RESOURCE ACTION TYPES

const LOAD_TSP_DEPENDENCIES = ReduxHelpers.generateAsyncActionTypes(
  loadTspDependenciesType,
);

// SINGLE-RESOURCE ACTION CREATORS

export const loadShipment = ReduxHelpers.generateAsyncActionCreator(
  loadShipmentType,
  LoadShipment,
);

// MULTIPLE-RESOURCE ACTION CREATORS
//
// These action types typically dispatch to other actions above to
// perform their work and exist to encapsulate when multiple requests
// need to be made in response to a user action.

export function loadShipmentDependencies(shipmentId) {
  const actions = ReduxHelpers.generateAsyncActions(loadTspDependenciesType);
  return async function(dispatch, getState) {
    dispatch(actions.start());
    try {
      await dispatch(loadShipment(shipmentId));
      return dispatch(actions.success());
    } catch (ex) {
      return dispatch(actions.error(ex));
    }
  };
}

// Selectors

// Reducer
const initialState = {};

export function tspReducer(state = initialState, action) {
  switch (action.type) {
    // SINGLE-RESOURCE ACTION TYPES

    // SHIPMENTS
    case LOAD_SHIPMENT.start:
      return Object.assign({}, state, {
        shipmentIsLoading: true,
        shipmentHasLoadSucces: false,
      });
    case LOAD_SHIPMENT.success:
      return Object.assign({}, state, {
        shipmentIsLoading: false,
        shipmentHasLoadSucces: true,
        shipmentHasLoadError: false,
        shipment: action.payload,
      });
    case LOAD_SHIPMENT.failure:
      return Object.assign({}, state, {
        shipmentIsLoading: false,
        shipmentHasLoadSucces: false,
        shipmentHasLoadError: null,
        shipment: null,
        error: action.error.message,
      });

    // MULTIPLE-RESOURCE ACTION TYPES
    //
    // These action types typically dispatch to other actions above to
    // perform their work and exist to encapsulate when multiple requests
    // need to be made in response to a user action.

    // ALL TSP DEPENDENCIES
    case LOAD_TSP_DEPENDENCIES.start:
      return Object.assign({}, state, {
        loadTspDependenciesHasSuccess: false,
        loadTspDependenciesHasError: false,
      });
    case LOAD_TSP_DEPENDENCIES.success:
      return Object.assign({}, state, {
        loadTspDependenciesHasSuccess: true,
        loadTspDependenciesHasError: false,
      });
    case LOAD_TSP_DEPENDENCIES.failure:
      return Object.assign({}, state, {
        loadTspDependenciesHasSuccess: false,
        loadTspDependenciesHasError: true,
        error: action.error.message,
      });
    default:
      return state;
  }
}
