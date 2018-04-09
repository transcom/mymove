import { CreateMove, UpdateMove, GetMove } from './api.js';

// Types
export const SET_PENDING_MOVE_TYPE = 'SET_PENDING_MOVE_TYPE';
export const CREATE_MOVE = 'CREATE_MOVE';
export const UPDATE_MOVE = 'UPDATE_MOVE';
export const CREATE_OR_UPDATE_MOVE_SUCCESS = 'CREATE_OR_UPDATE_MOVE_SUCCESS';
export const CREATE_OR_UPDATE_MOVE_FAILURE = 'CREATE_OR_UPDATE_MOVE_FAILURE';
export const GET_MOVE = 'GET_MOVE';
export const GET_MOVE_SUCCESS = 'GET_MOVE_SUCCESS';
export const GET_MOVE_FAILURE = 'GET_MOVE_FAILURE';

export const createMoveRequest = () => ({
  type: CREATE_MOVE,
});

export const updateMoveRequest = () => ({
  type: UPDATE_MOVE,
});

export const createOrUpdateMoveSuccess = item => ({
  type: CREATE_OR_UPDATE_MOVE_SUCCESS,
  item,
});

export const createOrUpdateMoveFailure = error => ({
  type: CREATE_OR_UPDATE_MOVE_FAILURE,
  error,
});

const getMoveRequest = () => ({
  type: GET_MOVE,
});

export const getMoveSuccess = item => ({
  type: GET_MOVE_SUCCESS,
  item,
  // item: items.length > 0 ? items[0] : null,
});

export const getMoveFailure = error => ({
  type: GET_MOVE_FAILURE,
  error,
});

// Action creation
export function setPendingMoveType(value) {
  return { type: SET_PENDING_MOVE_TYPE, payload: value };
}

export function createMove(moveType) {
  return function(dispatch, getState) {
    dispatch(createMoveRequest());
    CreateMove(moveType)
      .then(item => dispatch(createOrUpdateMoveSuccess(item)))
      .catch(error => dispatch(createOrUpdateMoveFailure(error)));
  };
}

export function updateMove(moveId, moveType) {
  return function(dispatch, getState) {
    dispatch(updateMoveRequest());
    UpdateMove(moveId, { selected_move_type: moveType })
      .then(item => dispatch(createOrUpdateMoveSuccess(item)))
      .catch(error => dispatch(createOrUpdateMoveFailure(error)));
  };
}

export function loadMove(moveId) {
  return function(dispatch, getState) {
    dispatch(getMoveRequest());
    const state = getState();
    const currentMove = state.submittedMoves.currentMove;
    if (!currentMove) {
      GetMove(moveId)
        .then(item => dispatch(getMoveSuccess(item)))
        .catch(error => dispatch(getMoveFailure(error)));
    }
  };
}

// Reducer
const initialState = {
  currentMove: null,
  pendingMoveType: null,
  hasSubmitError: false,
  hasSubmitSuccess: false,
  error: null,
};
export function moveReducer(state = initialState, action) {
  switch (action.type) {
    case SET_PENDING_MOVE_TYPE:
      return Object.assign({}, state, {
        pendingMoveType: action.payload,
      });
    case UPDATE_MOVE:
      return Object.assign({}, state, {
        hasSubmitSuccess: false,
      });
    case CREATE_OR_UPDATE_MOVE_SUCCESS:
      return Object.assign({}, state, {
        currentMove: action.item,
        pendingMoveType: null,
        hasSubmitSuccess: true,
        hasSubmitError: false,
        error: null,
      });
    case CREATE_OR_UPDATE_MOVE_FAILURE:
      return Object.assign({}, state, {
        currentMove: {},
        hasSubmitSuccess: false,
        hasSubmitError: true,
        error: action.error,
      });
    case GET_MOVE_SUCCESS:
      return Object.assign({}, state, {
        currentMove: action.item,
        hasSubmitSuccess: true,
        hasSubmitError: false,
        error: null,
      });
    case GET_MOVE_FAILURE:
      return Object.assign({}, state, {
        currentMove: {},
        hasSubmitSuccess: false,
        hasSubmitError: true,
        error: action.error,
      });
    default:
      return state;
  }
}
