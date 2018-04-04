import { CreatePpm, UpdatePpm, GetPpm } from './api.js';
import * as ReduxHelpers from 'shared/ReduxHelpers';

// Types
export const SET_PENDING_PPM_SIZE = 'SET_PENDING_PPM_SIZE';
export const SET_PENDING_PPM_WEIGHT = 'SET_PENDING_PPM_WEIGHT';
export const CREATE_OR_UPDATE_PPM = ReduxHelpers.generateAsyncActionTypes(
  'CREATE_OR_UPDATE_PPM',
);
export const GET_INCENTIVE = 'GET_INCENTIVE'; //TOOD: this should be async when rate engine is available
export const GET_PPM = ReduxHelpers.generateAsyncActionTypes('GET_PPM');

// Action creation
export function setPendingPpmSize(value) {
  return { type: SET_PENDING_PPM_SIZE, payload: value };
}

export function setPendingPpmWeight(value) {
  return { type: SET_PENDING_PPM_WEIGHT, payload: value };
}

export function getIncentive(weight) {
  // todo: this will probably need more information for real rate engince
  return {
    type: GET_INCENTIVE,
    payload: `$${0.75 * weight} - $${1.15 * weight}`,
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
      GetPpm(moveId)
        .then(item => dispatch(action.success(item)))
        .catch(error => dispatch(action.error(error)));
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
    case GET_INCENTIVE:
      return Object.assign({}, state, {
        incentive: action.payload,
      });
    case CREATE_OR_UPDATE_PPM.success:
      return Object.assign({}, state, {
        currentPpm: action.item,
        pendingPpmSize: null,
        hasSubmitSuccess: true,
        hasSubmitError: false,
      });
    case CREATE_OR_UPDATE_PPM.failure:
      return Object.assign({}, state, {
        currentPpm: null,
        hasSubmitSuccess: false,
        hasSubmitError: true,
        error: action.error,
      });
    case GET_PPM.success:
      return Object.assign({}, state, {
        currentPpm: action.item,
        hasSubmitSuccess: true,
        hasSubmitError: false,
      });
    case GET_PPM.failure:
      return Object.assign({}, state, {
        currentPpm: null,
        hasSubmitSuccess: false,
        hasSubmitError: true,
        error: action.error,
      });
    default:
      return state;
  }
}
