import { SET_ACTIVE_ROLE } from './actions';

export const initialState = {
  activeRole: null,
};

const authReducer = (state = initialState, action) => {
  switch (action?.type) {
    case SET_ACTIVE_ROLE: {
      return {
        ...initialState,
        activeRole: action.payload,
      };
    }

    default:
      return state;
  }
};

export default authReducer;
