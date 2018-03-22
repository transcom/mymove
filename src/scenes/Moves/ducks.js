import { CreateMove } from './api.js';

// Types
export const SET_PENDING_MOVE_TYPE = 'SET_PENDING_MOVE_TYPE';
export const CREATE_MOVE = 'CREATE_MOVE';
export const CREATE_MOVE_SUCCESS = 'CREATE_MOVE_SUCCESS';
export const CREATE_MOVE_FAILURE = 'CREATE_MOVE_FAILURE';

// creating move
export const createMoveRequest = () => ({
  type: CREATE_MOVE,
});

export const createMoveSuccess = item => ({
  type: CREATE_MOVE_SUCCESS,
  item,
});

export const createMoveFailure = error => ({
  type: CREATE_MOVE_FAILURE,
  error,
});

// Action creation
export function setPendingMoveType(value) {
  return { type: SET_PENDING_MOVE_TYPE, payload: value };
}

export function createMove(value) {
  return function(dispatch, getState) {
    dispatch(createMoveRequest());
    CreateMove(value)
      .then(item => dispatch(createMoveSuccess(item)))
      .catch(error => dispatch(createMoveFailure(error)));
  };
}

// Reducer
const initialState = {
  currentMove: null,
  pendingMoveType: null,
  hasSubmitError: false,
  hasSubmitSuccess: false,
};
export function moveReducer(state = initialState, action) {
  switch (action.type) {
    case SET_PENDING_MOVE_TYPE:
      return Object.assign({}, state, {
        pendingMoveType: action.payload,
      });
    case CREATE_MOVE_SUCCESS:
      return Object.assign({}, state, {
        currentMove: action.item,
        pendingMoveType: null,
        hasSubmitSuccess: true,
        hasSubmitError: false,
      });
    case CREATE_MOVE_FAILURE:
      return Object.assign({}, state, {
        currentMove: {},
        hasSubmitSuccess: false,
        hasSubmitError: true,
      });
    default:
      return state;
  }
}
