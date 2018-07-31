import { uploads } from '../schema';
import { ADD_ENTITIES } from '../actions';
import { denormalize } from 'normalizr';

export const STATE_KEY = 'uploads';

export default function reducer(state = {}, action) {
  switch (action.type) {
    case ADD_ENTITIES:
      return {
        ...state,
        ...action.payload.uploads,
      };

    default:
      return state;
  }
}

export const selectUpload = (state, id) => {
  return denormalize([id], uploads, state.entities)[0];
};
