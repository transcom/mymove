import { get, every, isNull, isNumber, isEmpty } from 'lodash';
import { CreatePpm, UpdatePpm, GetPpm, GetPpmWeightEstimate, GetPpmSitEstimate, RequestPayment } from './api.js';
import * as ReduxHelpers from 'shared/ReduxHelpers';
import { GET_LOGGED_IN_USER } from 'shared/Data/users';
import { fetchActive } from 'shared/utils';
import { loadEntitlementsFromState } from 'shared/entitlements';
import { formatCents } from 'shared/formatters';
import { selectShipment } from 'shared/Entities/modules/shipments';
import { getCurrentShipmentID } from 'shared/UI/ducks';
import { change } from 'redux-form';

// Types
export const SET_PENDING_PPM_SIZE = 'SET_PENDING_PPM_SIZE';
export const SET_PENDING_PPM_WEIGHT = 'SET_PENDING_PPM_WEIGHT';
const CLEAR_SIT_ESTIMATE = 'CLEAR_SIT_ESTIMATE';
export const CREATE_OR_UPDATE_PPM = ReduxHelpers.generateAsyncActionTypes('CREATE_OR_UPDATE_PPM');
export const GET_PPM = ReduxHelpers.generateAsyncActionTypes('GET_PPM');
export const GET_PPM_ESTIMATE = ReduxHelpers.generateAsyncActionTypes('GET_PPM_ESTIMATE');
export const GET_SIT_ESTIMATE = ReduxHelpers.generateAsyncActionTypes('GET_SIT_ESTIMATE');

// Action creation
export function setPendingPpmSize(value) {
  return { type: SET_PENDING_PPM_SIZE, payload: value };
}

export function setPendingPpmWeight(value) {
  return { type: SET_PENDING_PPM_WEIGHT, payload: value };
}

export function getPpmWeightEstimate(moveDate, originZip, destZip, weightEstimate) {
  const action = ReduxHelpers.generateAsyncActions('GET_PPM_ESTIMATE');
  return function(dispatch, getState) {
    dispatch(action.start());
    return GetPpmWeightEstimate(moveDate, originZip, destZip, weightEstimate)
      .then(item => dispatch(action.success(item)))
      .catch(error => dispatch(action.error(error)));
  };
}

export function getPpmSitEstimate(moveDate, sitDays, originZip, destZip, weightEstimate) {
  const action = ReduxHelpers.generateAsyncActions('GET_SIT_ESTIMATE');
  const canEstimate = every([moveDate, sitDays, originZip, destZip, weightEstimate]);
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

export function clearPpmSitEstimate() {
  return { type: CLEAR_SIT_ESTIMATE };
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

export function setInitialFormValues(originalMoveDate, pickupPostalCode, destinationPostalCode) {
  return function(dispatch) {
    dispatch(change('ppp_date_and_location', 'original_move_date', originalMoveDate));
    dispatch(change('ppp_date_and_location', 'pickup_postal_code', pickupPostalCode));
    dispatch(change('ppp_date_and_location', 'destination_postal_code', destinationPostalCode));
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

const REQUESTED_PAYMENT_ACTION = {
  type: 'REQUESTED_PAYMENT',
};

export function submitExpenseDocs(state) {
  const updateAction = ReduxHelpers.generateAsyncActions('CREATE_OR_UPDATE_PPM');
  return function(dispatch, getState) {
    dispatch(updateAction.start());
    const state = getState();
    const currentPpm = state.ppm.currentPpm;
    if (!currentPpm) {
      console.log('Attempted to request payment on a PPM that did not exist.');
      return Promise.reject();
    }
    return RequestPayment(currentPpm.id)
      .then(item => {
        dispatch(updateAction.success(item));
        dispatch(REQUESTED_PAYMENT_ACTION);
      })
      .catch(error => {
        dispatch(updateAction.error(error));
        return Promise.reject();
      });
  };
}

// Selectors
export function getRawWeightInfo(state) {
  const entitlement = loadEntitlementsFromState(state);
  if (isEmpty(entitlement)) {
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
  return weightInfo[size]; // eslint-disable-line security/detect-object-injection
}

export function isHHGPPMComboMove(state) {
  return get(state, 'moves.currentMove.selected_move_type') === 'HHG_PPM';
}

const estimatedRemainingWeight = (sum, weight) => {
  if (sum >= weight) {
    return sum - weight;
  } else {
    return sum;
  }
};

export function getEstimatedRemainingWeight(state) {
  const entitlements = loadEntitlementsFromState(state);

  if (!isHHGPPMComboMove(state) || isNull(entitlements)) {
    return null;
  }

  const { sum } = entitlements;

  const { pm_survey_weight_estimate, weight_estimate } = selectShipment(state, getCurrentShipmentID(state));

  if (pm_survey_weight_estimate) {
    return estimatedRemainingWeight(sum, pm_survey_weight_estimate);
  }

  if (sum && weight_estimate >= 0) {
    return estimatedRemainingWeight(sum, weight_estimate);
  }
}

export function getActualRemainingWeight(state) {
  const entitlements = loadEntitlementsFromState(state);

  if (!isHHGPPMComboMove(state) || isNull(entitlements)) {
    return null;
  }

  const { sum } = entitlements;
  const { tare_weight, gross_weight } = selectShipment(state, getCurrentShipmentID(state));

  if (sum && gross_weight && tare_weight) {
    return estimatedRemainingWeight(sum, gross_weight - tare_weight);
  }
}

export function getDestinationPostalCode(state) {
  const currentShipment = selectShipment(state, getCurrentShipmentID(state));
  const addresses = state.entities.addresses;
  const currentOrders = state.orders.currentOrders;

  return currentShipment.has_delivery_address && addresses
    ? addresses[currentShipment.delivery_address].postal_code
    : currentOrders.new_duty_station.address.postal_code;
}

export function getPPM(state) {
  const move = state.moves.currentMove || state.moves.latestMove || {};
  const moveId = move.id;
  const ppmFromEntities = Object.values(state.entities.personallyProcuredMoves).find(ppm => ppm.move_id === moveId);
  return ppmFromEntities || state.ppm.currentPpm;
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
      const activeOrders = fetchActive(get(action.payload, 'service_member.orders'));
      const activeMove = fetchActive(get(activeOrders, 'moves'));
      const activePpm = fetchActive(get(activeMove, 'personally_procured_moves'));
      return Object.assign({}, state, {
        currentPpm: activePpm,
        pendingPpmSize: get(activePpm, 'size', null),
        pendingPpmWeight: get(activePpm, 'weight_estimate', null),
        incentive_estimate_min: get(activePpm, 'incentive_estimate_min', null),
        incentive_estimate_max: get(activePpm, 'incentive_estimate_max', null),
        sitReimbursement: get(activePpm, 'estimated_storage_reimbursement', null),
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
        hasSubmitInProgress: true,
      });
    case CREATE_OR_UPDATE_PPM.success:
      return Object.assign({}, state, {
        currentPpm: action.payload,
        incentive_estimate_min: get(action.payload, 'incentive_estimate_min', null),
        incentive_estimate_max: get(action.payload, 'incentive_estimate_max', null),
        sitReimbursement: get(action.payload, 'estimated_storage_reimbursement', null),
        pendingPpmSize: null,
        pendingPpmWeight: null,
        hasSubmitSuccess: true,
        hasSubmitError: false,
        hasSubmitInProgress: false,
      });
    case CREATE_OR_UPDATE_PPM.failure:
      return Object.assign({}, state, {
        hasSubmitSuccess: false,
        hasSubmitError: true,
        hasSubmitInProgress: false,
        error: action.error,
      });
    case 'REQUESTED_PAYMENT':
      return Object.assign({}, state, {
        requestPaymentSuccess: true,
      });
    case GET_PPM.start:
      return Object.assign({}, state, {
        hasLoadSuccess: false,
      });
    case GET_PPM.success:
      return Object.assign({}, state, {
        currentPpm: get(action.payload, '0', null),
        pendingPpmWeight: get(action.payload, '0.weight_estimate', null),
        incentive_estimate_min: get(action.payload, '0.incentive_estimate_min', null),
        incentive_estimate_max: get(action.payload, '0.incentive_estimate_max', null),
        sitReimbursement: get(action.payload, '0.estimated_storage_reimbursement', null),
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
        estimate = '$' + formatCents(action.payload.estimate);
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
    case CLEAR_SIT_ESTIMATE:
      return Object.assign({}, state, {
        sitReimbursement: null,
        hasEstimateSuccess: true,
        hasEstimateError: false,
        hasEstimateInProgress: false,
        rateEngineError: null,
      });
    default:
      return state;
  }
}
