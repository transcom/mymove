import { takeLatest, put, call, all, select } from 'redux-saga/effects';
import { cloneableGenerator } from '@redux-saga/testing-utils';
import { push } from 'connected-react-router';

import {
  watchInitializeOnboarding,
  watchFetchCustomerData,
  fetchCustomerData,
  createServiceMember,
  initializeOnboarding,
  watchUpdateServiceMember,
  updateServiceMember,
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
} from 'services/internalApi';
import { addEntities } from 'shared/Entities/actions';
import { CREATE_SERVICE_MEMBER } from 'scenes/ServiceMembers/ducks';
import sampleLoggedInUserPayload from 'shared/User/sampleLoggedInUserPayload';
import { normalizeResponse } from 'services/swaggerRequest';
import { selectServiceMemberFromLoggedInUser, selectServiceMemberProfileState } from 'store/entities/selectors';
import { profileStates } from 'constants/customerStates';

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

    it('selects the service member from the store', () => {
      expect(generator.next().value).toEqual(select(selectServiceMemberFromLoggedInUser));
    });

    it('selects the service member’s profile state from the store', () => {
      expect(
        generator.next({
          id: 'testServiceMemberId',
        }).value,
      ).toEqual(select(selectServiceMemberProfileState));
    });

    it('redirects the user to the first step of the profile wizard', () => {
      const nextPath = `/service-member/testServiceMemberId/conus-status`;
      expect(generator.next(profileStates.EMPTY_PROFILE).value).toEqual(put(push(nextPath)));
    });

    it('puts action initOnboardingComplete', () => {
      expect(generator.next().value).toEqual(put(initOnboardingComplete()));
    });

    it('starts the watch saga', () => {
      expect(generator.next().value).toEqual(all([call(watchFetchCustomerData), call(watchUpdateServiceMember)]));
    });

    it('is done', () => {
      expect(generator.next().done).toEqual(true);
    });
  });

  describe('if the user is logged in and has a serviceMember', () => {
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

    const generator = cloneableGenerator(initializeOnboarding)();
    generator.next(); // fetchCustomerData
    generator.next(mockResponseData); // selectServiceMember

    describe('with no data', () => {
      const clone = generator.clone();

      it('selects the service member’s profile state from the store', () => {
        expect(clone.next(mockResponseData.serviceMembers.testServiceMemberId).value).toEqual(
          select(selectServiceMemberProfileState),
        );
      });

      it('redirects the user to the first step of the profile wizard', () => {
        const nextPath = `/service-member/testServiceMemberId/conus-status`;
        expect(clone.next(profileStates.EMPTY_PROFILE).value).toEqual(put(push(nextPath)));
      });

      it('puts action initOnboardingComplete', () => {
        expect(clone.next().value).toEqual(put(initOnboardingComplete()));
      });

      it('starts the watch saga', () => {
        expect(clone.next().value).toEqual(all([call(watchFetchCustomerData), call(watchUpdateServiceMember)]));
      });

      it('is done', () => {
        expect(clone.next().done).toEqual(true);
      });
    });

    describe('with DOD info complete', () => {
      const clone = generator.clone();

      const serviceMemberData = {
        ...mockResponseData.serviceMembers.testServiceMemberId,
        rank: 'test rank',
        edipi: '1234567890',
        affiliation: 'ARMY',
      };

      it('selects the service member’s profile state from the store', () => {
        expect(clone.next(serviceMemberData).value).toEqual(select(selectServiceMemberProfileState));
      });

      it('redirects the user to the name step of the profile wizard', () => {
        const nextPath = `/service-member/testServiceMemberId/name`;
        expect(clone.next(profileStates.DOD_INFO_COMPLETE).value).toEqual(put(push(nextPath)));
      });

      it('puts action initOnboardingComplete', () => {
        expect(clone.next().value).toEqual(put(initOnboardingComplete()));
      });

      it('starts the watch saga', () => {
        expect(clone.next().value).toEqual(all([call(watchFetchCustomerData), call(watchUpdateServiceMember)]));
      });

      it('is done', () => {
        expect(clone.next().done).toEqual(true);
      });
    });

    describe('with name complete', () => {
      const clone = generator.clone();

      const serviceMemberData = {
        ...mockResponseData.serviceMembers.testServiceMemberId,
        rank: 'test rank',
        edipi: '1234567890',
        affiliation: 'ARMY',
        first_name: 'Tester',
        last_name: 'Testperson',
      };

      it('selects the service member’s profile state from the store', () => {
        expect(clone.next(serviceMemberData).value).toEqual(select(selectServiceMemberProfileState));
      });

      it('redirects the user to the contact info step of the profile wizard', () => {
        const nextPath = `/service-member/testServiceMemberId/contact-info`;
        expect(clone.next(profileStates.NAME_COMPLETE).value).toEqual(put(push(nextPath)));
      });

      it('puts action initOnboardingComplete', () => {
        expect(clone.next().value).toEqual(put(initOnboardingComplete()));
      });

      it('starts the watch saga', () => {
        expect(clone.next().value).toEqual(all([call(watchFetchCustomerData), call(watchUpdateServiceMember)]));
      });

      it('is done', () => {
        expect(clone.next().done).toEqual(true);
      });
    });

    describe('with contact info complete', () => {
      const clone = generator.clone();

      const serviceMemberData = {
        ...mockResponseData.serviceMembers.testServiceMemberId,
        rank: 'test rank',
        edipi: '1234567890',
        affiliation: 'ARMY',
        first_name: 'Tester',
        last_name: 'Testperson',
        telephone: '1234567890',
        personal_email: 'test@example.com',
        email_is_preferred: true,
      };

      it('selects the service member’s profile state from the store', () => {
        expect(clone.next(serviceMemberData).value).toEqual(select(selectServiceMemberProfileState));
      });

      it('redirects the user to the duty station step of the profile wizard', () => {
        const nextPath = `/service-member/testServiceMemberId/duty-station`;
        expect(clone.next(profileStates.CONTACT_INFO_COMPLETE).value).toEqual(put(push(nextPath)));
      });

      it('puts action initOnboardingComplete', () => {
        expect(clone.next().value).toEqual(put(initOnboardingComplete()));
      });

      it('starts the watch saga', () => {
        expect(clone.next().value).toEqual(all([call(watchFetchCustomerData), call(watchUpdateServiceMember)]));
      });

      it('is done', () => {
        expect(clone.next().done).toEqual(true);
      });
    });

    describe('with duty station info complete', () => {
      const clone = generator.clone();

      const serviceMemberData = {
        ...mockResponseData.serviceMembers.testServiceMemberId,
        rank: 'test rank',
        edipi: '1234567890',
        affiliation: 'ARMY',
        first_name: 'Tester',
        last_name: 'Testperson',
        telephone: '1234567890',
        personal_email: 'test@example.com',
        email_is_preferred: true,
        current_station: {
          id: 'testDutyStationId',
        },
      };

      it('selects the service member’s profile state from the store', () => {
        expect(clone.next(serviceMemberData).value).toEqual(select(selectServiceMemberProfileState));
      });

      it('redirects the user to the address step of the profile wizard', () => {
        const nextPath = `/service-member/testServiceMemberId/residence-address`;
        expect(clone.next(profileStates.DUTY_STATION_COMPLETE).value).toEqual(put(push(nextPath)));
      });

      it('puts action initOnboardingComplete', () => {
        expect(clone.next().value).toEqual(put(initOnboardingComplete()));
      });

      it('starts the watch saga', () => {
        expect(clone.next().value).toEqual(all([call(watchFetchCustomerData), call(watchUpdateServiceMember)]));
      });

      it('is done', () => {
        expect(clone.next().done).toEqual(true);
      });
    });

    describe('with address info complete', () => {
      const clone = generator.clone();

      const serviceMemberData = {
        ...mockResponseData.serviceMembers.testServiceMemberId,
        rank: 'test rank',
        edipi: '1234567890',
        affiliation: 'ARMY',
        first_name: 'Tester',
        last_name: 'Testperson',
        telephone: '1234567890',
        personal_email: 'test@example.com',
        email_is_preferred: true,
        current_station: {
          id: 'testDutyStationId',
        },
        residential_address: {
          street: '123 Main St',
        },
      };

      it('selects the service member’s profile state from the store', () => {
        expect(clone.next(serviceMemberData).value).toEqual(select(selectServiceMemberProfileState));
      });

      it('redirects the user to the backup address step of the profile wizard', () => {
        const nextPath = `/service-member/testServiceMemberId/backup-mailing-address`;
        expect(clone.next(profileStates.ADDRESS_COMPLETE).value).toEqual(put(push(nextPath)));
      });

      it('puts action initOnboardingComplete', () => {
        expect(clone.next().value).toEqual(put(initOnboardingComplete()));
      });

      it('starts the watch saga', () => {
        expect(clone.next().value).toEqual(all([call(watchFetchCustomerData), call(watchUpdateServiceMember)]));
      });

      it('is done', () => {
        expect(clone.next().done).toEqual(true);
      });
    });

    describe('with backup address info complete', () => {
      const clone = generator.clone();

      const serviceMemberData = {
        ...mockResponseData.serviceMembers.testServiceMemberId,
        rank: 'test rank',
        edipi: '1234567890',
        affiliation: 'ARMY',
        first_name: 'Tester',
        last_name: 'Testperson',
        telephone: '1234567890',
        personal_email: 'test@example.com',
        email_is_preferred: true,
        current_station: {
          id: 'testDutyStationId',
        },
        residential_address: {
          street: '123 Main St',
        },
        backup_mailing_address: {
          street: '456 Main St',
        },
      };

      it('selects the service member’s profile state from the store', () => {
        expect(clone.next(serviceMemberData).value).toEqual(select(selectServiceMemberProfileState));
      });

      it('redirects the user to the backup contacts step of the profile wizard', () => {
        const nextPath = `/service-member/testServiceMemberId/backup-contacts`;
        expect(clone.next(profileStates.BACKUP_ADDRESS_COMPLETE).value).toEqual(put(push(nextPath)));
      });

      it('puts action initOnboardingComplete', () => {
        expect(clone.next().value).toEqual(put(initOnboardingComplete()));
      });

      it('starts the watch saga', () => {
        expect(clone.next().value).toEqual(all([call(watchFetchCustomerData), call(watchUpdateServiceMember)]));
      });

      it('is done', () => {
        expect(clone.next().done).toEqual(true);
      });
    });

    describe('with all profile info complete', () => {
      const clone = generator.clone();

      const serviceMemberData = {
        ...mockResponseData.serviceMembers.testServiceMemberId,
        rank: 'test rank',
        edipi: '1234567890',
        affiliation: 'ARMY',
        first_name: 'Tester',
        last_name: 'Testperson',
        telephone: '1234567890',
        personal_email: 'test@example.com',
        email_is_preferred: true,
        current_station: {
          id: 'testDutyStationId',
        },
        residential_address: {
          street: '123 Main St',
        },
        backup_mailing_address: {
          street: '456 Main St',
        },
        backup_contacts: [
          {
            id: 'testBackupContact',
          },
        ],
      };

      it('selects the service member’s profile state from the store', () => {
        expect(clone.next(serviceMemberData).value).toEqual(select(selectServiceMemberProfileState));
      });

      it('redirects the user to the Home page', () => {
        expect(clone.next(profileStates.BACKUP_CONTACTS_COMPLETE).value).toEqual(put(push('/')));
      });

      it('puts action initOnboardingComplete', () => {
        expect(clone.next().value).toEqual(put(initOnboardingComplete()));
      });

      it('starts the watch saga', () => {
        expect(clone.next().value).toEqual(all([call(watchFetchCustomerData), call(watchUpdateServiceMember)]));
      });

      it('is done', () => {
        expect(clone.next().done).toEqual(true);
      });
    });
  });
});

describe('createServiceMember saga', () => {
  describe('successful', () => {
    const generator = createServiceMember();
    const mockServiceMember = {
      id: 'testServiceMemberId',
      user_id: 'testUserId',
      is_profile_complete: false,
    };

    it('puts the CREATE_SERVICE_MEMBER.start action', () => {
      expect(generator.next().value).toEqual(
        put({
          type: CREATE_SERVICE_MEMBER.start,
        }),
      );
    });

    it('makes API call to createServiceMember', () => {
      expect(generator.next().value).toEqual(call(createServiceMemberApi));
    });

    it('puts the CREATE_SERVICE_MEMBER.success action', () => {
      expect(generator.next(mockServiceMember).value).toEqual(
        put({
          type: CREATE_SERVICE_MEMBER.success,
          payload: mockServiceMember,
        }),
      );
    });

    it('refetches user data', () => {
      expect(generator.next().value).toEqual(call(fetchCustomerData));
    });
  });

  describe('failure', () => {
    const generator = createServiceMember();

    it('puts the CREATE_SERVICE_MEMBER.start action', () => {
      expect(generator.next().value).toEqual(
        put({
          type: CREATE_SERVICE_MEMBER.start,
        }),
      );
    });

    it('makes API call to createServiceMember', () => {
      expect(generator.next().value).toEqual(call(createServiceMemberApi));
    });

    it('puts the CREATE_SERVICE_MEMBER.failure action', () => {
      const error = new Error('Service member already exists');
      expect(generator.throw(error).value).toEqual(
        put({
          type: CREATE_SERVICE_MEMBER.failure,
          error,
        }),
      );
    });

    it('sets the error flash message', () => {
      expect(generator.next().value).toEqual(
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

describe('watchUpdateServiceMember', () => {
  const generator = watchUpdateServiceMember();

  it('takes a UPDATE_SERVICE_MEMBER_SUCCESS action and calls updateServiceMember', () => {
    expect(generator.next().value).toEqual(takeLatest('UPDATE_SERVICE_MEMBER_SUCCESS', updateServiceMember));
  });
});

describe('updateServiceMember', () => {
  const action = {
    type: 'UPDATE_SERVICE_MEMBER_SUCCESS',
    payload: sampleLoggedInUserPayload.payload.service_member,
  };

  const generator = updateServiceMember(action);

  it('normalizes the data and puts it in entities', () => {
    const normalizedData = normalizeResponse(action.payload, 'serviceMember');
    expect(generator.next().value).toEqual(put(addEntities(normalizedData)));
  });
});
