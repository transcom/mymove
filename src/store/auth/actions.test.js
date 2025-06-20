import {
  setActiveRole,
  SET_ACTIVE_ROLE,
  loadUser,
  LOAD_USER,
  logOut,
  LOG_OUT,
  setActiveRoleSuccess,
  SET_ACTIVE_ROLE_SUCCESS,
  setActiveRoleFailure,
  SET_ACTIVE_ROLE_FAILURE,
} from './actions';

describe('auth actions', () => {
  it('setActiveRole returns the expected action', () => {
    const expectedAction = {
      type: SET_ACTIVE_ROLE,
      payload: 'myRole',
    };

    expect(setActiveRole('myRole')).toEqual(expectedAction);
  });

  it('setActiveRoleSuccess returns the expected action', () => {
    const expectedAction = {
      type: SET_ACTIVE_ROLE_SUCCESS,
      payload: 'myRole',
    };

    expect(setActiveRoleSuccess('myRole')).toEqual(expectedAction);
  });

  it('setActiveRoleFailure returns the expected action', () => {
    const expectedAction = {
      type: SET_ACTIVE_ROLE_FAILURE,
      error: 'error',
    };

    expect(setActiveRoleFailure('error')).toEqual(expectedAction);
  });

  it('loadUser returns the expected action', () => {
    const expectedAction = {
      type: LOAD_USER,
    };

    expect(loadUser()).toEqual(expectedAction);
  });

  it('logOut returns the expected action', () => {
    const expectedAction = {
      type: LOG_OUT,
    };

    expect(logOut()).toEqual(expectedAction);
  });
});
