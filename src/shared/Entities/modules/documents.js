import { documents } from '../schema';
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

// Selectors
export const selectDocument = (state, id) => {
  if (!id) {
    return {};
  }
  return denormalize([id], documents, state.entities)[0];
};
