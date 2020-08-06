export const SET_ACTIVE_ROLE = 'SET_ACTIVE_ROLE';

export const setActiveRole = (roleType) => ({
  type: SET_ACTIVE_ROLE,
  payload: roleType,
});

export const LOAD_USER = 'LOAD_USER';

export const loadUser = () => ({
  type: LOAD_USER,
});

export const LOG_OUT = 'LOG_OUT';

export const logOut = () => ({
  type: LOG_OUT,
});
