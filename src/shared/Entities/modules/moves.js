import { denormalize } from 'normalizr';

import { moves } from '../schema';
import { ADD_ENTITIES } from '../actions';

import { swaggerRequest } from 'shared/Swagger/request';
import { getClient, getGHCClient } from 'shared/Swagger/api';

/** REMAINING EXPORTS ARE USED BY PPM OFFICE */

export const STATE_KEY = 'moves';
const approveBasicsLabel = 'Moves.ApproveBasics';
const cancelMoveLabel = 'Moves.CancelMove';
export const loadMoveLabel = 'Moves.loadMove';
export const getMoveDatesSummaryLabel = 'Moves.getMoveDatesSummary';
export const getMoveByLocatorOperation = 'move.getMove';
export const getAllMovesLabel = 'move.getAllMoves';

export default function reducer(state = {}, action) {
  switch (action.type) {
    case ADD_ENTITIES:
      return {
        ...state,
        ...action.payload.moves,
      };

    default:
      return state;
  }
}

export function getAllMoves(serviceMemberId, label = getAllMovesLabel) {
  const swaggerTag = 'moves.getAllMoves';
  return swaggerRequest(getClient, swaggerTag, { serviceMemberId }, { label });
}

export function getMoveByLocator(locator, label = getMoveByLocatorOperation) {
  return swaggerRequest(getGHCClient, getMoveByLocatorOperation, { locator }, { label });
}

// TODO - migrate
export function loadMove(moveId, label = loadMoveLabel) {
  return swaggerRequest(getClient, 'moves.showMove', { moveId }, { label });
}

export function approveBasics(moveId, label = approveBasicsLabel) {
  const swaggerTag = 'office.approveMove';
  return swaggerRequest(getClient, swaggerTag, { moveId }, { label });
}

export function cancelMove(moveId, cancelReason, label = cancelMoveLabel) {
  const swaggerTag = 'office.cancelMove';
  const cancelMove = { cancel_reason: cancelReason };
  return swaggerRequest(getClient, swaggerTag, { moveId, cancelMove }, { label });
}

// Selectors
export const selectMove = (state, id) => {
  const emptyMove = {};
  if (!id) return emptyMove;
  return denormalize([id], moves, state.entities)[0] || emptyMove;
};

export const selectMoveByLocator = (state, locator) => {
  return Object.values(state.entities.moves).find((move) => move.locator === locator);
};

export function selectMoveStatus(state, moveId) {
  const move = selectMove(state, moveId);
  return move.status;
}
