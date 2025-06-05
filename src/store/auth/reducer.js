import {
  LOG_OUT,
  SET_ACTIVE_OFFICE,
  SET_ACTIVE_ROLE,
  SET_ACTIVE_ROLE_SUCCESS,
  SET_ACTIVE_ROLE_FAILURE,
} from './actions';

export const initialState = {
  activeRole: null,
  activeOffice: null,
  isLoggedIn: false,
  hasSucceeded: false,
  hasErrored: false,
  isLoading: true,
  underMaintenance: false,
  isSettingActiveRole: false,
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
      const {
        payload: { activeRole },
      } = action;
      const officeUser = action?.payload?.office_user;
      const primaryOffice = officeUser?.transportation_office_assignments?.find(
        (office) => office.primaryOffice === true,
      );

      return {
        ...state,
        activeRole: activeRole.roleType,
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
        isSettingActiveRole: true,
      };
    }
    case SET_ACTIVE_ROLE_SUCCESS: {
      return {
        ...state,
        activeRole: action.payload,
        isSettingActiveRole: false,
      };
    }
    case SET_ACTIVE_ROLE_FAILURE: {
      return {
        ...state,
        isSettingActiveRole: false,
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
