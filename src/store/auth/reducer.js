import { SET_ACTIVE_ROLE, LOG_OUT } from './actions';

import { officeRoles } from 'constants/userRoles';

export const initialState = {
  activeRole: null,
};

const authReducer = (state = initialState, action) => {
  switch (action?.type) {
    case LOG_OUT: {
      return initialState;
    }

    case 'GET_LOGGED_IN_USER_SUCCESS': {
      if (state.activeRole) return state;

      const {
        payload: { roles = [] },
      } = action;
      const firstOfficeRole = roles?.find((r) => officeRoles.indexOf(r.roleType) > -1)?.roleType;

      return {
        ...state,
        activeRole: firstOfficeRole,
      };
    }
    case SET_ACTIVE_ROLE: {
      return {
        ...state,
        activeRole: action.payload,
      };
    }

    default:
      return state;
  }
};

export default authReducer;
