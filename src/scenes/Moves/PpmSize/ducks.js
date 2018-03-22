import { CreatePpm } from './api.js';

// Types
export const SET_PENDING_PPM_SIZE = 'SET_PENDING_PPM_SIZE';
export const CREATE_PPM = 'CREATE_PPM';
export const CREATE_PPM_SUCCESS = 'CREATE_PPM_SUCCESS';
export const CREATE_PPM_FAILURE = 'CREATE_PPM_FAILURE';

// Creating ppm
export const createPpmRequest = () => ({
  type: CREATE_PPM,
});

export const createPpmSuccess = item => ({
  type: CREATE_PPM_SUCCESS,
  item,
});

export const createPpmFailure = error => ({
  type: CREATE_PPM_FAILURE,
  error,
});

// Action creation
export function setPendingPpmSize(value) {
  return { type: SET_PENDING_PPM_SIZE, payload: value };
}
export function createPpm(moveId, size) {
  return function(dispatch) {
    dispatch(createPpmRequest());
    CreatePpm(moveId, size)
      .then(item => dispatch(createPpmSuccess(item)))
      .catch(error => dispatch(createPpmFailure(error)));
  };
}

// Reducer
const initialState = {
  pendingPpmSize: null,
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
    case CREATE_PPM_SUCCESS:
      return Object.assign({}, state, {
        currentPpm: action.item,
        hasSubmitSuccess: true,
        hasSubmitError: false,
      });
    case CREATE_PPM_FAILURE:
      return Object.assign({}, state, {
        currentPpm: {},
        hasSubmitSuccess: false,
        hasSubmitError: true,
      });
    default:
      return state;
  }
}
