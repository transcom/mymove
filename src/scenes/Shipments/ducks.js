import { ShipmentsIndex, CreateShipment, UpdateShipment } from './api';

import * as ReduxHelpers from 'shared/ReduxHelpers';

// SINGLE RESOURCE ACTION TYPES
const createShipmentType = 'CREATE_SHIPMENT';
const updateShipmentType = 'UPDATE_SHIPMENT';

// SINGLE RESOURCE ACTION TYPES

const CREATE_SHIPMENT = ReduxHelpers.generateAsyncActionTypes(
  createShipmentType,
);
const UPDATE_SHIPMENT = ReduxHelpers.generateAsyncActionTypes(
  updateShipmentType,
);

// SINGLE-RESOURCE ACTION CREATORS

export const createShipment = ReduxHelpers.generateAsyncActionCreator(
  createShipmentType,
  CreateShipment,
);

export const updateShipment = ReduxHelpers.generateAsyncActionCreator(
  updateShipmentType,
  UpdateShipment,
);

// Types
export const SHOW_SHIPMENTS = 'SHOW_SHIPMENTS';
export const SHOW_SHIPMENTS_SUCCESS = 'SHOW_SHIPMENTS_SUCCESS';
export const SHOW_SHIPMENTS_FAILURE = 'SHOW_SHIPMENTS_FAILURE';

// Actions
export const createShowShipmentsRequest = () => ({
  type: SHOW_SHIPMENTS,
});

export const createShowShipmentsSuccess = shipments => ({
  type: SHOW_SHIPMENTS_SUCCESS,
  shipments,
});

export const createShowShipmentsFailure = error => ({
  type: SHOW_SHIPMENTS_FAILURE,
  error,
});

// Action Creator
export function loadShipments() {
  // Interpreted by the thunk middleware:
  return function(dispatch, getState) {
    dispatch(createShowShipmentsRequest());
    return ShipmentsIndex()
      .then(shipments => dispatch(createShowShipmentsSuccess(shipments)))
      .catch(error => dispatch(createShowShipmentsFailure(error)));
  };
}

// Reducer

const initialState = {
  shipments: [],
  hasError: false,
  shipmentIsCreating: false,
  shipmentIsUpdating: false,
  shipmentHasCreateSuccess: false,
  shipmentHasUpdateSuccess: false,
  shipmentHasCreateError: null,
  shipmentHasUpdateError: null,
};

export function shipmentsReducer(state = initialState, action) {
  switch (action.type) {
    case SHOW_SHIPMENTS_SUCCESS:
      return { shipments: action.shipments, hasError: false };
    case SHOW_SHIPMENTS_FAILURE:
      return { shipments: [], hasError: true };
    case CREATE_SHIPMENT.start:
      return Object.assign({}, state, {
        shipmentIsCreating: true,
        shipmentHasCreateSuccess: false,
      });
    case CREATE_SHIPMENT.success:
      return Object.assign({}, state, {
        shipmentIsCreating: false,
        shipment: action.payload,
        shipmentHasCreateSuccess: true,
        shipmentHasCreateError: false,
      });
    case CREATE_SHIPMENT.failure:
      return Object.assign({}, state, {
        shipmentIsCreating: false,
        shipment: null,
        shipmentHasCreateSuccess: false,
        shipmentHasCreateError: true,
        error: action.error,
      });
    case UPDATE_SHIPMENT.start:
      return Object.assign({}, state, {
        shipmentIsUpdating: true,
        shipmentHasUpdateSuccess: false,
      });
    case UPDATE_SHIPMENT.success:
      return Object.assign({}, state, {
        shipmentIsUpdating: false,
        shipment: action.payload,
        shipmentHasUpdateSuccess: true,
        shipmentHasUpdateError: false,
      });
    case UPDATE_SHIPMENT.failure:
      return Object.assign({}, state, {
        shipmentIsUpdating: false,
        shipment: null,
        shipmentHasUpdateSuccess: false,
        shipmentHasUpdateError: true,
        error: action.error,
      });
    default:
      return state;
  }
}
