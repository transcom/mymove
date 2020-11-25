import { selectLoggedInUser, selectServiceMemberFromLoggedInUser, selectIsProfileComplete } from './selectors';

describe('selectLoggedInUser', () => {
  it('returns the first user stored in entities', () => {
    const testState = {
      entities: {
        user: {
          userId123: {
            id: 'userId123',
          },
        },
      },
    };

    expect(selectLoggedInUser(testState)).toEqual(testState.entities.user.userId123);
  });

  it('returns null if there is no user in entities', () => {
    const testState = {
      entities: {},
    };

    expect(selectLoggedInUser(testState)).toEqual(null);
  });
});

describe('selectServiceMemberFromLoggedInUser', () => {
  it('returns the service member associated with the logged in user', () => {
    const testState = {
      entities: {
        user: {
          userId123: {
            id: 'userId123',
            service_member: 'serviceMemberId456',
          },
        },
        serviceMembers: {
          serviceMemberId456: {
            id: 'serviceMemberId456',
          },
        },
      },
    };

    expect(selectServiceMemberFromLoggedInUser(testState)).toEqual(
      testState.entities.serviceMembers.serviceMemberId456,
    );
  });

  it('returns null if there is no user in entities', () => {
    const testState = {
      entities: {},
    };

    expect(selectServiceMemberFromLoggedInUser(testState)).toEqual(null);
  });

  it('returns null if the user has no service member', () => {
    const testState = {
      entities: {
        user: {
          userId123: {
            id: 'userId123',
          },
        },
      },
    };

    expect(selectServiceMemberFromLoggedInUser(testState)).toEqual(null);
  });

  it('returns null if the service member is not in entities', () => {
    const testState = {
      entities: {
        user: {
          userId123: {
            id: 'userId123',
            service_member: 'serviceMemberId456',
          },
        },
        serviceMembers: {},
      },
    };

    expect(selectServiceMemberFromLoggedInUser(testState)).toEqual(null);
  });
});

describe('selectIsProfileComplete', () => {
  const testServiceMember = {
    affiliation: 'ARMY',
    backup_mailing_address: {
      city: 'Washington',
      postal_code: '20021',
      state: 'DC',
      street_address_1: '200 K St',
    },
    created_at: '2018-05-25T15:48:49.918Z',
    current_station: {
      address: {
        city: 'Colorado Springs',
        country: 'United States',
        postal_code: '80913',
        state: 'CO',
        street_address_1: 'n/a',
      },
      affiliation: 'ARMY',
      created_at: '2018-05-20T18:36:45.034Z',
      id: '28f63a9d-8fff-4a0f-84ef-661c5c8c354e',
      name: 'Ft Carson',
      updated_at: '2018-05-20T18:36:45.034Z',
    },
    edipi: '1234567890',
    email_is_preferred: false,
    first_name: 'Erin',
    id: '1694e00e-17ff-43fe-af6d-ab0519a18ff2',
    is_profile_complete: true,
    last_name: 'Stanfill',
    middle_name: '',
    personal_email: 'erin@truss.works',
    phone_is_preferred: true,
    rank: 'O_4_W_4',
    residential_address: {
      city: 'Somewhere',
      postal_code: '80913',
      state: 'CO',
      street_address_1: '123 Main',
    },
    telephone: '555-555-5556',
    updated_at: '2018-05-25T21:39:10.484Z',
    user_id: 'b46e651e-9d1c-4be5-bb88-bba58e817696',
  };

  it('returns false if the required attributes are not complete', () => {
    const testState = {
      entities: {
        user: {
          userId123: {
            id: 'userId123',
            service_member: 'serviceMemberId456',
          },
        },
        serviceMembers: {
          serviceMemberId456: {
            id: 'serviceMemberId456',
            backup_contacts: [],
          },
        },
      },
    };

    expect(selectIsProfileComplete(testState)).toEqual(false);
  });

  it('returns false if there are no backupContacts', () => {
    const testState = {
      entities: {
        user: {
          userId123: {
            id: 'userId123',
            service_member: 'serviceMemberId456',
          },
        },
        serviceMembers: {
          serviceMemberId456: {
            ...testServiceMember,
            backup_contacts: [],
          },
        },
      },
    };

    expect(selectIsProfileComplete(testState)).toEqual(false);
  });

  it('returns true if all required attributes are complete', () => {
    const testState = {
      entities: {
        user: {
          userId123: {
            id: 'userId123',
            service_member: 'serviceMemberId456',
          },
        },
        serviceMembers: {
          serviceMemberId456: {
            ...testServiceMember,
            backup_contacts: ['backupContact1'],
          },
        },
      },
    };

    expect(selectIsProfileComplete(testState)).toEqual(true);
  });
});
