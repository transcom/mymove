import { isNull } from 'lodash';
import { denormalize } from 'normalizr';

import { moves } from '../schema';
import { ADD_ENTITIES } from '../actions';

import { swaggerRequest } from 'shared/Swagger/request';
import { getClient, getGHCClient } from 'shared/Swagger/api';
import { selectEntitlements } from 'shared/entitlements.js';
import { selectOrdersForMove } from 'shared/Entities/modules/orders';

/** REMAINING EXPORTS ARE USED BY PPM OFFICE */

export const STATE_KEY = 'moves';
const approveBasicsLabel = 'Moves.ApproveBasics';
const cancelMoveLabel = 'Moves.CancelMove';
export const loadMoveLabel = 'Moves.loadMove';
export const getMoveDatesSummaryLabel = 'Moves.getMoveDatesSummary';
export const getMoveByLocatorOperation = 'move.getMove';

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

export function calculateEntitlementsForMove(state, moveId) {
  const orders = selectOrdersForMove(state, moveId);
  const hasDependents = orders.has_dependents;
  const spouseHasProGear = orders.spouse_has_pro_gear;
  if (isNull(hasDependents) || isNull(spouseHasProGear) || isNull(orders.authorizedWeight)) {
    return null;
  }
  return selectEntitlements(orders, hasDependents, spouseHasProGear);
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
