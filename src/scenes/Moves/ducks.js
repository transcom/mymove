import { CreateMove, UpdateMove } from './api.js';

// Types
export const SET_PENDING_MOVE_TYPE = 'SET_PENDING_MOVE_TYPE';
export const CREATE_OR_UPDATE_MOVE = 'CREATE_OR_UPDATE_MOVE';
export const CREATE_OR_UPDATE_MOVE_SUCCESS = 'CREATE_OR_UPDATE_MOVE_SUCCESS';
export const CREATE_OR_UPDATE_MOVE_FAILURE = 'CREATE_OR_UPDATE_MOVE_FAILURE';

// creating move
export const createOrUpdateMoveRequest = () => ({
  type: CREATE_OR_UPDATE_MOVE,
});

export const createOrUpdateMoveSuccess = item => ({
  type: CREATE_OR_UPDATE_MOVE_SUCCESS,
  item,
});

export const createOrUpdateMoveFailure = error => ({
  type: CREATE_OR_UPDATE_MOVE_FAILURE,
  error,
});

// Action creation
export function setPendingMoveType(value) {
  return { type: SET_PENDING_MOVE_TYPE, payload: value };
}

// TODO : add loadMove
export function createOrUpdateMove(moveId, moveType) {
  return function(dispatch, getState) {
    dispatch(createOrUpdateMoveRequest());
    const state = getState();
    const currentMove = state.submittedMoves.currentMove;
    if (currentMove) {
      console.log('update', currentMove);
      UpdateMove(moveId, { selected_move_type: moveType })
        .then(item => dispatch(createOrUpdateMoveSuccess(item)))
        .catch(error => dispatch(createOrUpdateMoveFailure(error)));
    } else {
      CreateMove(moveId, { moveType })
        .then(item => dispatch(createOrUpdateMoveSuccess(item)))
        .catch(error => dispatch(createOrUpdateMoveFailure(error)));
    }
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
    case CREATE_OR_UPDATE_MOVE_SUCCESS:
      return Object.assign({}, state, {
        currentMove: action.item,
        pendingMoveType: null,
        hasSubmitSuccess: true,
        hasSubmitError: false,
      });
    case CREATE_OR_UPDATE_MOVE_FAILURE:
      return Object.assign({}, state, {
        currentMove: {},
        hasSubmitSuccess: false,
        hasSubmitError: true,
      });
    default:
      return state;
  }
}
