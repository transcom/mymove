import { setActiveRole, SET_ACTIVE_ROLE, loadUser, LOAD_USER, logOut, LOG_OUT } from './actions';

describe('auth actions', () => {
  it('setActiveRole returns the expected action', () => {
    const expectedAction = {
      type: SET_ACTIVE_ROLE,
      payload: 'myRole',
    };

    expect(setActiveRole('myRole')).toEqual(expectedAction);
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
