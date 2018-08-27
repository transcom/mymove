import { LoadShipment, PatchShipment, AcceptShipment } from './api.js';

import * as ReduxHelpers from 'shared/ReduxHelpers';

// SINGLE RESOURCE ACTION TYPES
const loadShipmentType = 'LOAD_SHIPMENT';
const patchShipmentType = 'PATCH_SHIPMENT';
const acceptShipmentType = 'ACCEPT_SHIPMENT';
const REMOVE_BANNER = 'REMOVE_BANNER';

// MULTIPLE-RESOURCE ACTION TYPES
const loadTspDependenciesType = 'LOAD_TSP_DEPENDENCIES';

// SINGLE RESOURCE ACTION TYPES
const LOAD_SHIPMENT = ReduxHelpers.generateAsyncActionTypes(loadShipmentType);
const PATCH_SHIPMENT = ReduxHelpers.generateAsyncActionTypes(patchShipmentType);
const ACCEPT_SHIPMENT = ReduxHelpers.generateAsyncActionTypes(
  acceptShipmentType,
);

// MULTIPLE-RESOURCE ACTION TYPES

const LOAD_TSP_DEPENDENCIES = ReduxHelpers.generateAsyncActionTypes(
  loadTspDependenciesType,
);

// SINGLE-RESOURCE ACTION CREATORS

export const loadShipment = ReduxHelpers.generateAsyncActionCreator(
  loadShipmentType,
  LoadShipment,
);

export const patchShipment = ReduxHelpers.generateAsyncActionCreator(
  patchShipmentType,
  PatchShipment,
);

export const acceptShipment = shipmentId => {
  const actions = ReduxHelpers.generateAsyncActions(acceptShipmentType);
  return async function(dispatch, getState) {
    dispatch(actions.start());
    return AcceptShipment(shipmentId)
      .then(
        item => dispatch(actions.success(item)),
        error => dispatch(actions.error(error)),
      )
      .then(() => {
        setTimeout(() => dispatch(removeBanner()), 10000);
      });
  };
};

export const removeBanner = () => {
  return {
    type: REMOVE_BANNER,
  };
};

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
const initialState = {
  shipmentIsAccepting: false,
  shipmentHasAcceptError: null,
  shipmentHasAcceptSuccess: false,
  flashMessage: false,
};

export function tspReducer(state = initialState, action) {
  switch (action.type) {
    // SINGLE-RESOURCE ACTION TYPES

    // SHIPMENTS
    case LOAD_SHIPMENT.start:
      return Object.assign({}, state, {
        shipmentIsLoading: true,
        shipmentHasLoadSuccess: false,
      });
    case LOAD_SHIPMENT.success:
      return Object.assign({}, state, {
        shipmentIsLoading: false,
        shipmentHasLoadSuccess: true,
        shipmentHasLoadError: false,
        shipment: action.payload,
      });
    case LOAD_SHIPMENT.failure:
      return Object.assign({}, state, {
        shipmentIsLoading: false,
        shipmentHasLoadSuccess: false,
        shipmentHasLoadError: null,
        shipment: null,
        error: action.error.message,
      });
    case PATCH_SHIPMENT.start:
      return Object.assign({}, state, {
        shipmentPatchSuccess: false,
      });
    case PATCH_SHIPMENT.success:
      return Object.assign({}, state, {
        shipmentPatchSuccess: true,
        shipmentPatchError: false,
        shipment: action.payload,
      });
    case PATCH_SHIPMENT.failure:
      return Object.assign({}, state, {
        shipmentPatchSuccess: false,
        shipmentPatchError: null,
        error: action.error.message,
      });
    case ACCEPT_SHIPMENT.start:
      return Object.assign({}, state, {
        shipmentIsAccepting: true,
        shipmentHasAcceptSuccess: false,
      });
    case ACCEPT_SHIPMENT.success:
      return Object.assign({}, state, {
        shipmentIsAccepting: false,
        shipmentHasAcceptSuccess: true,
        shipmentHasAcceptError: false,
        shipment: action.payload,
        flashMessage: true,
      });
    case ACCEPT_SHIPMENT.failure:
      return Object.assign({}, state, {
        shipmentIsAccepting: false,
        shipmentHasAcceptSuccess: false,
        shipmentHasAcceptError: null,
        error: action.error.message,
        flashMessage: false,
      });
    case REMOVE_BANNER:
      return Object.assign({}, state, {
        flashMessage: false,
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
