import React from 'react';
import { mount } from 'enzyme';
import * as reactRedux from 'react-redux';
import { push } from 'connected-react-router';

import requireCustomerStateHOC, { getIsAllowedProfileState } from './requireCustomerState';

import { MockProviders } from 'testUtils';
import { profileStates } from 'constants/customerStates';

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
  const useDispatchMock = jest.spyOn(reactRedux, 'useDispatch');
  const mockDispatch = jest.fn();

  beforeEach(() => {
    useDispatchMock.mockClear();
    mockDispatch.mockClear();
    useDispatchMock.mockReturnValue(mockDispatch);
  });

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
    expect(mockDispatch).toHaveBeenCalledWith(push('/service-member/conus-oconus'));
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
            current_station: {
              id: 'testDutyStationId',
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
    expect(mockDispatch).not.toHaveBeenCalled();
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
            current_station: {
              id: 'testDutyStationId',
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
    expect(mockDispatch).not.toHaveBeenCalled();
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
    expect(mockDispatch).toHaveBeenCalledWith(push('/'));
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
    expect(mockDispatch).not.toHaveBeenCalled();
  });
});
