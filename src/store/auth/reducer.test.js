import authReducer, { initialState } from './reducer';
import { setActiveRole, logOut } from './actions';
import { selectIsLoggedIn } from './selectors';

import { roleTypes } from 'constants/userRoles';

describe('authReducer', () => {
  it('returns the initial state by default', () => {
    expect(authReducer(undefined, undefined)).toEqual(initialState);
  });

  it('returns the existing state if activeRole is set for unhandled actions', () => {
    const currentState = {
      ...initialState,
      activeRole: 'myRole',
    };

    expect(authReducer(currentState, undefined)).toEqual(currentState);
  });

  it('handles the setActiveRole action', () => {
    expect(authReducer(initialState, setActiveRole('myRole'))).toEqual({
      ...initialState,
      activeRole: 'myRole',
    });
  });

  it('handles the logOut action', () => {
    const currentState = {
      ...initialState,
      activeRole: 'myRole',
    };

    expect(authReducer(currentState, logOut())).toEqual(initialState);
  });

  it('handles the GET_LOGGED_IN_USER_SUCCESS action with no activeRole set', () => {
    const action = {
      type: 'GET_LOGGED_IN_USER_SUCCESS',
      payload: {
        roles: [
          {
            roleType: roleTypes.CUSTOMER,
          },
          {
            roleType: roleTypes.TOO,
          },
        ],
      },
    };

    expect(authReducer(initialState, action)).toEqual({
      ...initialState,
      activeRole: roleTypes.TOO,
      hasSucceeded: true,
      hasErrored: false,
      isLoading: false,
      isLoggedIn: true,
    });
  });

  it('handles the GET_LOGGED_IN_USER_SUCCESS action with an activeRole already set', () => {
    const currentState = {
      ...initialState,
      activeRole: roleTypes.TOO,
      hasSucceeded: true,
      hasErrored: false,
      isLoading: false,
      isLoggedIn: true,
    };

    const action = {
      type: 'GET_LOGGED_IN_USER_SUCCESS',
      payload: {
        roles: [
          {
            roleType: roleTypes.CUSTOMER,
          },
          {
            roleType: roleTypes.TOO,
          },
        ],
      },
    };

    expect(authReducer(currentState, action)).toEqual(currentState);
  });
});

describe('selectIsLoggedIn', () => {
  it('returns boolean as to whether user is logged in or not', () => {
    const testState = {
      auth: { isLoggedIn: true },
    };

    expect(selectIsLoggedIn(testState)).toEqual(testState.auth.isLoggedIn);
  });
});
