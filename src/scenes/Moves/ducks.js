import { get } from 'lodash';
import {
  CreateMove,
  UpdateMove,
  GetMove,
  SubmitMoveForApproval,
} from './api.js';
import { GET_LOGGED_IN_USER } from 'shared/User/ducks';

import * as ReduxHelpers from 'shared/ReduxHelpers';
// Types
const SET_PENDING_MOVE_TYPE = 'SET_PENDING_MOVE_TYPE';

export const getMoveType = 'GET_MOVE';
export const GET_MOVE = ReduxHelpers.generateAsyncActionTypes(getMoveType);

export const createOrUpdateMoveType = 'CREATE_OR_UPDATE_MOVE';
export const CREATE_OR_UPDATE_MOVE = ReduxHelpers.generateAsyncActionTypes(
  createOrUpdateMoveType,
);

export const submitForApprovalType = 'SUBMIT_FOR_APPROVAL';
export const SUBMIT_FOR_APPROVAL = ReduxHelpers.generateAsyncActionTypes(
  submitForApprovalType,
);

// Action creation
export function setPendingMoveType(value) {
  return { type: SET_PENDING_MOVE_TYPE, payload: value };
}

export function createMove(ordersId, movePayload = {}) {
  return function(dispatch) {
    const action = ReduxHelpers.generateAsyncActions(createOrUpdateMoveType);
    dispatch(action.start());
    return CreateMove(ordersId, movePayload)
      .then(item => dispatch(action.success(item)))
      .catch(error => dispatch(action.failure(error)));
  };
}

export function updateMove(moveId, moveType) {
  return function(dispatch) {
    const action = ReduxHelpers.generateAsyncActions(createOrUpdateMoveType);
    dispatch(action.start());
    return UpdateMove(moveId, { selected_move_type: moveType })
      .then(item => dispatch(action.success(item)))
      .catch(error => dispatch(action.failure(error)));
  };
}

export function loadMove(moveId) {
  return function(dispatch, getState) {
    const action = ReduxHelpers.generateAsyncActions(getMoveType);
    dispatch(action.start());
    return GetMove(moveId)
      .then(item => dispatch(action.success(item)))
      .catch(error => dispatch(action.success(error)));
  };
}

export const SubmitForApproval = ReduxHelpers.generateAsyncActionCreator(
  submitForApprovalType,
  SubmitMoveForApproval,
);

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
    case GET_LOGGED_IN_USER.success:
      const move = get(action.payload, 'service_member.orders.0.moves.0', null);
      return Object.assign({}, state, {
        currentMove: move,
      });
    case SET_PENDING_MOVE_TYPE:
      return Object.assign({}, state, {
        pendingMoveType: action.payload,
      });
    case CREATE_OR_UPDATE_MOVE.success:
      return Object.assign({}, state, {
        currentMove: action.payload,
        pendingMoveType: null,
        hasSubmitSuccess: true,
        hasSubmitError: false,
        error: null,
      });
    case CREATE_OR_UPDATE_MOVE.failure:
      return Object.assign({}, state, {
        currentMove: {},
        hasSubmitSuccess: false,
        hasSubmitError: true,
        error: action.error,
      });
    case GET_MOVE.success:
      return Object.assign({}, state, {
        currentMove: action.payload,
        hasSubmitSuccess: true,
        hasSubmitError: false,
        error: null,
      });
    case GET_MOVE.failure:
      return Object.assign({}, state, {
        currentMove: {},
        hasSubmitSuccess: false,
        hasSubmitError: true,
        error: action.error,
      });
    case SUBMIT_FOR_APPROVAL.start:
      return Object.assign({}, state, {
        submittedForApproval: false,
      });
    case SUBMIT_FOR_APPROVAL.success:
      return Object.assign({}, state, {
        currentMove: action.payload,
        submittedForApproval: true,
      });
    case SUBMIT_FOR_APPROVAL.failure:
      return Object.assign({}, state, {
        submittedForApproval: false,
        error: action.error,
      });
    default:
      return state;
  }
}
