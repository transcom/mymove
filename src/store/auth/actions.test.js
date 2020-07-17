import { setActiveRole, SET_ACTIVE_ROLE, loadUser, LOAD_USER } from './actions';

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
});
