import { moves } from '../schema';
import { ADD_ENTITIES } from '../actions';
import { denormalize } from 'normalizr';

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
