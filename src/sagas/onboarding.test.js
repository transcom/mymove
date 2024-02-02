import { takeLatest, put, call, all } from 'redux-saga/effects';

import {
  watchInitializeOnboarding,
  watchFetchCustomerData,
  fetchCustomerData,
  createServiceMember,
  initializeOnboarding,
} from './onboarding';

import {
  INIT_ONBOARDING,
  FETCH_CUSTOMER_DATA,
  initOnboardingFailed,
  initOnboardingComplete,
} from 'store/onboarding/actions';
import { setFlashMessage } from 'store/flash/actions';
import {
  getLoggedInUser,
  createServiceMember as createServiceMemberApi,
  getMTOShipmentsForMove,
  getAllMoves,
} from 'services/internalApi';
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
      serviceMembers: {
        serviceMemberId: {
          id: 'serviceMemberId',
        },
      },
    };
    const mockMultipleMoveResponseData = {
      currentMove: [],
      previousMoves: [],
    };

    it('makes an API call to request the logged in user', () => {
      expect(generator.next().value).toEqual(call(getLoggedInUser));
    });

    it('stores the user data in entities', () => {
      expect(generator.next(mockResponseData).value).toEqual(put(addEntities(mockResponseData)));
    });

    it('makes an API call to request multiple moves', () => {
      expect(generator.next().value).toEqual(call(getAllMoves, 'serviceMemberId'));
    });

    it('stores the multiple move data in entities', () => {
      expect(generator.next(mockMultipleMoveResponseData).value).toEqual(
        put(
          addEntities({
            serviceMemberMoves: mockMultipleMoveResponseData,
          }),
        ),
      );
    });

    it('yields the user data to the caller', () => {
      expect(generator.next().value).toEqual(mockResponseData);
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
      serviceMembers: {
        serviceMemberId: {
          id: 'serviceMemberId',
        },
      },
      moves: {
        testMoveId: {
          id: 'testMoveId',
        },
      },
      serviceMemberMoves: {
        currentMove: [],
        previousMoves: [],
      },
    };

    const mockMTOResponseData = {
      mtoShipments: {
        testMTOShipmentId: {
          id: 'testMTOShipmentId',
        },
      },
    };

    const mockMultipleMoveResponseData = {
      currentMove: [],
      previousMoves: [],
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

    it('makes an API call to request multiple moves', () => {
      expect(generator.next().value).toEqual(call(getAllMoves, 'serviceMemberId'));
    });

    it('stores the multiple move data in entities', () => {
      expect(generator.next(mockMultipleMoveResponseData).value).toEqual(
        put(
          addEntities({
            serviceMemberMoves: mockMultipleMoveResponseData,
          }),
        ),
      );
    });

    it('yields the user data to the caller', () => {
      expect(generator.next().value).toEqual(mockResponseData);
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

  describe('if the user is logged in and does not have a serviceMember', () => {
    const generator = initializeOnboarding();

    const mockResponseData = {
      user: {
        testUserId: {
          id: 'testUserId',
          email: 'testuser@example.com',
        },
      },
    };

    it('calls the fetchCustomerData saga', () => {
      expect(generator.next().value).toEqual(call(fetchCustomerData));
    });

    it('calls the createServiceMember saga', () => {
      expect(generator.next(mockResponseData).value).toEqual(call(createServiceMember));
    });

    it('puts action initOnboardingComplete', () => {
      expect(generator.next().value).toEqual(put(initOnboardingComplete()));
    });

    it('starts the watch saga', () => {
      expect(generator.next().value).toEqual(all([call(watchFetchCustomerData)]));
    });

    it('is done', () => {
      expect(generator.next().done).toEqual(true);
    });
  });

  describe('if the user is logged in and has a serviceMember', () => {
    const generator = initializeOnboarding();

    const mockResponseData = {
      user: {
        testUserId: {
          id: 'testUserId',
          email: 'testuser@example.com',
          service_member: 'testServiceMemberId',
        },
      },
      serviceMembers: {
        testServiceMemberId: {
          id: 'testServiceMemberId',
        },
      },
    };

    it('calls the fetchCustomerData saga', () => {
      expect(generator.next().value).toEqual(call(fetchCustomerData));
    });

    it('puts action initOnboardingComplete', () => {
      expect(generator.next(mockResponseData).value).toEqual(put(initOnboardingComplete()));
    });

    it('starts the watch saga', () => {
      expect(generator.next().value).toEqual(all([call(watchFetchCustomerData)]));
    });

    it('is done', () => {
      expect(generator.next().done).toEqual(true);
    });
  });
});

describe('createServiceMember saga', () => {
  describe('successful', () => {
    const generator = createServiceMember();

    it('makes API call to createServiceMember', () => {
      expect(generator.next().value).toEqual(call(createServiceMemberApi));
    });

    it('refetches user data', () => {
      expect(generator.next().value).toEqual(call(fetchCustomerData));
    });
  });

  describe('failure', () => {
    const generator = createServiceMember();

    it('makes API call to createServiceMember', () => {
      expect(generator.next().value).toEqual(call(createServiceMemberApi));
    });

    it('sets the error flash message', () => {
      const error = new Error('Service member already exists');
      expect(generator.throw(error).value).toEqual(
        put(
          setFlashMessage(
            'SERVICE_MEMBER_CREATE_ERROR',
            'error',
            'There was an error creating your profile information.',
            'An error occurred',
          ),
        ),
      );
    });

    it('is done', () => {
      expect(generator.next().done).toEqual(true);
    });
  });
});
