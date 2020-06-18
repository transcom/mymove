import { setActiveRole, SET_ACTIVE_ROLE } from './actions';

describe('auth actions', () => {
  it('setActiveRole returns the expected action', () => {
    const expectedAction = {
      type: SET_ACTIVE_ROLE,
      payload: 'myRole',
    };

    expect(setActiveRole('myRole')).toEqual(expectedAction);
  });
});
