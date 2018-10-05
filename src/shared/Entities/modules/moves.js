import { denormalize } from 'normalizr';
import { swaggerRequest } from 'shared/Swagger/request';
import { getClient } from 'shared/Swagger/api';

import { moves } from '../schema';
import { ADD_ENTITIES } from '../actions';

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

export function getMove(label, moveId) {
  return swaggerRequest(getClient, 'moves.showMove', { moveId }, { label });
}

export const selectMove = (state, id) => {
  return denormalize([id], moves, state.entities)[0];
};
