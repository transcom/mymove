import {
  selectLoggedInUser,
  selectServiceMemberFromLoggedInUser,
  selectServiceMemberProfileState,
  selectIsProfileComplete,
  selectBackupContacts,
  selectCurrentDutyStation,
  selectOrdersForLoggedInUser,
  selectCurrentOrders,
  selectMovesForLoggedInUser,
  selectMovesForCurrentOrders,
  selectCurrentMove,
  selectCurrentPPM,
  selectPPMForMove,
  selectPPMEstimateRange,
  selectPPMSitEstimate,
  selectReimbursementById,
  selectEntitlementsForLoggedInUser,
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
            rank: 'O_4_W_4',
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
            rank: 'O_4_W_4',
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

  it('returns CONTACT_INFO_COMPLETE if there is no duty station data', () => {
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
            rank: 'O_4_W_4',
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

  it('returns DUTY_STATION_COMPLETE if there is no address data', () => {
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
            rank: 'O_4_W_4',
            edipi: '1234567890',
            first_name: 'Erin',
            last_name: 'Stanfill',
            middle_name: '',
            personal_email: 'erin@truss.works',
            phone_is_preferred: true,
            telephone: '555-555-5556',
            email_is_preferred: false,
            current_station: {
              id: 'testDutyStationId',
              address: {
                city: 'Colorado Springs',
                country: 'United States',
                postal_code: '80913',
                state: 'CO',
                street_address_1: 'n/a',
              },
            },
          },
        },
      },
    };

    expect(selectServiceMemberProfileState(testState)).toEqual(profileStates.DUTY_STATION_COMPLETE);
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
            rank: 'O_4_W_4',
            edipi: '1234567890',
            first_name: 'Erin',
            last_name: 'Stanfill',
            middle_name: '',
            personal_email: 'erin@truss.works',
            phone_is_preferred: true,
            telephone: '555-555-5556',
            email_is_preferred: false,
            current_station: {
              id: 'testDutyStationId',
              address: {
                city: 'Colorado Springs',
                country: 'United States',
                postal_code: '80913',
                state: 'CO',
                street_address_1: 'n/a',
              },
            },
            residential_address: {
              city: 'Somewhere',
              postal_code: '80913',
              state: 'CO',
              street_address_1: '123 Main',
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
            rank: 'O_4_W_4',
            edipi: '1234567890',
            first_name: 'Erin',
            last_name: 'Stanfill',
            middle_name: '',
            personal_email: 'erin@truss.works',
            phone_is_preferred: true,
            telephone: '555-555-5556',
            email_is_preferred: false,
            current_station: {
              id: 'testDutyStationId',
              address: {
                city: 'Colorado Springs',
                country: 'United States',
                postal_code: '80913',
                state: 'CO',
                street_address_1: 'n/a',
              },
            },
            residential_address: {
              city: 'Somewhere',
              postal_code: '80913',
              state: 'CO',
              street_address_1: '123 Main',
            },
            backup_mailing_address: {
              city: 'Washington',
              postal_code: '20021',
              state: 'DC',
              street_address_1: '200 K St',
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
            backup_contacts: ['backupContact789', 'backupContact8910'],
          },
        },
      },
    };

    expect(selectServiceMemberProfileState(testState)).toEqual(profileStates.BACKUP_CONTACTS_COMPLETE);
  });
});

describe('selectCurrentDutyStation', () => {
  it('returns the service member’s current duty station', () => {
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
            current_station: {
              id: 'dutyStationId890',
            },
          },
        },
      },
    };

    expect(selectCurrentDutyStation(testState)).toEqual(
      testState.entities.serviceMembers.serviceMemberId456.current_station,
    );
  });

  it('returns null if there is the service member has no current station', () => {
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

    expect(selectCurrentDutyStation(testState)).toEqual(null);
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

describe('selectPPMEstimateRange', () => {
  it('returns the only PPM estimate range stored in entities', () => {
    const testState = {
      entities: {
        ppmEstimateRanges: {
          undefined: {
            range_min: 1000,
            range_max: 2400,
          },
        },
      },
    };

    expect(selectPPMEstimateRange(testState)).toEqual(testState.entities.ppmEstimateRanges.undefined);
  });

  it('returns null if there is no PPM estimate range in entities', () => {
    const testState = {
      entities: {
        ppmEstimateRanges: {},
      },
    };

    expect(selectPPMEstimateRange(testState)).toEqual(null);
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

describe('selectReimbursementById', () => {
  it('returns the only PPM SIT estimate stored in entities', () => {
    const testState = {
      entities: {
        reimbursements: {
          testReimbursement123: {
            id: 'testReimbursement123',
          },
        },
      },
    };

    expect(selectReimbursementById(testState, 'testReimbursement123')).toEqual(
      testState.entities.reimbursements.testReimbursement123,
    );
  });

  it('returns null if there is no reimbursement in entities', () => {
    const testState = {
      entities: {
        ppmSitEstimate: {},
      },
    };

    expect(selectReimbursementById(testState, 'testReimbursement123')).toEqual(null);
  });
});

describe('selectEntitlementsForLoggedInUser', () => {
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
                weight_allotment: {
                  total_weight_self: 5000,
                  total_weight_self_plus_dependents: 8000,
                  pro_gear_weight: 2000,
                  pro_gear_weight_spouse: 500,
                },
              },
            },
          },
        };

        expect(selectEntitlementsForLoggedInUser(testState)).toEqual({
          pro_gear: 2000,
          pro_gear_spouse: 500,
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
                weight_allotment: {
                  total_weight_self: 5000,
                  total_weight_self_plus_dependents: 8000,
                  pro_gear_weight: 2000,
                  pro_gear_weight_spouse: 500,
                },
              },
            },
          },
        };

        expect(selectEntitlementsForLoggedInUser(testState)).toEqual({
          pro_gear: 2000,
          pro_gear_spouse: 0,
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
              weight_allotment: {
                total_weight_self: 5000,
                total_weight_self_plus_dependents: 8000,
                pro_gear_weight: 2000,
                pro_gear_weight_spouse: 500,
              },
            },
          },
        },
      };

      expect(selectEntitlementsForLoggedInUser(testState)).toEqual({
        pro_gear: 2000,
        pro_gear_spouse: 0,
        sum: 7000,
        weight: 5000,
      });
    });
  });
});
