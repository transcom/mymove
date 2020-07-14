import { takeLatest, put, call } from 'redux-saga/effects';

import watchFetchUser, { fetchUser, getLoggedInActions } from './auth';

import { GetIsLoggedIn, GetLoggedInUser } from 'shared/User/api';

describe('watchFetchUser saga', () => {
  const generator = watchFetchUser();

  it('takes the latest GET_LOGGED_IN_USER_START action and calls fetchUser', () => {
    expect(generator.next().value).toEqual(takeLatest('GET_LOGGED_IN_USER_START', fetchUser));
  });

  it('is done', () => {
    expect(generator.next().done).toEqual(true);
  });
});

describe('fetchUser saga', () => {
  describe('if the user is logged in and is not a service member', () => {
    const generator = fetchUser();

    it('makes the GetIsLoggedIn API call', () => {
      expect(generator.next().value).toEqual(call(GetIsLoggedIn));
    });

    it('makes the GetLoggedInUser API call', () => {
      expect(generator.next(true).value).toEqual(call(GetLoggedInUser));
    });

    it('stores the user data in Redux', () => {
      const testUser = { id: 'testUserId' };
      expect(generator.next(testUser).value).toEqual(put(getLoggedInActions.success(testUser)));
    });

    it('is done', () => {
      expect(generator.next().done).toEqual(true);
    });
  });

  describe('if the user is logged in and is a service member', () => {
    const generator = fetchUser();

    it('makes the GetIsLoggedIn API call', () => {
      expect(generator.next().value).toEqual(call(GetIsLoggedIn));
    });

    it('makes the GetLoggedInUser API call', () => {
      expect(generator.next(true).value).toEqual(call(GetLoggedInUser));
    });

    it('stores the user data in Redux', () => {
      const testUser = { id: 'testUserId', service_member: true };
      expect(generator.next(testUser).value).toEqual(put(getLoggedInActions.success(testUser)));
    });

    it('is done', () => {
      expect(generator.next().done).toEqual(true);
    });
  });
});
