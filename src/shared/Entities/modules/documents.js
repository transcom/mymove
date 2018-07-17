import { documentModel } from '../schema';
import { ADD_ENTITIES } from '../actions';
import { denormalize } from 'normalizr';

export const STATE_KEY = 'documents';

export default function reducer(state = {}, action) {
  switch (action.type) {
    case ADD_ENTITIES:
      return {
        ...state,
        ...action.payload.documents,
      };

    default:
      return state;
  }
}

export const selectHydrated = (state, id) => {
  return denormalize([id], documentModel, state.entities)[0];
};
