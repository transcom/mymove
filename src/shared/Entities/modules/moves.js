import { isNull, get } from 'lodash';
import { moves } from '../schema';
import { ADD_ENTITIES } from '../actions';
import { denormalize } from 'normalizr';
import { swaggerRequest } from 'shared/Swagger/request';
import { getClient } from 'shared/Swagger/api';
import { getEntitlements } from 'shared/entitlements.js';
import { selectOrdersForMove } from 'shared/Entities/modules/orders';
import { selectServiceMemberForMove } from 'shared/Entities/modules/serviceMembers';

export const STATE_KEY = 'moves';
const approveBasicsLabel = 'Moves.ApproveBasics';
const cancelMoveLabel = 'Moves.CancelMove';
export const loadMoveLabel = 'Moves.loadMove';
export const getMoveDatesSummaryLabel = 'Moves.getMoveDatesSummary';

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
  const rank = serviceMember.rank;
  if (isNull(hasDependents) || isNull(spouseHasProGear) || isNull(rank)) {
    return null;
  }
  return getEntitlements(rank, hasDependents, spouseHasProGear);
}

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

export function selectMoveStatus(state, moveId) {
  const move = selectMove(state, moveId);
  return move.status;
}
