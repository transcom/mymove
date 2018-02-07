import { AvailableShipmentsIndex, AwardedShipmentsIndex } from 'shared/api';

// AVAILABLE SHIPMENTS

// Types
export const SHOW_AVAILABLE_SHIPMENTS = 'SHOW_AVAILABLE_SHIPMENTS';
export const SHOW_AVAILABLE_SHIPMENTS_SUCCESS =
  'SHOW_AVAILABLE_SHIPMENTS_SUCCESS';
export const SHOW_AVAILABLE_SHIPMENTS_FAILURE =
  'SHOW_AVAILABLE_SHIPMENTS_FAILURE';

// Actions
export const createShowAvailableShipmentsRequest = () => ({
  type: SHOW_AVAILABLE_SHIPMENTS,
});

export const createShowAvailableShipmentsSuccess = shipments => ({
  type: SHOW_AVAILABLE_SHIPMENTS_SUCCESS,
  shipments,
});

export const createShowAvailableShipmentsFailure = error => ({
  type: SHOW_AVAILABLE_SHIPMENTS_FAILURE,
  error,
});

// Action Creator
export function loadAvailableShipments() {
  // Interpreted by the thunk middleware:
  return function(dispatch, getState) {
    dispatch(createShowAvailableShipmentsRequest());
    AvailableShipmentsIndex()
      .then(shipments =>
        dispatch(createShowAvailableShipmentsSuccess(shipments)),
      )
      .catch(error => dispatch(createShowAvailableShipmentsFailure(error)));
  };
}

// Reducer
export function availableShipmentsReducer(
  state = { shipments: null, hasError: false },
  action,
) {
  switch (action.type) {
    case SHOW_AVAILABLE_SHIPMENTS_SUCCESS:
      return { shipments: action.shipments, hasError: false };
    case SHOW_AVAILABLE_SHIPMENTS_FAILURE:
      return { shipments: null, hasError: true };
    default:
      return state;
  }
}

// AWARDED SHIPMENTS

// Types
export const SHOW_AWARDED_SHIPMENTS = 'SHOW_AWARDED_SHIPMENTS';
export const SHOW_AWARDED_SHIPMENTS_SUCCESS = 'SHOW_AWARDED_SHIPMENTS_SUCCESS';
export const SHOW_AWARDED_SHIPMENTS_FAILURE = 'SHOW_AWARDED_SHIPMENTS_FAILURE';

// Actions
export const createShowAwardedShipmentsRequest = () => ({
  type: SHOW_AWARDED_SHIPMENTS,
});

export const createShowAwardedShipmentsSuccess = shipments => ({
  type: SHOW_AWARDED_SHIPMENTS_SUCCESS,
  shipments,
});

export const createShowAwardedShipmentsFailure = error => ({
  type: SHOW_AWARDED_SHIPMENTS_FAILURE,
  error,
});

// Action Creator
export function loadAwardedShipments() {
  // Interpreted by the thunk middleware:
  return function(dispatch, getState) {
    dispatch(createShowAwardedShipmentsRequest());
    AwardedShipmentsIndex()
      .then(shipments => dispatch(createShowAwardedShipmentsSuccess(shipments)))
      .catch(error => dispatch(createShowAwardedShipmentsFailure(error)));
  };
}

// Reducer
export function awardedShipmentsReducer(
  state = { shipments: null, hasError: false },
  action,
) {
  switch (action.type) {
    case SHOW_AWARDED_SHIPMENTS_SUCCESS:
      return { shipments: action.shipments, hasError: false };
    case SHOW_AWARDED_SHIPMENTS_FAILURE:
      return { shipments: null, hasError: true };
    default:
      return state;
  }
}
