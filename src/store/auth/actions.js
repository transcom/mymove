export const SET_ACTIVE_ROLE = 'SET_ACTIVE_ROLE';
export const SET_ACTIVE_ROLE_SUCCESS = 'SET_ACTIVE_ROLE_SUCCESS';
export const SET_ACTIVE_ROLE_FAILURE = 'SET_ACTIVE_ROLE_FAILURE';

export const setActiveRole = (roleType) => ({
  type: SET_ACTIVE_ROLE,
  payload: roleType,
});

export const setActiveRoleSuccess = (roleType) => ({
  type: SET_ACTIVE_ROLE_SUCCESS,
  payload: roleType,
});

export const setActiveRoleFailure = (error) => ({
  type: SET_ACTIVE_ROLE_FAILURE,
  error,
});

export const LOAD_USER = 'LOAD_USER';

export const loadUser = () => ({
  type: LOAD_USER,
});

export const LOG_OUT = 'LOG_OUT';

export const logOut = () => ({
  type: LOG_OUT,
});

export const GET_LOGGED_IN_USER_START = 'GET_LOGGED_IN_USER_START';

export const getLoggedInUserStart = () => ({
  type: GET_LOGGED_IN_USER_START,
});

export const GET_LOGGED_IN_USER_SUCCESS = 'GET_LOGGED_IN_USER_SUCCESS';

export const getLoggedInUserSuccess = (user) => ({
  type: GET_LOGGED_IN_USER_SUCCESS,
  payload: user,
});

export const GET_LOGGED_IN_USER_FAILURE = 'GET_LOGGED_IN_USER_FAILURE';

export const getLoggedInUserFailure = (error) => ({
  type: GET_LOGGED_IN_USER_FAILURE,
  error,
});

export const SET_UNDER_MAINTENANCE = 'SET_UNDER_MAINTENANCE';

export const setUnderMaintenance = () => ({
  type: SET_UNDER_MAINTENANCE,
});
