import { SET_ACTIVE_ROLE, LOG_OUT } from './actions';

import { officeRoles } from 'constants/userRoles';

export const initialState = {
  activeRole: null,
};

export const selectIsLoggedIn = (state) => {
  return state.auth?.isLoggedIn || null;
};

export function selectGetCurrentUserIsLoading(state) {
  return state.auth.isLoading;
}

export function selectGetCurrentUserIsSuccess(state) {
  return state.auth.hasSucceeded;
}

export function selectGetCurrentUserIsError(state) {
  return state.auth.hasErrored;
}

const authReducer = (state = initialState, action) => {
  switch (action?.type) {
    case LOG_OUT: {
      return initialState;
    }
    case 'GET_LOGGED_IN_USER_START': {
      return {
        ...state,
        hasSucceeded: false,
        hasErrored: false,
        isLoading: true,
      };
    }
    case 'GET_LOGGED_IN_USER_SUCCESS': {
      if (state.activeRole)
        return {
          ...state,
          hasSucceeded: true,
          hasErrored: false,
          isLoading: false,
          isLoggedIn: true,
        };

      const {
        payload: { roles = [] },
      } = action;
      const firstOfficeRole = roles?.find((r) => officeRoles.indexOf(r.roleType) > -1)?.roleType;

      return {
        ...state,
        activeRole: firstOfficeRole,
        hasSucceeded: true,
        hasErrored: false,
        isLoading: false,
      };
    }
    case 'GET_LOGGED_IN_USER_FAILURE': {
      return {
        ...state,
        isLoading: false,
        hasErrored: true,
        hasSucceeded: false,
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
