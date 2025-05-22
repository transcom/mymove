import { LOG_OUT, SET_ACTIVE_OFFICE, SET_ACTIVE_ROLE } from './actions';

import { officeRoles } from 'constants/userRoles';

export const initialState = {
  activeRole: null,
  activeOffice: null,
  isLoggedIn: false,
  hasSucceeded: false,
  hasErrored: false,
  isLoading: true,
  underMaintenance: false,
};

const authReducer = (state = initialState, action = {}) => {
  switch (action?.type) {
    case LOG_OUT: {
      return initialState;
    }
    case 'SET_UNDER_MAINTENANCE': {
      return {
        ...state,
        isLoading: false,
        underMaintenance: true,
      };
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
      const officeUser = action?.payload?.office_user;
      const firstOfficeRole = roles?.find((r) => officeRoles.indexOf(r.roleType) > -1)?.roleType;
      const primaryOffice = officeUser?.transportation_office_assignments?.find(
        (office) => office.primaryOffice === true,
      );

      return {
        ...state,
        activeRole: firstOfficeRole,
        activeOffice: primaryOffice?.transportationOffice,
        hasSucceeded: true,
        hasErrored: false,
        isLoading: false,
        isLoggedIn: true,
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
    case SET_ACTIVE_OFFICE: {
      return {
        ...state,
        activeOffice: action.payload,
      };
    }

    default:
      return state;
  }
};

export default authReducer;
