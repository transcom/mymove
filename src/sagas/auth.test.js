import { takeLatest, put, call } from 'redux-saga/effects';

import watchFetchUser, { fetchUser } from './auth';

import { LOAD_USER } from 'store/auth/actions';
import { GetIsLoggedIn, GetLoggedInUser } from 'shared/User/api';
import { getLoggedInActions } from 'shared/Data/users';
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
  describe('if the user is not logged in', () => {
    const generator = fetchUser();

    it('dispatches the GET_LOGGED_IN_USER_START action', () => {
      expect(generator.next().value).toEqual(put(getLoggedInActions.start()));
    });

    it('makes the GetIsLoggedIn API call', () => {
      expect(generator.next().value).toEqual(call(GetIsLoggedIn));
    });

    it('dispatches the User is not logged in error action', () => {
      expect(generator.next(false).value).toEqual(put(getLoggedInActions.error('User is not logged in')));
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
      expect(generator.next().value).toEqual(put(getLoggedInActions.start()));
    });

    it('makes the GetIsLoggedIn API call', () => {
      expect(generator.next().value).toEqual(call(GetIsLoggedIn));
    });

    it('makes the GetLoggedInUser API call', () => {
      expect(generator.next(true).value).toEqual(call(GetLoggedInUser));
    });

    it('stores the user data in the entities reducer', () => {
      const normalizedUser = {
        orders: {
          testorder1: { id: 'testorder1' },
          testorder2: { id: 'testorder2' },
        },
        roles: {
          testRole: { id: 'testRole', roleType: 'customer' },
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
            roles: ['testRole'],
            service_member: 'testServiceMemberId',
          },
        },
      };

      expect(generator.next(testUser).value).toEqual(put(addEntities(normalizedUser)));
    });

    it('stores the user data in the user reducer (legacy)', () => {
      // TODO - delete when deprecating the user reducer
      expect(generator.next().value).toEqual(put(getLoggedInActions.success(testUser)));
    });

    it('is done', () => {
      expect(generator.next().done).toEqual(true);
    });
  });
});
