import { selectLoggedInUser, selectServiceMemberFromLoggedInUser } from './selectors';

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
