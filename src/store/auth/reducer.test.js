import authReducer, { initialState } from './reducer';
import { setActiveRole, logOut, setActiveRoleSuccess, setActiveRoleFailure } from './actions';
import { selectIsLoggedIn, selectUnderMaintenance } from './selectors';

import { roleTypes } from 'constants/userRoles';

const primaryOffice = {
  address: null,
  created_at: '2025-06-05T15:23:29.086Z',
  gbloc: 'KKFA',
  id: '171b54fa-4c89-45d8-8111-a2d65818ff8c',
  name: 'JPPSO - North Central (KKFA) - USAF',
  phone_lines: [],
  updated_at: '2025-06-05T15:23:29.086Z',
};
const secondaryOffice = {
  address: null,
  created_at: '2025-06-05T15:23:29.086Z',
  gbloc: 'AGFM',
  id: '3132b512-1889-4776-a666-9c08a24afe20',
  name: 'JPPSO - North East (AGFM) - USAF',
  phone_lines: [],
  updated_at: '2025-06-05T15:23:29.086Z',
};

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
      isSettingActiveRole: true,
    });
  });

  it('handles the setActiveRoleSucces action', () => {
    expect(authReducer(initialState, setActiveRoleSuccess('myRole'))).toEqual({
      ...initialState,
      activeRole: 'myRole',
      isSettingActiveRole: false,
    });
  });

  it('handles the setActiveRoleFailure action', () => {
    expect(authReducer(initialState, setActiveRoleFailure())).toEqual({
      ...initialState,
      isSettingActiveRole: false,
    });
  });

  it('handles the logOut action', () => {
    const currentState = {
      ...initialState,
      activeRole: 'myRole',
    };

    expect(authReducer(currentState, logOut())).toEqual(initialState);
  });

  it('handles the GET_LOGGED_IN_USER_SUCCESS action with no activeRole or activeOffice set', () => {
    const action = {
      type: 'GET_LOGGED_IN_USER_SUCCESS',
      payload: {
        activeRole: {
          roleType: roleTypes.TOO,
        },
        office_user: {
          transportation_office_assignments: [
            { primaryOffice: true, transportationOffice: primaryOffice },
            { primaryOffice: false, transportationOffice: secondaryOffice },
          ],
        },
      },
    };

    expect(authReducer(initialState, action)).toEqual({
      ...initialState,
      activeRole: roleTypes.TOO,
      activeOffice: primaryOffice,
      hasSucceeded: true,
      hasErrored: false,
      isLoading: false,
      isLoggedIn: true,
    });
  });

  it('handles the GET_LOGGED_IN_USER_SUCCESS action with an activeRole and activeOffice already set', () => {
    const currentState = {
      ...initialState,
      activeRole: roleTypes.TOO,
      activeOffice: primaryOffice,
      hasSucceeded: true,
      hasErrored: false,
      isLoading: false,
      isLoggedIn: true,
    };

    const action = {
      type: 'GET_LOGGED_IN_USER_SUCCESS',
      payload: {
        activeRole: {
          roleType: roleTypes.TOO,
        },
        office_user: {
          transportation_office_assignments: [
            { primaryOffice: true, transportationOffice: primaryOffice },
            { primaryOffice: false, transportationOffice: secondaryOffice },
          ],
        },
      },
    };

    expect(authReducer(currentState, action)).toEqual(currentState);
  });

  it('handles the GET_LOGGED_IN_USER_SUCCESS action when a non-primary office is already in state', () => {
    const currentState = {
      ...initialState,
      activeRole: roleTypes.TOO,
      activeOffice: secondaryOffice,
      hasSucceeded: true,
      hasErrored: false,
      isLoading: false,
      isLoggedIn: true,
    };

    const action = {
      type: 'GET_LOGGED_IN_USER_SUCCESS',
      payload: {
        activeRole: {
          roleType: roleTypes.TOO,
        },
        office_user: {
          transportation_office_assignments: [
            { primaryOffice: true, transportationOffice: primaryOffice },
            { primaryOffice: false, transportationOffice: secondaryOffice },
          ],
        },
      },
    };

    expect(authReducer(currentState, action)).toEqual(currentState);
  });

  it('SET_ACTIVE_OFFICE sets activeOffice to the passed in office', () => {
    const currentState = {
      ...initialState,
      activeOffice: primaryOffice,
    };

    const action = {
      type: 'SET_ACTIVE_OFFICE',
      payload: secondaryOffice,
    };

    expect(authReducer(currentState, action).activeOffice).toEqual(secondaryOffice);
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

describe('setUnderMaintenance', () => {
  it('returns boolean as to whether or not app is under maintenance', () => {
    const testState = {
      auth: { underMaintenance: true },
    };

    expect(selectUnderMaintenance(testState)).toEqual(testState.auth.underMaintenance);
  });
});
