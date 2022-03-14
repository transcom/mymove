import { takeLatest, put, call } from 'redux-saga/effects';

import watchFetchUser, { fetchUser } from './auth';

import { setFlashMessage } from 'store/flash/actions';
import { GetIsLoggedIn, GetLoggedInUser } from 'utils/api';
import { LOAD_USER, getLoggedInUserStart, getLoggedInUserSuccess, getLoggedInUserFailure } from 'store/auth/actions';
import { addEntities } from 'shared/Entities/actions';

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
    const testUser = {
      id: 'testUserId',
      email: 'test@example.com',
      first_name: 'Tester',
      roles: [{ id: 'testRole', roleType: 'customer' }],
      service_member: {
        id: 'testServiceMemberId',
        orders: [{ id: 'testorder1' }, { id: 'testorder2' }],
      },
    };

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

    it('stores the user data in the entities reducer', () => {
      const normalizedUser = {
        orders: {
          testorder1: { id: 'testorder1' },
          testorder2: { id: 'testorder2' },
        },
        serviceMembers: {
          testServiceMemberId: {
            id: 'testServiceMemberId',
            orders: ['testorder1', 'testorder2'],
          },
        },
        user: {
          testUserId: {
            id: 'testUserId',
            email: 'test@example.com',
            first_name: 'Tester',
            roles: [{ id: 'testRole', roleType: 'customer' }],
            service_member: 'testServiceMemberId',
          },
        },
      };

      expect(generator.next(testUser).value).toEqual(put(addEntities(normalizedUser)));
    });

    it('stores the user auth data in the auth reducer', () => {
      expect(generator.next().value).toEqual(put(getLoggedInUserSuccess(testUser)));
    });

    it('is done', () => {
      expect(generator.next().done).toEqual(true);
    });
  });
});
