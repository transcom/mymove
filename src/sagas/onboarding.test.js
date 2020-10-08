import { takeLatest, put, call } from 'redux-saga/effects';

import {
  watchInitializeOnboarding,
  watchFetchCustomerData,
  fetchCustomerData,
  initializeOnboarding,
} from './onboarding';

import {
  INIT_ONBOARDING,
  FETCH_CUSTOMER_DATA,
  initOnboardingFailed,
  initOnboardingComplete,
} from 'store/onboarding/actions';
import { getLoggedInUser, getMTOShipmentsForMove } from 'services/internalApi';
import { addEntities } from 'shared/Entities/actions';

describe('watchInitializeOnboarding', () => {
  const generator = watchInitializeOnboarding();

  it('takes the latest INIT_ONBOARDING action and calls initializeOnboarding', () => {
    expect(generator.next().value).toEqual(takeLatest(INIT_ONBOARDING, initializeOnboarding));
  });

  it('is done', () => {
    expect(generator.next().done).toEqual(true);
  });
});

describe('watchFetchCustomerData', () => {
  const generator = watchFetchCustomerData();

  it('takes a FETCH_CUSTOMER_DATA action and calls fetchCustomerData', () => {
    expect(generator.next().value).toEqual(takeLatest(FETCH_CUSTOMER_DATA, fetchCustomerData));
  });
});

describe('fetchCustomerData', () => {
  describe('if the user doesn’t have a move', () => {
    const generator = fetchCustomerData();

    const mockResponseData = {
      user: {
        testUserId: {
          id: 'testUserId',
          email: 'testuser@example.com',
        },
      },
    };

    it('makes an API call to request the logged in user', () => {
      expect(generator.next().value).toEqual(call(getLoggedInUser));
    });

    it('stores the user data in entities', () => {
      expect(generator.next(mockResponseData).value).toEqual(put(addEntities(mockResponseData)));
    });

    it('is done', () => {
      expect(generator.next().done).toEqual(true);
    });
  });

  describe('if the user has a move', () => {
    const generator = fetchCustomerData();

    const mockResponseData = {
      user: {
        testUserId: {
          id: 'testUserId',
          email: 'testuser@example.com',
        },
      },
      moves: {
        testMoveId: {
          id: 'testMoveId',
        },
      },
    };

    const mockMTOResponseData = {
      mtoShipments: {
        testMTOShipmentId: {
          id: 'testMTOShipmentId',
        },
      },
    };

    it('makes an API call to request the logged in user', () => {
      expect(generator.next().value).toEqual(call(getLoggedInUser));
    });

    it('stores the user data in entities', () => {
      expect(generator.next(mockResponseData).value).toEqual(put(addEntities(mockResponseData)));
    });

    it('makes an API call to request the MTO shipments', () => {
      expect(generator.next().value).toEqual(call(getMTOShipmentsForMove, 'testMoveId'));
    });

    it('stores the MTO shipment data in entities', () => {
      expect(generator.next(mockMTOResponseData).value).toEqual(put(addEntities(mockMTOResponseData)));
    });

    it('is done', () => {
      expect(generator.next().done).toEqual(true);
    });
  });
});

describe('initializeOnboarding', () => {
  describe('if the user is not logged in', () => {
    const generator = initializeOnboarding();

    it('calls the fetchCustomerData saga', () => {
      expect(generator.next().value).toEqual(call(fetchCustomerData));
    });

    it('puts action initOnboardingFailed with the error', () => {
      const error = new Error('User not logged in');
      expect(generator.throw(error).value).toEqual(put(initOnboardingFailed(error)));
    });

    it('is done', () => {
      expect(generator.next().done).toEqual(true);
    });
  });

  describe('if the user is logged in', () => {
    const generator = initializeOnboarding();

    it('calls the fetchCustomerData saga', () => {
      expect(generator.next().value).toEqual(call(fetchCustomerData));
    });

    it('puts action initOnboardingComplete', () => {
      expect(generator.next().value).toEqual(put(initOnboardingComplete()));
    });

    it('starts the watchFetchCustomerData saga', () => {
      expect(generator.next().value).toEqual(call(watchFetchCustomerData));
    });

    it('is done', () => {
      expect(generator.next().done).toEqual(true);
    });
  });
});
