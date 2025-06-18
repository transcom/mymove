/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { fireEvent, render, screen } from '@testing-library/react';
import { useLocation } from 'react-router-dom';

import ConnectedProfile from './Profile';

import { customerRoutes } from 'constants/routes';
import { MockProviders } from 'testUtils';

jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useLocation: jest.fn(),
}));

describe('Profile component', () => {
  const testProps = {};

  it('renders the Profile Page', async () => {
    const mockState = {
      entities: {
        user: {
          testUserId: {
            id: 'testUserId',
            email: 'testuser@example.com',
            service_member: 'testServiceMemberId',
          },
        },
        orders: {
          test: {
            new_duty_location: {
              name: 'Test Duty Location',
            },
            status: 'DRAFT',
            moves: ['testMove'],
          },
        },
        moves: {
          testMove: {
            created_at: '2020-12-17T15:54:48.873Z',
            id: 'testMove',
            locator: 'test',
            orders_id: 'test',
            selected_move_type: '',
            service_member_id: 'testServiceMemberId',
            status: 'DRAFT',
          },
        },
        serviceMembers: {
          testServiceMemberId: {
            id: 'testServiceMemberId',
            edipi: '1234567890',
            affiliation: 'ARMY',
            first_name: 'Tester',
            last_name: 'Testperson',
            telephone: '1234567890',
            personal_email: 'test@example.com',
            email_is_preferred: true,
            residential_address: {
              city: 'San Diego',
              state: 'CA',
              postalCode: '92131',
              streetAddress1: 'Some Street',
              country: 'USA',
            },
            backup_mailing_address: {
              city: 'San Diego',
              state: 'CA',
              postalCode: '92131',
              streetAddress1: 'Some Backup Street',
              country: 'USA',
            },
            current_location: {
              origin_duty_location: {
                name: 'Current Station',
              },
              grade: 'E-5',
            },
            backup_contacts: [
              {
                name: 'Backup Contact',
                telephone: '555-555-5555',
                email: 'backup@test.com',
              },
            ],
            orders: ['test'],
          },
        },
      },
    };
    useLocation.mockReturnValue({ state: { moveId: 'test' } });

    render(
      <MockProviders initialState={mockState} path={customerRoutes.MOVE_HOME_PATH} params={{ moveId: 'testMoveId' }}>
        <ConnectedProfile {...testProps} />
      </MockProviders>,
    );

    const mainHeader = await screen.findByRole('heading', { name: 'Profile', level: 1 });

    expect(mainHeader).toBeInTheDocument();

    const contactInfoHeader = await screen.findByRole('heading', { name: 'Contact info', level: 2 });

    expect(contactInfoHeader).toBeInTheDocument();

    const serviceInfoHeader = await screen.findByRole('heading', { name: 'Service info', level: 2 });

    expect(serviceInfoHeader).toBeInTheDocument();

    const editLinks = screen.getAllByText('Edit');

    expect(editLinks.length).toBe(4);

    const homeLink = screen.getByText('Return to Move');

    expect(homeLink).toBeInTheDocument();

    // these should be false since needsToVerifyProfile is not true
    const returnToDashboardLink = screen.queryByText('Return to Dashboard');
    expect(returnToDashboardLink).not.toBeInTheDocument();

    const createMoveBtn = screen.queryByText('createMoveBtn');
    expect(createMoveBtn).not.toBeInTheDocument();

    const profileConfirmAlert = screen.queryByText('profileConfirmAlert');
    expect(profileConfirmAlert).not.toBeInTheDocument();
  });

  it('renders the Profile Page when there are no orders', async () => {
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
            edipi: '1234567890',
            affiliation: 'ARMY',
            first_name: 'Tester',
            last_name: 'Testperson',
            telephone: '1234567890',
            personal_email: 'test@example.com',
            email_is_preferred: true,
            residential_address: {
              city: 'San Diego',
              state: 'CA',
              postalCode: '92131',
              streetAddress1: 'Some Street',
              country: 'USA',
            },
            backup_mailing_address: {
              city: 'San Diego',
              state: 'CA',
              postalCode: '92131',
              streetAddress1: 'Some Backup Street',
              country: 'USA',
            },
            current_location: {
              origin_duty_location: {
                name: 'Current Station',
              },
              grade: 'E-5',
            },
            backup_contacts: [
              {
                name: 'Backup Contact',
                telephone: '555-555-5555',
                email: 'backup@test.com',
              },
            ],
          },
        },
      },
    };
    useLocation.mockReturnValue({ state: { moveId: 'test' } });

    render(
      <MockProviders initialState={mockState} path={customerRoutes.MOVE_HOME_PATH} params={{ moveId: 'testMoveId' }}>
        <ConnectedProfile {...testProps} />
      </MockProviders>,
    );

    const mainHeader = await screen.findByRole('heading', { name: 'Profile', level: 1 });

    expect(mainHeader).toBeInTheDocument();

    const contactInfoHeader = await screen.findByRole('heading', { name: 'Contact info', level: 2 });

    expect(contactInfoHeader).toBeInTheDocument();

    const serviceInfoHeader = await screen.findByRole('heading', { name: 'Service info', level: 2 });

    expect(serviceInfoHeader).toBeInTheDocument();

    const editLinks = screen.getAllByText('Edit');

    expect(editLinks.length).toBe(3);

    const homeLink = screen.getByText('Return to Move');

    expect(homeLink).toBeInTheDocument();

    expect(screen.queryByText('Contact your movers if you need to make changes to your move.')).not.toBeInTheDocument();

    expect(screen.queryByText(/To change information in this section, contact the/)).not.toBeInTheDocument();
  });

  it('does not allow the user to edit the service info information after a move has been submitted', async () => {
    const mockState = {
      entities: {
        user: {
          testUserId: {
            id: 'testUserId',
            email: 'testuser@example.com',
            service_member: 'testServiceMemberId',
          },
        },
        orders: {
          test: {
            new_duty_location: {
              name: 'Test Duty Location',
            },
            status: 'DRAFT',
            moves: ['testMove'],
            id: 'testOrder',
          },
        },
        moves: {
          testMove: {
            created_at: '2020-12-17T15:54:48.873Z',
            id: 'testMove',
            locator: 'test',
            orders_id: 'test',
            selected_move_type: '',
            service_member_id: 'testServiceMemberId',
            status: 'SUBMITTED',
          },
        },
        serviceMembers: {
          testServiceMemberId: {
            id: 'testServiceMemberId',
            edipi: '1234567890',
            affiliation: 'ARMY',
            first_name: 'Tester',
            last_name: 'Testperson',
            telephone: '1234567890',
            personal_email: 'test@example.com',
            email_is_preferred: true,
            residential_address: {
              city: 'San Diego',
              state: 'CA',
              postalCode: '92131',
              streetAddress1: 'Some Street',
              country: 'USA',
            },
            backup_mailing_address: {
              city: 'San Diego',
              state: 'CA',
              postalCode: '92131',
              streetAddress1: 'Some Backup Street',
              country: 'USA',
            },
            current_location: {
              origin_duty_location: {
                name: 'Current Station',
              },
              grade: 'E-5',
            },
            backup_contacts: [
              {
                name: 'Backup Contact',
                telephone: '555-555-5555',
                email: 'backup@test.com',
              },
            ],
            orders: ['test'],
          },
        },
      },
    };
    useLocation.mockReturnValue({ state: { moveId: 'test' } });

    render(
      <MockProviders initialState={mockState} path={customerRoutes.MOVE_HOME_PATH} params={{ moveId: 'testMoveId' }}>
        <ConnectedProfile {...testProps} />
      </MockProviders>,
    );

    const alert = screen.getByText(
      'You can change these details later by talking to a move counselor or customer care representative.',
    );

    expect(alert).toBeInTheDocument();

    const whoToContact = screen.getByText(/To change information in this section, contact the/);

    expect(whoToContact).toBeInTheDocument();

    const editLinks = screen.getAllByText('Edit');

    expect(editLinks.length).toBe(3);

    const homeLink = screen.getByText('Return to Move');

    expect(homeLink).toBeInTheDocument();
  });

  it('renders the Profile Page with needsToVerifyProfile set to true', async () => {
    const mockState = {
      entities: {
        user: {
          testUserId: {
            id: 'testUserId',
            email: 'testuser@example.com',
            service_member: 'testServiceMemberId',
          },
        },
        orders: {
          test: {
            new_duty_location: {
              name: 'Test Duty Location',
            },
            status: 'DRAFT',
            moves: ['testMove'],
          },
        },
        moves: {
          testMove: {
            created_at: '2020-12-17T15:54:48.873Z',
            id: 'testMove',
            locator: 'test',
            orders_id: 'test',
            selected_move_type: '',
            service_member_id: 'testServiceMemberId',
            status: 'DRAFT',
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
            residential_address: {
              city: 'San Diego',
              state: 'CA',
              postalCode: '92131',
              streetAddress1: 'Some Street',
              country: 'USA',
            },
            backup_mailing_address: {
              city: 'San Diego',
              state: 'CA',
              postalCode: '92131',
              streetAddress1: 'Some Backup Street',
              country: 'USA',
            },
            current_location: {
              origin_duty_location: {
                name: 'Current Station',
              },
              grade: 'E-5',
            },
            backup_contacts: [
              {
                name: 'Backup Contact',
                telephone: '555-555-5555',
                email: 'backup@test.com',
              },
            ],
            orders: ['test'],
          },
        },
      },
    };

    useLocation.mockReturnValue({ state: { needsToVerifyProfile: true, moveId: 'test' } });

    render(
      <MockProviders initialState={mockState} path={customerRoutes.MOVE_HOME_PATH} params={{ moveId: 'testMoveId' }}>
        <ConnectedProfile {...testProps} />
      </MockProviders>,
    );

    const returnToDashboardLink = screen.getByText('Return to Dashboard');
    expect(returnToDashboardLink).toBeInTheDocument();

    const validateProfileContainer = screen.getByTestId('validateProfileContainer');
    expect(validateProfileContainer).toBeInTheDocument();

    const createMoveBtn = screen.getByTestId('createMoveBtn');
    expect(createMoveBtn).toBeInTheDocument();
    expect(createMoveBtn).toBeDisabled();

    const validateProfileBtn = screen.getByTestId('validateProfileBtn');
    expect(validateProfileBtn).toBeInTheDocument();
    expect(validateProfileBtn).toBeEnabled();

    const profileConfirmAlert = screen.getByTestId('profileConfirmAlert');
    expect(profileConfirmAlert).toBeInTheDocument();

    // user validates their profile, which enables create move btn
    fireEvent.click(validateProfileBtn);
    expect(createMoveBtn).toBeEnabled();
    expect(validateProfileBtn).toBeDisabled();
    expect(screen.getByText('Profile Validated')).toBeInTheDocument();
  });
});
