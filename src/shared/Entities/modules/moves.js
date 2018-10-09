import { get } from 'lodash';
import { moves } from '../schema';
import { ADD_ENTITIES } from '../actions';
import { denormalize } from 'normalizr';
import { swaggerRequest } from 'shared/Swagger/request';
import { getClient } from 'shared/Swagger/api';

export const STATE_KEY = 'moves';

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

export const selectMove = (state, id) => {
  return denormalize([id], moves, state.entities)[0];
};

export function getMoveDatesSummary(label, moveId, move_date) {
  return swaggerRequest(
    getClient,
    'moves.showMoveDatesSummary',
    { moveId, move_date },
    { label },
  );
}

export function selectMoveDatesSummary(state, idAttribute) {
  if (!idAttribute) {
    return null;
  }
  return get(state, `entities.moveDatesSummaries.${idAttribute}`);
}
