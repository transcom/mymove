import { isNull, get, isEmpty } from 'lodash';
import { moves } from '../schema';
import { ADD_ENTITIES } from '../actions';
import { denormalize } from 'normalizr';
import { swaggerRequest } from 'shared/Swagger/request';
import { getClient } from 'shared/Swagger/api';
import { selectEntitlements } from 'shared/entitlements.js';
import { selectOrdersForMove, selectActiveOrLatestOrders } from 'shared/Entities/modules/orders';
import { selectServiceMemberForMove } from 'shared/Entities/modules/serviceMembers';
import { getGHCClient } from 'shared/Swagger/api';
import { filter } from 'lodash';
import { fetchActive } from 'shared/utils';

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

export function loadMove(moveId, label = loadMoveLabel) {
  return swaggerRequest(getClient, 'moves.showMove', { moveId }, { label });
}

export function getMoveDatesSummary(moveId, moveDate, label = getMoveDatesSummaryLabel) {
  return swaggerRequest(getClient, 'moves.showMoveDatesSummary', { moveId, moveDate }, { label });
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
  const serviceMember = selectServiceMemberForMove(state, moveId);
  const weightAllotment = serviceMember.weight_allotment;
  if (isNull(hasDependents) || isNull(spouseHasProGear) || isNull(weightAllotment)) {
    return null;
  }
  return selectEntitlements(weightAllotment, hasDependents, spouseHasProGear);
}

// Selectors
export function selectMoveDatesSummary(state, moveId, moveDate) {
  if (!moveId || !moveDate) {
    return null;
  }
  return get(state, `entities.moveDatesSummaries.${moveId}:${moveDate}`);
}

export const selectMove = (state, id) => {
  const emptyMove = {};
  if (!id) return emptyMove;
  return denormalize([id], moves, state.entities)[0] || emptyMove;
};

export function selectMoveByLocator(state, locator) {
  const moves = filter(state.entities.moves, (move) => move.locator === locator);
  return moves[0];
}

export function selectActiveMoveByOrdersId(state, ordersId) {
  let emptymove = {};
  const move = fetchActive(filter(state.entities.moves, (move) => move.orders_id === ordersId));
  return move || emptymove;
}

export function selectMoveStatus(state, moveId) {
  const move = selectMove(state, moveId);
  return move.status;
}

export function selectActiveOrLatestMove(state) {
  // temp until full redux refactor: gets active (or latest move) from entities if it exists.  If not, gets it from currentMove
  let activeOrLatestOrders = selectActiveOrLatestOrders(state);
  if (isEmpty(activeOrLatestOrders)) {
    return {};
  }

  // get move from entities if it's there
  let move = selectActiveMoveByOrdersId(state, activeOrLatestOrders.id);
  if (isEmpty(move)) {
    move = get(state, 'moves.currentMove') || get(state, 'moves.latestMove') || {};
    return move;
  }

  return move;
}
