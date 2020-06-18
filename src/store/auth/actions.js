export const SET_ACTIVE_ROLE = 'SET_ACTIVE_ROLE';

export const setActiveRole = (roleType) => ({
  type: SET_ACTIVE_ROLE,
  payload: roleType,
});
