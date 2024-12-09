import { takeLatest, put, call } from 'redux-saga/effects';

import watchFetchUser, { fetchUser } from './auth';

import { setFlashMessage } from 'store/flash/actions';
import { GetIsLoggedIn, GetLoggedInUser, GetOktaUser, GetAdminUser } from 'utils/api';
import { LOAD_USER, getLoggedInUserStart, getLoggedInUserFailure } from 'store/auth/actions';
import { setAdminUser } from 'shared/Entities/actions';
import { serviceName } from 'shared/constants';

jest.mock('shared/constants', () => ({
  ...jest.requireActual('shared/constants'),
  serviceName: jest.fn(),
}));

beforeEach(() => {
  jest.clearAllMocks();
});

beforeEach(() => {
  jest.clearAllMocks();
});

describe('watchFetchUser saga', () => {
  const generator = watchFetchUser();

  it('takes the latest LOAD_USER action and calls fetchUser', () => {
    expect(generator.next().value).toEqual(takeLatest(LOAD_USER, fetchUser));
  });

  it('is done', () => {
    expect(generator.next().done).toEqual(true);
  });
});

describe('fetchUser saga', () => {
  describe('if the get logged in request fails', () => {
    const generator = fetchUser();

    it('dispatches the GET_LOGGED_IN_USER_START action', () => {
      expect(generator.next().value).toEqual(put(getLoggedInUserStart()));
    });

    it('makes the GetIsLoggedIn API call', () => {
      expect(generator.next().value).toEqual(call(GetIsLoggedIn));
    });

    it('sets the flash error', () => {
      const error = new Error('Logged In request failed');
      expect(generator.throw(error).value).toEqual(
        put(
          setFlashMessage(
            'LOGGED_IN_GET_ERROR',
            'error',
            'There was an error loading your user information.',
            'An error occurred',
          ),
        ),
      );
    });

    it('dispatches the User is not logged in error action', () => {
      const error = new Error('Logged In request failed');
      expect(generator.next(false).value).toEqual(put(getLoggedInUserFailure(error)));
    });
  });

  describe('if the user is not logged in', () => {
    const generator = fetchUser();

    it('dispatches the GET_LOGGED_IN_USER_START action', () => {
      expect(generator.next().value).toEqual(put(getLoggedInUserStart()));
    });

    it('makes the GetIsLoggedIn API call', () => {
      expect(generator.next().value).toEqual(call(GetIsLoggedIn));
    });

    it('dispatches the User is not logged in error action', () => {
      expect(generator.next({ isLoggedIn: false }).value).toEqual(put(getLoggedInUserFailure('User is not logged in')));
    });

    it('is done', () => {
      expect(generator.next().done).toEqual(true);
    });
  });

  describe('if the get user data request fails', () => {
    const generator = fetchUser();

    it('dispatches the GET_LOGGED_IN_USER_START action', () => {
      expect(generator.next().value).toEqual(put(getLoggedInUserStart()));
    });

    it('makes the GetIsLoggedIn API call', () => {
      expect(generator.next().value).toEqual(call(GetIsLoggedIn));
    });

    it('makes the GetLoggedInUser API call', () => {
      expect(generator.next({ isLoggedIn: true }).value).toEqual(call(GetLoggedInUser));
    });

    it('sets the flash error', () => {
      const error = new Error('Get user request failed');
      expect(generator.throw(error).value).toEqual(
        put(
          setFlashMessage(
            'USER_GET_ERROR',
            'error',
            'There was an error loading your user information.',
            'An error occurred',
          ),
        ),
      );
    });

    it('dispatches the User is not logged in error action', () => {
      const error = new Error('Get user request failed');
      expect(generator.next().value).toEqual(put(getLoggedInUserFailure(error)));
    });

    it('is done', () => {
      expect(generator.next().done).toEqual(true);
    });
  });

  describe('if the user is logged in', () => {
    let generator;
    beforeEach(() => {
      serviceName.mockResolvedValue('my');
      generator = fetchUser();
    });

    it('gets logged in user and okta data and sets to state', () => {
      expect(generator.next().value).toEqual(put(getLoggedInUserStart()));
      expect(generator.next().value).toEqual(call(GetIsLoggedIn));
      expect(generator.next({ isLoggedIn: true }).value).toEqual(call(GetLoggedInUser));
      expect(generator.next().value).toEqual(call(GetOktaUser));
    });
  });

  describe('if the user is logged in and isAdminSite is true', () => {
    let generator;
    beforeEach(() => {
      serviceName.mockReturnValue('admin');
      generator = fetchUser();
    });

    it('gets logged in user, okta, and admin user data and sets to state', () => {
      expect(generator.next().value).toEqual(put(getLoggedInUserStart()));
      expect(generator.next().value).toEqual(call(GetIsLoggedIn));
      expect(generator.next({ isLoggedIn: true }).value).toEqual(call(GetLoggedInUser));
      expect(generator.next().value).toEqual(call(GetOktaUser));
      expect(generator.next().value).toEqual(call(GetAdminUser));
      const adminUser = { id: 1, name: 'Admin' };
      expect(generator.next(adminUser).value).toEqual(put(setAdminUser(adminUser)));
    });
  });
});

describe('fetch underMaintenance', () => {
  const generator = fetchUser();

  it('makes the GetIsLoggedIn API call', () => {
    expect(generator.next().value).toEqual(put(getLoggedInUserStart()));
    expect(generator.next().value).toEqual(call(GetIsLoggedIn));
    expect(generator.next({ isLoggedIn: true, underMaintenance: false }).value).toEqual(call(GetLoggedInUser));
    // expect(generator.next().value).toEqual(put(setUnderMaintenance));
  });
});
