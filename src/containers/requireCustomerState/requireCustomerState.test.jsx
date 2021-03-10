import React from 'react';
import { mount } from 'enzyme';
import * as reactRedux from 'react-redux';
import { push } from 'connected-react-router';

import requireCustomerState from './requireCustomerState';

import { MockProviders } from 'testUtils';
import { profileStates } from 'constants/customerStates';

describe('requireCustomerState HOC', () => {
  const useDispatchMock = jest.spyOn(reactRedux, 'useDispatch');
  const mockDispatch = jest.fn();

  beforeEach(() => {
    useDispatchMock.mockClear();
    mockDispatch.mockClear();
    useDispatchMock.mockReturnValue(mockDispatch);
  });

  const TestComponent = () => <div>My test component</div>;
  const TestComponentWithHOC = requireCustomerState(TestComponent, profileStates.BACKUP_CONTACTS_COMPLETE);

  it('dispatches a redirect if the current state does not equal the required state', () => {
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
    expect(mockDispatch).not.toHaveBeenCalled();
  });
});
