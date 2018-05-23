import { get, find } from 'lodash';
import { CreatePpm, UpdatePpm, GetPpm, GetPpmWeightEstimate } from './api.js';
import * as ReduxHelpers from 'shared/ReduxHelpers';

// Types
export const SET_PENDING_PPM_SIZE = 'SET_PENDING_PPM_SIZE';
export const SET_PENDING_PPM_WEIGHT = 'SET_PENDING_PPM_WEIGHT';
export const CREATE_OR_UPDATE_PPM = ReduxHelpers.generateAsyncActionTypes(
  'CREATE_OR_UPDATE_PPM',
);
export const GET_PPM = ReduxHelpers.generateAsyncActionTypes('GET_PPM');
export const GET_PPM_ESTIMATE = ReduxHelpers.generateAsyncActionTypes(
  'GET_PPM_ESTIMATE',
);

function formatPpmEstimate(estimate) {
  // Range values arrive in cents, so convert to dollars
  return `$${(estimate.range_min / 100).toFixed(2)} - ${(
    estimate.range_max / 100
  ).toFixed(2)}`;
}

// Action creation
export function setPendingPpmSize(value) {
  return { type: SET_PENDING_PPM_SIZE, payload: value };
}

export function setPendingPpmWeight(value) {
  return { type: SET_PENDING_PPM_WEIGHT, payload: value };
}

export function getPpmWeightEstimate(
  moveDate,
  originZip,
  destZip,
  weightEstimate,
) {
  const action = ReduxHelpers.generateAsyncActions('GET_PPM_ESTIMATE');
  return function(dispatch, getState) {
    dispatch(action.start());
    GetPpmWeightEstimate(moveDate, originZip, destZip, weightEstimate)
      .then(item => dispatch(action.success(item)))
      .catch(error => dispatch(action.error(error)));
  };
}

export function createOrUpdatePpm(moveId, ppm) {
  const action = ReduxHelpers.generateAsyncActions('CREATE_OR_UPDATE_PPM');
  return function(dispatch, getState) {
    dispatch(action.start());
    const state = getState();
    const currentPpm = state.ppm.currentPpm;
    if (currentPpm) {
      UpdatePpm(moveId, currentPpm.id, ppm)
        .then(item =>
          dispatch(action.success(Object.assign({}, currentPpm, item))),
        )
        .catch(error => dispatch(action.error(error)));
    } else {
      CreatePpm(moveId, ppm)
        .then(item => dispatch(action.success(item)))
        .catch(error => dispatch(action.error(error)));
    }
  };
}

export function loadPpm(moveId) {
  const action = ReduxHelpers.generateAsyncActions('GET_PPM');
  return function(dispatch, getState) {
    dispatch(action.start);
    const state = getState();
    const currentPpm = state.ppm.currentPpm;
    if (!currentPpm) {
      // Load PPM from loggedInUser if available
      const loadedMoves = get(
        state,
        'loggedInUser.loggedInUser.service_member.orders.0.moves',
        [],
      );
      const matchingMove = find(loadedMoves, ['id', moveId]);
      if (get(matchingMove, 'personally_procured_moves.length')) {
        dispatch(action.success(matchingMove.personally_procured_moves));
      } else {
        GetPpm(moveId)
          .then(item => dispatch(action.success(item)))
          .catch(error => dispatch(action.error(error)));
      }
    }
  };
}
// Reducer
const initialState = {
  pendingPpmSize: null,
  incentive: null,
  pendingPpmWeight: null,
  currentPpm: null,
  hasSubmitError: false,
  hasSubmitSuccess: false,
  hasLoadSuccess: false,
  hasLoadError: false,
  hasEstimateSuccess: false,
  hasEstimateError: false,
  hasEstimateInProgress: false,
};
export function ppmReducer(state = initialState, action) {
  switch (action.type) {
    case SET_PENDING_PPM_SIZE:
      return Object.assign({}, state, {
        pendingPpmSize: action.payload,
      });
    case SET_PENDING_PPM_WEIGHT:
      return Object.assign({}, state, {
        pendingPpmWeight: action.payload,
      });
    case CREATE_OR_UPDATE_PPM.start:
      return Object.assign({}, state, {
        hasSubmitSuccess: false,
      });
    case CREATE_OR_UPDATE_PPM.success:
      return Object.assign({}, state, {
        currentPpm: action.payload,
        pendingPpmSize: null,
        pendingPpmWeight: null,
        hasSubmitSuccess: true,
        hasSubmitError: false,
      });
    case CREATE_OR_UPDATE_PPM.failure:
      return Object.assign({}, state, {
        hasSubmitSuccess: false,
        hasSubmitError: true,
        error: action.error,
      });
    case GET_PPM.start:
      return Object.assign({}, state, {
        hasLoadSuccess: false,
      });
    case GET_PPM.success:
      return Object.assign({}, state, {
        currentPpm: get(action.payload, '0', null),
        pendingPpmWeight: get(action.payload, '0.weight_estimate', null),
        incentive: get(action.payload, '0.estimated_incentive', null),
        hasLoadSuccess: true,
        hasLoadError: false,
      });
    case GET_PPM.failure:
      return Object.assign({}, state, {
        currentPpm: null,
        hasLoadSuccess: false,
        hasLoadError: true,
        error: action.error,
      });
    case GET_PPM_ESTIMATE.start:
      return Object.assign({}, state, {
        hasEstimateSuccess: false,
        hasEstimateInProgress: true,
      });
    case GET_PPM_ESTIMATE.success:
      return Object.assign({}, state, {
        incentive: formatPpmEstimate(action.payload),
        hasEstimateSuccess: true,
        hasEstimateError: false,
        hasEstimateInProgress: false,
        error: null,
      });
    case GET_PPM_ESTIMATE.failure:
      return Object.assign({}, state, {
        incentive: null,
        hasEstimateSuccess: false,
        hasEstimateError: true,
        hasEstimateInProgress: false,
        error: action.error,
      });
    default:
      return state;
  }
}
