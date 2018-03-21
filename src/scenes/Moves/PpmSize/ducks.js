import { CreatePpm } from './api.js';

// Types
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
export function createPpm(value) {
  return function(dispatch) {
    dispatch(createPpmRequest());
    CreatePpm(value)
      .then(item => dispatch(createPpmSuccess(item)))
      .catch(error => dispatch(createPpmFailure(error)));
  };
}

// Reducer
const initialState = {
  currentPpm: null,
  hasSubmitError: false,
  hasSubmitSuccess: false,
};
export function ppmReducer(state = initialState, action) {
  switch (action.type) {
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
