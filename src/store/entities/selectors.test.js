import {
  selectLoggedInUser,
  selectServiceMemberFromLoggedInUser,
  selectServiceMemberProfileState,
  selectIsProfileComplete,
  selectBackupContacts,
  selectCurrentDutyLocation,
  selectOrdersForLoggedInUser,
  selectCurrentOrders,
  selectMovesForLoggedInUser,
  selectMovesForCurrentOrders,
  selectCurrentMove,
  selectCurrentPPM,
  selectPPMForMove,
  selectPPMSitEstimate,
  selectWeightAllotmentsForLoggedInUser,
  selectWeightTicketAndIndexById,
} from './selectors';

import { profileStates } from 'constants/customerStates';

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

describe('selectServiceMemberProfileState', () => {
  it('returns EMPTY_PROFILE if there is no DoD data', () => {
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

    expect(selectServiceMemberProfileState(testState)).toEqual(profileStates.EMPTY_PROFILE);
  });

  it('returns DOD_INFO_COMPLETE if there is no name data', () => {
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
            affiliation: 'ARMY',
            edipi: '1234567890',
          },
        },
      },
    };

    expect(selectServiceMemberProfileState(testState)).toEqual(profileStates.DOD_INFO_COMPLETE);
  });

  it('returns NAME_COMPLETE if there is no contact info data', () => {
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
            affiliation: 'ARMY',
            edipi: '1234567890',
            first_name: 'Erin',
            last_name: 'Stanfill',
            middle_name: '',
          },
        },
      },
    };

    expect(selectServiceMemberProfileState(testState)).toEqual(profileStates.NAME_COMPLETE);
  });

  it('returns CONTACT_INFO_COMPLETE if there is no residential address data', () => {
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
            affiliation: 'ARMY',
            edipi: '1234567890',
            first_name: 'Erin',
            last_name: 'Stanfill',
            middle_name: '',
            personal_email: 'erin@truss.works',
            phone_is_preferred: true,
            telephone: '555-555-5556',
            email_is_preferred: false,
          },
        },
      },
    };

    expect(selectServiceMemberProfileState(testState)).toEqual(profileStates.CONTACT_INFO_COMPLETE);
  });

  it('returns ADDRESS_COMPLETE if there is no backup address data', () => {
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
            affiliation: 'ARMY',
            edipi: '1234567890',
            first_name: 'Erin',
            last_name: 'Stanfill',
            middle_name: '',
            personal_email: 'erin@truss.works',
            phone_is_preferred: true,
            telephone: '555-555-5556',
            email_is_preferred: false,
            current_location: {
              id: 'testDutyLocationId',
              address: {
                city: 'Colorado Springs',
                country: 'United States',
                postalCode: '80913',
                state: 'CO',
                streetAddress1: 'n/a',
              },
            },
            residential_address: {
              city: 'Somewhere',
              postalCode: '80913',
              state: 'CO',
              streetAddress1: '123 Main',
            },
          },
        },
      },
    };

    expect(selectServiceMemberProfileState(testState)).toEqual(profileStates.ADDRESS_COMPLETE);
  });

  it('returns BACKUP_ADDRESS_COMPLETE if there is no backup contact data', () => {
    const testState = {
      entities: {
        backupContacts: {
          backupContact789: {
            id: 'backupContact789',
            service_member_id: 'serviceMemberId456',
          },
          backupContact8910: {
            id: 'backupContact8910',
            service_member_id: 'serviceMemberId456',
          },
        },
        user: {
          userId123: {
            id: 'userId123',
            service_member: 'serviceMemberId456',
          },
        },
        serviceMembers: {
          serviceMemberId456: {
            id: 'serviceMemberId456',
            affiliation: 'ARMY',
            edipi: '1234567890',
            first_name: 'Erin',
            last_name: 'Stanfill',
            middle_name: '',
            personal_email: 'erin@truss.works',
            phone_is_preferred: true,
            telephone: '555-555-5556',
            email_is_preferred: false,
            current_location: {
              id: 'testDutyLocationId',
              address: {
                city: 'Colorado Springs',
                country: 'United States',
                postalCode: '80913',
                state: 'CO',
                streetAddress1: 'n/a',
              },
            },
            residential_address: {
              city: 'Somewhere',
              postalCode: '80913',
              state: 'CO',
              streetAddress1: '123 Main',
            },
            backup_mailing_address: {
              city: 'Washington',
              postalCode: '20021',
              state: 'DC',
              streetAddress1: '200 K St',
            },
          },
        },
      },
    };

    expect(selectServiceMemberProfileState(testState)).toEqual(profileStates.BACKUP_ADDRESS_COMPLETE);
  });

  it('returns BACKUP_CONTACTS_COMPLETE if all data is complete', () => {
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
            affiliation: 'ARMY',
            backup_mailing_address: {
              city: 'Washington',
              postalCode: '20021',
              state: 'DC',
              streetAddress1: '200 K St',
            },
            created_at: '2018-05-25T15:48:49.918Z',
            current_location: {
              address: {
                city: 'Colorado Springs',
                country: 'United States',
                postalCode: '80913',
                state: 'CO',
                streetAddress1: 'n/a',
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
            residential_address: {
              city: 'Somewhere',
              postalCode: '80913',
              state: 'CO',
              streetAddress1: '123 Main',
            },
            telephone: '555-555-5556',
            updated_at: '2018-05-25T21:39:10.484Z',
            user_id: 'b46e651e-9d1c-4be5-bb88-bba58e817696',
            backup_contacts: ['backupContact789', 'backupContact8910'],
          },
        },
      },
    };

    expect(selectServiceMemberProfileState(testState)).toEqual(profileStates.BACKUP_CONTACTS_COMPLETE);
  });
});

describe('selectCurrentDutyLocation', () => {
  it('returns the service member’s current duty location', () => {
    const testState = {
      entities: {
        orders: {
          orders789: {
            id: 'orders789',
            service_member_id: 'serviceMemberId456',
            origin_duty_location: {
              address: {
                city: 'Colorado Springs',
                country: 'United States',
                postalCode: '80913',
                state: 'CO',
                streetAddress1: 'n/a',
              },
              affiliation: 'ARMY',
              created_at: '2018-05-20T18:36:45.034Z',
              id: '28f63a9d-8fff-4a0f-84ef-661c5c8c354e',
              name: 'Ft Carson',
              updated_at: '2018-05-20T18:36:45.034Z',
            },
          },
        },
        user: {
          userId123: {
            id: 'userId123',
            service_member: 'serviceMemberId456',
          },
        },
        serviceMembers: {
          serviceMemberId456: {
            id: 'serviceMemberId456',
            orders: ['orders789', 'orders8910'],
          },
        },
      },
    };

    expect(selectCurrentDutyLocation(testState)).toEqual(testState.entities.orders.orders789.origin_duty_location);
  });

  it('returns null if there is the service member has no current location', () => {
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

    expect(selectCurrentDutyLocation(testState)).toEqual(null);
  });
});

describe('selectBackupContacts', () => {
  it('returns the backup contacts associated with the logged in user', () => {
    const testState = {
      entities: {
        backupContacts: {
          backupContact789: {
            id: 'backupContact789',
            service_member_id: 'serviceMemberId456',
          },
          backupContact8910: {
            id: 'backupContact8910',
            service_member_id: 'serviceMemberId456',
          },
        },
        user: {
          userId123: {
            id: 'userId123',
            service_member: 'serviceMemberId456',
          },
        },
        serviceMembers: {
          serviceMemberId456: {
            id: 'serviceMemberId456',
            backup_contacts: ['backupContact789', 'backupContact8910'],
          },
        },
      },
    };

    expect(selectBackupContacts(testState)).toEqual([
      testState.entities.backupContacts.backupContact789,
      testState.entities.backupContacts.backupContact8910,
    ]);
  });

  it('returns an empty array if the service member has no backup contacts', () => {
    const testState = {
      entities: {
        backupContacts: {
          backupContact789: {
            id: 'backupContact789',
            service_member_id: 'serviceMemberId123',
          },
          backupContact8910: {
            id: 'backupContact8910',
            service_member_id: 'serviceMemberId123',
          },
        },
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

    expect(selectBackupContacts(testState)).toEqual([]);
  });
});

describe('selectIsProfileComplete', () => {
  const testServiceMember = {
    affiliation: 'ARMY',
    backup_mailing_address: {
      city: 'Washington',
      postalCode: '20021',
      state: 'DC',
      streetAddress1: '200 K St',
    },
    created_at: '2018-05-25T15:48:49.918Z',
    current_location: {
      address: {
        city: 'Colorado Springs',
        country: 'United States',
        postalCode: '80913',
        state: 'CO',
        streetAddress1: 'n/a',
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
    residential_address: {
      city: 'Somewhere',
      postalCode: '80913',
      state: 'CO',
      streetAddress1: '123 Main',
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

describe('selectOrdersForLoggedInUser', () => {
  it('returns the orders associated with the logged in user', () => {
    const testState = {
      entities: {
        orders: {
          orders789: {
            id: 'orders789',
            service_member_id: 'serviceMemberId456',
          },
          orders8910: {
            id: 'orders8910',
            service_member_id: 'serviceMemberId456',
          },
        },
        user: {
          userId123: {
            id: 'userId123',
            service_member: 'serviceMemberId456',
          },
        },
        serviceMembers: {
          serviceMemberId456: {
            id: 'serviceMemberId456',
            orders: ['orders789', 'orders8910'],
          },
        },
      },
    };

    expect(selectOrdersForLoggedInUser(testState)).toEqual([
      testState.entities.orders.orders789,
      testState.entities.orders.orders8910,
    ]);
  });

  it('returns an empty array if the service member has no orders', () => {
    const testState = {
      entities: {
        orders: {
          orders789: {
            id: 'orders789',
            service_member_id: 'serviceMemberId456',
          },
          orders8910: {
            id: 'orders8910',
            service_member_id: 'serviceMemberId456',
          },
        },
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

    expect(selectOrdersForLoggedInUser(testState)).toEqual([]);
  });
});

describe('selectCurrentOrders', () => {
  it('returns the current orders associated with the logged in user', () => {
    const testState = {
      entities: {
        moves: {
          move1029: {
            id: 'move1029',
            orders_id: 'orders789',
          },
          move2938: {
            id: 'move2938',
            orders_id: 'orders8910',
          },
        },
        orders: {
          orders789: {
            id: 'orders789',
            service_member_id: 'serviceMemberId456',
            moves: ['move1029'],
            status: 'CANCELED',
          },
          orders8910: {
            id: 'orders8910',
            service_member_id: 'serviceMemberId456',
            moves: ['move2938'],
            status: 'APPROVED',
          },
        },
        user: {
          userId123: {
            id: 'userId123',
            service_member: 'serviceMemberId456',
          },
        },
        serviceMembers: {
          serviceMemberId456: {
            id: 'serviceMemberId456',
            orders: ['orders789', 'orders8910'],
          },
        },
      },
    };

    expect(selectCurrentOrders(testState)).toEqual(testState.entities.orders.orders8910);
  });

  it('returns the first orders if none of the orders statuses are active', () => {
    const testState = {
      entities: {
        moves: {
          move1029: {
            id: 'move1029',
            orders_id: 'orders789',
          },
          move2938: {
            id: 'move2938',
            orders_id: 'orders8910',
          },
        },
        orders: {
          orders789: {
            id: 'orders789',
            service_member_id: 'serviceMemberId456',
            moves: ['move1029'],
            status: 'CANCELED',
          },
          orders8910: {
            id: 'orders8910',
            service_member_id: 'serviceMemberId456',
            moves: ['move2938'],
            status: 'CANCELED',
          },
        },
        user: {
          userId123: {
            id: 'userId123',
            service_member: 'serviceMemberId456',
          },
        },
        serviceMembers: {
          serviceMemberId456: {
            id: 'serviceMemberId456',
            orders: ['orders789', 'orders8910'],
          },
        },
      },
    };

    expect(selectCurrentOrders(testState)).toEqual(testState.entities.orders.orders789);
  });

  it('returns null if there are no orders', () => {
    const testState = {
      entities: {
        moves: {
          move1029: {
            id: 'move1029',
            orders_id: 'orders789',
          },
          move2938: {
            id: 'move2938',
            orders_id: 'orders8910',
          },
        },
        orders: {},
        user: {
          userId123: {
            id: 'userId123',
            service_member: 'serviceMemberId456',
          },
        },
        serviceMembers: {
          serviceMemberId456: {
            id: 'serviceMemberId456',
            orders: ['orders789', 'orders8910'],
          },
        },
      },
    };

    expect(selectCurrentOrders(testState)).toEqual(null);
  });
});

describe('selectMovesForLoggedInUser', () => {
  it('returns the moves associated with the logged in user', () => {
    const testState = {
      entities: {
        moves: {
          move1029: {
            id: 'move1029',
            orders_id: 'orders789',
            status: 'CANCELED',
          },
          move2938: {
            id: 'move2938',
            orders_id: 'orders8910',
          },
        },
        orders: {
          orders789: {
            id: 'orders789',
            service_member_id: 'serviceMemberId456',
            moves: ['move1029'],
          },
          orders8910: {
            id: 'orders8910',
            service_member_id: 'serviceMemberId456',
            moves: ['move2938'],
            status: 'SUBMITTED',
          },
        },
        user: {
          userId123: {
            id: 'userId123',
            service_member: 'serviceMemberId456',
          },
        },
        serviceMembers: {
          serviceMemberId456: {
            id: 'serviceMemberId456',
            orders: ['orders789', 'orders8910'],
          },
        },
      },
    };

    expect(selectMovesForLoggedInUser(testState)).toEqual([
      testState.entities.moves.move1029,
      testState.entities.moves.move2938,
    ]);
  });

  it('returns an empty array if the logged in user has no moves', () => {
    const testState = {
      entities: {
        orders: {
          orders789: {
            id: 'orders789',
            service_member_id: 'serviceMemberId456',
            moves: [],
          },
          orders8910: {
            id: 'orders8910',
            service_member_id: 'serviceMemberId456',
            status: 'SUBMITTED',
          },
        },
        user: {
          userId123: {
            id: 'userId123',
            service_member: 'serviceMemberId456',
          },
        },
        serviceMembers: {
          serviceMemberId456: {
            id: 'serviceMemberId456',
            orders: ['orders789', 'orders8910'],
          },
        },
      },
    };

    expect(selectMovesForLoggedInUser(testState)).toEqual([]);
  });
});

describe('selectMovesForCurrentOrders', () => {
  it('returns the moves associated with the current orders', () => {
    const testState = {
      entities: {
        moves: {
          move1029: {
            id: 'move1029',
            orders_id: 'orders789',
          },
          move2938: {
            id: 'move2938',
            orders_id: 'orders8910',
          },
        },
        orders: {
          orders789: {
            id: 'orders789',
            service_member_id: 'serviceMemberId456',
            moves: ['move1029'],
          },
          orders8910: {
            id: 'orders8910',
            service_member_id: 'serviceMemberId456',
            moves: ['move2938'],
            status: 'SUBMITTED',
          },
        },
        user: {
          userId123: {
            id: 'userId123',
            service_member: 'serviceMemberId456',
          },
        },
        serviceMembers: {
          serviceMemberId456: {
            id: 'serviceMemberId456',
            orders: ['orders789', 'orders8910'],
          },
        },
      },
    };

    expect(selectMovesForCurrentOrders(testState)).toEqual([testState.entities.moves.move2938]);
  });

  it('returns an empty array if the current orders have no moves', () => {
    const testState = {
      entities: {
        moves: {
          move1029: {
            id: 'move1029',
            orders_id: 'orders789',
          },
          move2938: {
            id: 'move2938',
            orders_id: 'orders8910',
          },
        },
        orders: {
          orders789: {
            id: 'orders789',
            service_member_id: 'serviceMemberId456',
            moves: ['move1029'],
          },
          orders8910: {
            id: 'orders8910',
            service_member_id: 'serviceMemberId456',
            moves: [],
            status: 'SUBMITTED',
          },
        },
        user: {
          userId123: {
            id: 'userId123',
            service_member: 'serviceMemberId456',
          },
        },
        serviceMembers: {
          serviceMemberId456: {
            id: 'serviceMemberId456',
            orders: ['orders789', 'orders8910'],
          },
        },
      },
    };

    expect(selectMovesForCurrentOrders(testState)).toEqual([]);
  });
});

describe('selectCurrentMove', () => {
  it('returns the current move associated with the current orders', () => {
    const testState = {
      entities: {
        moves: {
          move1029: {
            id: 'move1029',
            orders_id: 'orders789',
            status: 'CANCELED',
          },
          move2938: {
            id: 'move2938',
            orders_id: 'orders8910',
            status: 'DRAFT',
          },
        },
        orders: {
          orders789: {
            id: 'orders789',
            service_member_id: 'serviceMemberId456',
            moves: ['move1029'],
            status: 'CANCELED',
          },
          orders8910: {
            id: 'orders8910',
            service_member_id: 'serviceMemberId456',
            moves: ['move2938'],
            status: 'DRAFT',
          },
        },
        user: {
          userId123: {
            id: 'userId123',
            service_member: 'serviceMemberId456',
          },
        },
        serviceMembers: {
          serviceMemberId456: {
            id: 'serviceMemberId456',
            orders: ['orders789', 'orders8910'],
          },
        },
      },
    };

    expect(selectCurrentMove(testState)).toEqual(testState.entities.moves.move2938);
  });

  it('returns the current move associated with the first orders if none of the orders statuses are active', () => {
    const testState = {
      entities: {
        moves: {
          move1029: {
            id: 'move1029',
            orders_id: 'orders789',
            status: 'SUBMITTED',
          },
          move2938: {
            id: 'move2938',
            orders_id: 'orders8910',
            status: 'APPROVED',
          },
        },
        orders: {
          orders789: {
            id: 'orders789',
            service_member_id: 'serviceMemberId456',
            moves: ['move1029'],
            status: 'CANCELED',
          },
          orders8910: {
            id: 'orders8910',
            service_member_id: 'serviceMemberId456',
            moves: ['move2938'],
            status: 'CANCELED',
          },
        },
        user: {
          userId123: {
            id: 'userId123',
            service_member: 'serviceMemberId456',
          },
        },
        serviceMembers: {
          serviceMemberId456: {
            id: 'serviceMemberId456',
            orders: ['orders789', 'orders8910'],
          },
        },
      },
    };

    expect(selectCurrentMove(testState)).toEqual(testState.entities.moves.move1029);
  });

  it('returns the first move if there are no active moves', () => {
    const testState = {
      entities: {
        moves: {
          move1029: {
            id: 'move1029',
            orders_id: 'orders789',
            status: 'CANCELED',
          },
          move2938: {
            id: 'move2938',
            orders_id: 'orders8910',
            status: 'CANCELED',
          },
        },
        orders: {
          orders789: {
            id: 'orders789',
            service_member_id: 'serviceMemberId456',
            moves: ['move1029'],
            status: 'CANCELED',
          },
          orders8910: {
            id: 'orders8910',
            service_member_id: 'serviceMemberId456',
            moves: ['move2938'],
            status: 'CANCELED',
          },
        },
        user: {
          userId123: {
            id: 'userId123',
            service_member: 'serviceMemberId456',
          },
        },
        serviceMembers: {
          serviceMemberId456: {
            id: 'serviceMemberId456',
            orders: ['orders789', 'orders8910'],
          },
        },
      },
    };

    expect(selectCurrentMove(testState)).toEqual(testState.entities.moves.move1029);
  });
});
describe('selectPPMForMove', () => {
  it('returns the PPM associated with the given move ID', () => {
    const testState = {
      entities: {
        moves: {
          move1029: {
            id: 'move1029',
            orders_id: 'orders789',
            status: 'CANCELED',
          },
          move2938: {
            id: 'move2938',
            orders_id: 'orders8910',
            status: 'DRAFT',
          },
        },
        orders: {
          orders789: {
            id: 'orders789',
            service_member_id: 'serviceMemberId456',
            moves: ['move1029'],
            status: 'CANCELED',
          },
          orders8910: {
            id: 'orders8910',
            service_member_id: 'serviceMemberId456',
            moves: ['move2938'],
            status: 'DRAFT',
          },
        },
        personallyProcuredMoves: {
          ppmId789: {
            id: 'ppmId789',
            move_id: 'move2938',
            status: 'DRAFT',
          },
          ppmId910: {
            id: 'ppmId910',
            move_id: 'move1029',
            status: 'CANCELED',
          },
        },
        user: {
          userId123: {
            id: 'userId123',
            service_member: 'serviceMemberId456',
          },
        },
        serviceMembers: {
          serviceMemberId456: {
            id: 'serviceMemberId456',
            orders: ['orders789', 'orders8910'],
          },
        },
      },
    };

    expect(selectPPMForMove(testState, 'move2938')).toEqual(testState.entities.personallyProcuredMoves.ppmId789);
  });

  it('returns null if the PPM status is not active', () => {
    const testState = {
      entities: {
        moves: {
          move1029: {
            id: 'move1029',
            orders_id: 'orders789',
            status: 'CANCELED',
          },
          move2938: {
            id: 'move2938',
            orders_id: 'orders8910',
            status: 'DRAFT',
          },
        },
        orders: {
          orders789: {
            id: 'orders789',
            service_member_id: 'serviceMemberId456',
            moves: ['move1029'],
            status: 'CANCELED',
          },
          orders8910: {
            id: 'orders8910',
            service_member_id: 'serviceMemberId456',
            moves: ['move2938'],
            status: 'DRAFT',
          },
        },
        personallyProcuredMoves: {
          ppmId789: {
            id: 'ppmId789',
            move_id: 'move2938',
            status: 'CANCELED',
          },
          ppmId910: {
            id: 'ppmId910',
            move_id: 'move1029',
            status: 'CANCELED',
          },
        },
        user: {
          userId123: {
            id: 'userId123',
            service_member: 'serviceMemberId456',
          },
        },
        serviceMembers: {
          serviceMemberId456: {
            id: 'serviceMemberId456',
            orders: ['orders789', 'orders8910'],
          },
        },
      },
    };

    expect(selectPPMForMove(testState, 'move1029')).toEqual(null);
  });

  it('returns null if there is no PPM associated with the given move', () => {
    const testState = {
      entities: {
        moves: {
          move1029: {
            id: 'move1029',
            orders_id: 'orders789',
            status: 'CANCELED',
          },
          move2938: {
            id: 'move2938',
            orders_id: 'orders8910',
            status: 'DRAFT',
          },
        },
        orders: {
          orders789: {
            id: 'orders789',
            service_member_id: 'serviceMemberId456',
            moves: ['move1029'],
            status: 'CANCELED',
          },
          orders8910: {
            id: 'orders8910',
            service_member_id: 'serviceMemberId456',
            moves: ['move2938'],
            status: 'DRAFT',
          },
        },
        personallyProcuredMoves: {
          ppmId789: {
            id: 'ppmId789',
            move_id: 'move1111',
            status: 'DRAFT',
          },
          ppmId910: {
            id: 'ppmId910',
            move_id: 'move2222',
            status: 'SUBMITTED',
          },
        },
        user: {
          userId123: {
            id: 'userId123',
            service_member: 'serviceMemberId456',
          },
        },
        serviceMembers: {
          serviceMemberId456: {
            id: 'serviceMemberId456',
            orders: ['orders789', 'orders8910'],
          },
        },
      },
    };

    expect(selectPPMForMove(testState, 'move1029')).toEqual(null);
  });
});

describe('selectWeightTicketAndIndexById', () => {
  it('return the correct weight ticket and index', () => {
    const weightTicketId = '71422b71-a40b-41a7-b2ff-4da922a9c7f2';
    const testState = {
      entities: {
        mtoShipments: {
          '2ed2998e-ae36-46cd-af83-c3ecee55fe3e': {
            createdAt: '2022-07-01T01:10:51.224Z',
            eTag: 'MjAyMi0wNy0xMVQxODoyMDoxOC43MjQ1NzRa',
            id: '2ed2998e-ae36-46cd-af83-c3ecee55fe3e',
            moveTaskOrderID: '26b960d8-a96d-4450-a441-673ccd7cc3c7',
            ppmShipment: {
              actualDestinationPostalCode: '30813',
              actualMoveDate: '2022-07-31',
              actualPickupPostalCode: '90210',
              advanceAmountReceived: null,
              advanceAmountRequested: 598700,
              approvedAt: '2022-04-15T12:30:00.000Z',
              createdAt: '2022-07-01T01:10:51.231Z',
              eTag: 'MjAyMi0wNy0xMVQxODoyMDoxOC43NTIwMDNa',
              estimatedIncentive: 1000000,
              estimatedWeight: 4000,
              expectedDepartureDate: '2020-03-15',
              hasProGear: true,
              hasReceivedAdvance: false,
              hasRequestedAdvance: true,
              id: 'b9ae4c25-1376-4b9b-8781-106b5ae7ecab',
              proGearWeight: 1987,
              reviewedAt: null,
              shipmentId: '2ed2998e-ae36-46cd-af83-c3ecee55fe3e',
              sitEstimatedCost: null,
              sitEstimatedDepartureDate: null,
              sitEstimatedEntryDate: null,
              sitEstimatedWeight: null,
              sitExpected: false,
              spouseProGearWeight: 498,
              status: 'WAITING_ON_CUSTOMER',
              submittedAt: null,
              updatedAt: '2022-07-11T18:20:18.752Z',
              weightTickets: [
                {
                  id: 'd35d835f-8258-4266-87aa-54d61c917780',
                  emptyWeightDocumentId: '000676ac-c5ff-4630-8768-ef238f04e706',
                  fullWeightDocumentId: '7eeb270b-dc97-4f95-94c3-709c082cbf94',
                  trailerOwnershipDocumentId: 'd6b68bba-fe81-4402-82ac-6c02bf7cb660',
                },
                {
                  id: weightTicketId,
                  emptyWeightDocumentId: '15fdd562-82a9-4892-85d7-81cc3a85e68e',
                  fullWeightDocumentId: '4a7f7fd9-15d1-468f-9184-53d7c0c1ccdc',
                  trailerOwnershipDocumentId: 'f9ed20ad-bebd-4b5d-a59b-e3a86d273b78',
                },
              ],
            },
            shipmentType: 'PPM',
            status: 'APPROVED',
            updatedAt: '2022-07-11T18:20:18.724Z',
          },
        },
      },
    };
    const mtoShipmentID = Object.keys(testState.entities.mtoShipments)[0];

    expect(selectWeightTicketAndIndexById(testState, mtoShipmentID, weightTicketId)).toEqual({
      weightTicket: testState.entities.mtoShipments[mtoShipmentID].ppmShipment.weightTickets[1],
      index: 1,
    });
  });
});

describe('selectCurrentPPM', () => {
  it('returns the current PPM associated with the current move', () => {
    const testState = {
      entities: {
        moves: {
          move1029: {
            id: 'move1029',
            orders_id: 'orders789',
            status: 'CANCELED',
          },
          move2938: {
            id: 'move2938',
            orders_id: 'orders8910',
            status: 'DRAFT',
          },
        },
        orders: {
          orders789: {
            id: 'orders789',
            service_member_id: 'serviceMemberId456',
            moves: ['move1029'],
            status: 'CANCELED',
          },
          orders8910: {
            id: 'orders8910',
            service_member_id: 'serviceMemberId456',
            moves: ['move2938'],
            status: 'DRAFT',
          },
        },
        personallyProcuredMoves: {
          ppmId789: {
            id: 'ppmId789',
            move_id: 'move2938',
            status: 'DRAFT',
          },
          ppmId910: {
            id: 'ppmId910',
            move_id: 'move1029',
            status: 'CANCELED',
          },
        },
        user: {
          userId123: {
            id: 'userId123',
            service_member: 'serviceMemberId456',
          },
        },
        serviceMembers: {
          serviceMemberId456: {
            id: 'serviceMemberId456',
            orders: ['orders789', 'orders8910'],
          },
        },
      },
    };

    expect(selectCurrentPPM(testState)).toEqual(testState.entities.personallyProcuredMoves.ppmId789);
  });

  it('returns null if the PPM status is not active', () => {
    const testState = {
      entities: {
        moves: {
          move1029: {
            id: 'move1029',
            orders_id: 'orders789',
            status: 'CANCELED',
          },
          move2938: {
            id: 'move2938',
            orders_id: 'orders8910',
            status: 'DRAFT',
          },
        },
        orders: {
          orders789: {
            id: 'orders789',
            service_member_id: 'serviceMemberId456',
            moves: ['move1029'],
            status: 'CANCELED',
          },
          orders8910: {
            id: 'orders8910',
            service_member_id: 'serviceMemberId456',
            moves: ['move2938'],
            status: 'DRAFT',
          },
        },
        personallyProcuredMoves: {
          ppmId789: {
            id: 'ppmId789',
            move_id: 'move2938',
            status: 'CANCELED',
          },
          ppmId910: {
            id: 'ppmId910',
            move_id: 'move1029',
            status: 'CANCELED',
          },
        },
        user: {
          userId123: {
            id: 'userId123',
            service_member: 'serviceMemberId456',
          },
        },
        serviceMembers: {
          serviceMemberId456: {
            id: 'serviceMemberId456',
            orders: ['orders789', 'orders8910'],
          },
        },
      },
    };

    expect(selectCurrentPPM(testState)).toEqual(null);
  });

  it('returns null if there is no PPM associated with the current move', () => {
    const testState = {
      entities: {
        moves: {
          move1029: {
            id: 'move1029',
            orders_id: 'orders789',
            status: 'CANCELED',
          },
          move2938: {
            id: 'move2938',
            orders_id: 'orders8910',
            status: 'DRAFT',
          },
        },
        orders: {
          orders789: {
            id: 'orders789',
            service_member_id: 'serviceMemberId456',
            moves: ['move1029'],
            status: 'CANCELED',
          },
          orders8910: {
            id: 'orders8910',
            service_member_id: 'serviceMemberId456',
            moves: ['move2938'],
            status: 'DRAFT',
          },
        },
        personallyProcuredMoves: {
          ppmId789: {
            id: 'ppmId789',
            move_id: 'move1111',
            status: 'DRAFT',
          },
          ppmId910: {
            id: 'ppmId910',
            move_id: 'move2222',
            status: 'SUBMITTED',
          },
        },
        user: {
          userId123: {
            id: 'userId123',
            service_member: 'serviceMemberId456',
          },
        },
        serviceMembers: {
          serviceMemberId456: {
            id: 'serviceMemberId456',
            orders: ['orders789', 'orders8910'],
          },
        },
      },
    };

    expect(selectCurrentPPM(testState)).toEqual(null);
  });
});

describe('selectPPMSitEstimate', () => {
  it('returns the only PPM SIT estimate stored in entities', () => {
    const testState = {
      entities: {
        ppmSitEstimate: {
          undefined: {
            estimate: 12500,
          },
        },
      },
    };

    expect(selectPPMSitEstimate(testState)).toEqual(testState.entities.ppmSitEstimate.undefined.estimate);
  });

  it('returns null if there is no PPM SIT estimate in entities', () => {
    const testState = {
      entities: {
        ppmSitEstimate: {},
      },
    };

    expect(selectPPMSitEstimate(testState)).toEqual(null);
  });
});

describe('selectWeightAllotmentsForLoggedInUser', () => {
  describe('when I have dependents', () => {
    describe('when my spouse has pro gear', () => {
      it('should include spouse progear', () => {
        const testState = {
          entities: {
            orders: {
              orders8910: {
                id: 'orders8910',
                service_member_id: 'serviceMemberId456',
                moves: ['move2938'],
                status: 'DRAFT',
                has_dependents: true,
                spouse_has_pro_gear: true,
                authorizedWeight: 8000,
                entitlement: {
                  proGear: 2000,
                  proGearSpouse: 500,
                },
              },
            },
            user: {
              userId123: {
                id: 'userId123',
                service_member: 'serviceMemberId456',
              },
            },
            serviceMembers: {
              serviceMemberId456: {
                id: 'serviceMemberId456',
                orders: ['orders8910'],
              },
            },
          },
        };

        expect(selectWeightAllotmentsForLoggedInUser(testState)).toEqual({
          proGear: 2000,
          proGearSpouse: 500,
          sum: 10500,
          weight: 8000,
        });
      });
    });

    describe('when my spouse does not have pro gear', () => {
      it('should not include spouse progear', () => {
        const testState = {
          entities: {
            orders: {
              orders8910: {
                id: 'orders8910',
                service_member_id: 'serviceMemberId456',
                moves: ['move2938'],
                status: 'DRAFT',
                has_dependents: true,
                spouse_has_pro_gear: false,
                authorizedWeight: 8000,
                entitlement: {
                  proGear: 2000,
                  proGearSpouse: 0,
                },
              },
            },
            user: {
              userId123: {
                id: 'userId123',
                service_member: 'serviceMemberId456',
              },
            },
            serviceMembers: {
              serviceMemberId456: {
                id: 'serviceMemberId456',
                orders: ['orders8910'],
              },
            },
          },
        };

        expect(selectWeightAllotmentsForLoggedInUser(testState)).toEqual({
          proGear: 2000,
          proGearSpouse: 0,
          sum: 10000,
          weight: 8000,
        });
      });
    });
  });

  describe("when I don't have dependents", () => {
    it('should exclude spouse progear', () => {
      const testState = {
        entities: {
          orders: {
            orders8910: {
              id: 'orders8910',
              service_member_id: 'serviceMemberId456',
              moves: ['move2938'],
              status: 'DRAFT',
              has_dependents: false,
              spouse_has_pro_gear: false,
              authorizedWeight: 5000,
              entitlement: {
                proGear: 2000,
                proGearSpouse: 0,
              },
            },
          },
          user: {
            userId123: {
              id: 'userId123',
              service_member: 'serviceMemberId456',
            },
          },
          serviceMembers: {
            serviceMemberId456: {
              id: 'serviceMemberId456',
              orders: ['orders8910'],
            },
          },
        },
      };

      expect(selectWeightAllotmentsForLoggedInUser(testState)).toEqual({
        proGear: 2000,
        proGearSpouse: 0,
        sum: 7000,
        weight: 5000,
      });
    });
  });
});
