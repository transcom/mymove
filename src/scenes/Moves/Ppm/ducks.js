import { get, every, isNull, isNumber } from 'lodash';
import {
  CreatePpm,
  UpdatePpm,
  GetPpm,
  GetPpmWeightEstimate,
  GetPpmSitEstimate,
} from './api.js';
import * as ReduxHelpers from 'shared/ReduxHelpers';
import { GET_LOGGED_IN_USER } from 'shared/User/ducks';
import { loadEntitlementsFromState } from 'shared/entitlements';

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
export const GET_SIT_ESTIMATE = ReduxHelpers.generateAsyncActionTypes(
  'GET_SIT_ESTIMATE',
);

function formatSitEstimate(estimate) {
  // Range values arrive in cents, so convert to dollars
  return `$${(estimate / 100).toFixed(2)}`;
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
    return GetPpmWeightEstimate(moveDate, originZip, destZip, weightEstimate)
      .then(item => dispatch(action.success(item)))
      .catch(error => dispatch(action.error(error)));
  };
}

export function getPpmSitEstimate(
  moveDate,
  sitDays,
  originZip,
  destZip,
  weightEstimate,
) {
  const action = ReduxHelpers.generateAsyncActions('GET_SIT_ESTIMATE');
  const canEstimate = every([
    moveDate,
    sitDays,
    originZip,
    destZip,
    weightEstimate,
  ]);
  return function(dispatch, getState) {
    if (!canEstimate) {
      return dispatch(action.success({ estimate: null }));
    }
    dispatch(action.start());
    GetPpmSitEstimate(moveDate, sitDays, originZip, destZip, weightEstimate)
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
      return UpdatePpm(moveId, currentPpm.id, ppm)
        .then(item => dispatch(action.success(item)))
        .catch(error => dispatch(action.error(error)));
    } else {
      return CreatePpm(moveId, ppm)
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
      return GetPpm(moveId)
        .then(item => dispatch(action.success(item)))
        .catch(error => dispatch(action.error(error)));
    }
    return Promise.resolve();
  };
}

// Selectors
export function getRawWeightInfo(state) {
  const entitlement = loadEntitlementsFromState(state);
  if (isNull(entitlement)) {
    return null;
  }

  return {
    S: {
      min: 50,
      max: 1000,
    },
    M: {
      min: 500,
      max: 2500,
    },
    L: {
      min: 1500,
      max: entitlement.sum,
    },
  };
}

export function getMaxAdvance(state) {
  const maxIncentive = get(state, 'ppm.incentive_estimate_max');
  // we are using 20000000 since it is the largest number MacRae found that could be stored in table
  // and we don't want to block the user from requesting an advance if the rate engine fails
  return maxIncentive ? 0.6 * maxIncentive : 20000000;
}
export function getSelectedWeightInfo(state) {
  const weightInfo = getRawWeightInfo(state);
  const ppm = get(state, 'ppm.currentPpm', null);
  if (isNull(weightInfo) || isNull(ppm)) {
    return null;
  }

  const size = ppm ? ppm.size : 'L';
  return weightInfo[size];
}

// Reducer
const initialState = {
  pendingPpmSize: null,
  incentive: null,
  sitReimbursement: null,
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
    case GET_LOGGED_IN_USER.success:
      // Initialize state when we get the logged in user
      const user = action.payload;
      const currentPpm = get(
        user,
        'service_member.orders.0.moves.0.personally_procured_moves.0',
        null,
      );
      return Object.assign({}, state, {
        currentPpm: currentPpm,
        pendingPpmSize: get(currentPpm, 'size', null),
        pendingPpmWeight: get(currentPpm, 'weight_estimate', null),
        incentive_estimate_min: get(currentPpm, 'incentive_estimate_min', null),
        incentive_estimate_max: get(currentPpm, 'incentive_estimate_max', null),
        sitReimbursement: get(
          currentPpm,
          'estimated_storage_reimbursement',
          null,
        ),
        hasLoadSuccess: true,
        hasLoadError: false,
      });
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
        incentive_estimate_min: get(
          action.payload,
          'incentive_estimate_min',
          null,
        ),
        incentive_estimate_max: get(
          action.payload,
          'incentive_estimate_max',
          null,
        ),
        sitReimbursement: get(
          action.payload,
          'estimated_storage_reimbursement',
          null,
        ),
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
        incentive_estimate_min: get(
          action.payload,
          '0.incentive_estimate_min',
          null,
        ),
        incentive_estimate_max: get(
          action.payload,
          '0.incentive_estimate_max',
          null,
        ),
        sitReimbursement: get(
          action.payload,
          '0.estimated_storage_reimbursement',
          null,
        ),
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
        incentive_estimate_min: action.payload.range_min,
        incentive_estimate_max: action.payload.range_max,
        hasEstimateSuccess: true,
        hasEstimateError: false,
        hasEstimateInProgress: false,
        rateEngineError: null,
      });
    case GET_PPM_ESTIMATE.failure:
      return Object.assign({}, state, {
        incentive_estimate_min: null,
        incentive_estimate_max: null,
        hasEstimateSuccess: false,
        hasEstimateError: true,
        hasEstimateInProgress: false,
        rateEngineError: action.error,
        error: null,
      });
    case GET_SIT_ESTIMATE.start:
      return Object.assign({}, state, {
        hasEstimateSuccess: false,
        hasEstimateInProgress: true,
      });
    case GET_SIT_ESTIMATE.success:
      let estimate = null;
      if (isNumber(action.payload.estimate)) {
        // Convert from cents
        estimate = formatSitEstimate(action.payload.estimate);
      }
      return Object.assign({}, state, {
        sitReimbursement: estimate,
        hasEstimateSuccess: true,
        hasEstimateError: false,
        hasEstimateInProgress: false,
        rateEngineError: null,
      });
    case GET_SIT_ESTIMATE.failure:
      return Object.assign({}, state, {
        sitReimbursement: null,
        hasEstimateSuccess: false,
        hasEstimateError: true,
        hasEstimateInProgress: false,
        rateEngineError: action.error,
      });
    default:
      return state;
  }
}
