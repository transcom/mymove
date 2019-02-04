import { get } from 'lodash';
import { moves } from '../schema';
import { ADD_ENTITIES } from '../actions';
import { denormalize } from 'normalizr';
import { swaggerRequest } from 'shared/Swagger/request';
import { getClient } from 'shared/Swagger/api';

export const STATE_KEY = 'moves';
const approveBasicsLabel = 'Moves.ApproveBasics';
const cancelMoveLabel = 'Moves.CancelMove';

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

export function getMove(label, moveId) {
  return swaggerRequest(getClient, 'moves.showMove', { moveId }, { label });
}

export const selectMove = (state, id) => {
  return denormalize([id], moves, state.entities)[0];
};

export function getMoveDatesSummary(label, moveId, moveDate) {
  return swaggerRequest(getClient, 'moves.showMoveDatesSummary', { moveId, moveDate }, { label });
}

export function selectMoveDatesSummary(state, moveId, moveDate) {
  if (!moveId || !moveDate) {
    return null;
  }
  return get(state, `entities.moveDatesSummaries.${moveId}:${moveDate}`);
}

export function approveBasics(moveId) {
  const label = approveBasicsLabel;
  const swaggerTag = 'office.approveMove';
  return swaggerRequest(getClient, swaggerTag, { moveId }, { label });
}

export function cancelMove(moveId, cancelReason) {
  const label = cancelMoveLabel;
  const swaggerTag = 'office.cancelMove';
  const cancelMove = { cancel_reason: cancelReason };
  return swaggerRequest(getClient, swaggerTag, { moveId, cancelMove }, { label });
}

export function selectMoveStatus(state, moveId) {
  const entitiesMove = get(state, `entities.moves.${moveId}`);
  if (entitiesMove) {
    return entitiesMove.status;
  } else {
    return get(state, 'office.officeMove.status', '');
  }
}
