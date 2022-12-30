import React from 'react';
import { mount } from 'enzyme';

import requireCustomerStateHOC, { getIsAllowedProfileState } from './requireCustomerState';

import { MockProviders } from 'testUtils';
import { profileStates } from 'constants/customerStates';

const mockNavigate = jest.fn();
jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useNavigate: () => mockNavigate,
}));

beforeEach(() => {
  jest.resetAllMocks();
});

describe('getIsAllowedProfileState', () => {
  it('returns true for a required state that is before the current state', () => {
    const requiredState = profileStates.DOD_INFO_COMPLETE;
    const currentState = profileStates.ADDRESS_COMPLETE;
    const result = getIsAllowedProfileState(requiredState, currentState);
    expect(result).toBe(true);
  });

  it('returns false if the required state is after the current state', () => {
    const requiredState = profileStates.ADDRESS_COMPLETE;
    const currentState = profileStates.DOD_INFO_COMPLETE;
    const result = getIsAllowedProfileState(requiredState, currentState);
    expect(result).toBe(false);
  });
  it('returns true if the required state and current state are the same and profile is not complete', () => {
    const requiredState = profileStates.ADDRESS_COMPLETE;
    const currentState = profileStates.ADDRESS_COMPLETE;
    const result = getIsAllowedProfileState(requiredState, currentState);
    expect(result).toBe(true);
  });
  it('returns false if the current state is a completed profile and required state is not', () => {
    const requiredState = profileStates.ADDRESS_COMPLETE;
    const currentState = profileStates.BACKUP_CONTACTS_COMPLETE;
    const result = getIsAllowedProfileState(requiredState, currentState);
    expect(result).toBe(false);
  });
  it('returns true if the required state is a completed profile', () => {
    const requiredState = profileStates.BACKUP_CONTACTS_COMPLETE;
    const currentState = profileStates.BACKUP_CONTACTS_COMPLETE;
    const result = getIsAllowedProfileState(requiredState, currentState);
    expect(result).toBe(true);
  });
});

describe('requireCustomerState HOC', () => {
  const TestComponent = () => <div>My test component</div>;
  const TestComponentWithHOC = requireCustomerStateHOC(TestComponent, profileStates.ADDRESS_COMPLETE);

  it('dispatches a redirect if the current state is earlier than the required state', () => {
    const mockState = {
      entities: {
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
      },
    };

    const wrapper = mount(
      <MockProviders initialState={mockState}>
        <TestComponentWithHOC />
      </MockProviders>,
    );

    expect(wrapper.exists()).toBe(true);
    expect(mockNavigate).toHaveBeenCalledWith('/service-member/conus-oconus');
  });

  it('does not redirect if the current state equals the required state', () => {
    const mockState = {
      entities: {
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
            rank: 'test rank',
            edipi: '1234567890',
            affiliation: 'ARMY',
            first_name: 'Tester',
            last_name: 'Testperson',
            telephone: '1234567890',
            personal_email: 'test@example.com',
            email_is_preferred: true,
            current_location: {
              id: 'testDutyLocationId',
            },
            residential_address: {
              street: '123 Main St',
            },
          },
        },
      },
    };

    const wrapper = mount(
      <MockProviders initialState={mockState}>
        <TestComponentWithHOC />
      </MockProviders>,
    );

    expect(wrapper.exists()).toBe(true);
    expect(mockNavigate).not.toHaveBeenCalled();
  });
  it('does not redirect if the current state is after the required state but profile is not complete', () => {
    const mockState = {
      entities: {
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
            rank: 'test rank',
            edipi: '1234567890',
            affiliation: 'ARMY',
            first_name: 'Tester',
            last_name: 'Testperson',
            telephone: '1234567890',
            personal_email: 'test@example.com',
            email_is_preferred: true,
            current_location: {
              id: 'testDutyLocationId',
            },
            residential_address: {
              street: '123 Main St',
            },
            backup_mailing_address: {
              street: '456 Main St',
            },
          },
        },
      },
    };

    const wrapper = mount(
      <MockProviders initialState={mockState}>
        <TestComponentWithHOC />
      </MockProviders>,
    );

    expect(wrapper.exists()).toBe(true);
    expect(mockNavigate).not.toHaveBeenCalled();
  });

  it('does redirect if profile is complete and required state is not the completed profile state', () => {
    const mockState = {
      entities: {
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
            rank: 'test rank',
            edipi: '1234567890',
            affiliation: 'ARMY',
            first_name: 'Tester',
            last_name: 'Testperson',
            telephone: '1234567890',
            personal_email: 'test@example.com',
            email_is_preferred: true,
            current_location: {
              id: 'testDutyLocationId',
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
          },
        },
      },
    };

    const wrapper = mount(
      <MockProviders initialState={mockState}>
        <TestComponentWithHOC />
      </MockProviders>,
    );

    expect(wrapper.exists()).toBe(true);
    expect(mockNavigate).toHaveBeenCalledWith('/');
  });

  it('does not redirect if profile is complete and required state is the completed profile state', () => {
    const mockState = {
      entities: {
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
            rank: 'test rank',
            edipi: '1234567890',
            affiliation: 'ARMY',
            first_name: 'Tester',
            last_name: 'Testperson',
            telephone: '1234567890',
            personal_email: 'test@example.com',
            email_is_preferred: true,
            current_location: {
              id: 'testDutyLocationId',
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
          },
        },
      },
    };

    const TestComponentCompletedProfileWithHOC = requireCustomerStateHOC(
      TestComponent,
      profileStates.BACKUP_CONTACTS_COMPLETE,
    );
    const wrapper = mount(
      <MockProviders initialState={mockState}>
        <TestComponentCompletedProfileWithHOC />
      </MockProviders>,
    );

    expect(wrapper.exists()).toBe(true);
    expect(mockNavigate).not.toHaveBeenCalled();
  });
});
