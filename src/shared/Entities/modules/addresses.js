import { addresses } from '../schema';
import { ADD_ENTITIES } from '../actions';
import { denormalize } from 'normalizr';

export const STATE_KEY = 'addresses';

export default function reducer(state = {}, action) {
  switch (action.type) {
    case ADD_ENTITIES:
      return {
        ...state,
        ...action.payload.addresses,
      };

    default:
      return state;
  }
}

export const selectAddress = (state, id) => {
  return denormalize([id], addresses, state.entities)[0];
};
