/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { fireEvent, render, screen } from '@testing-library/react';
import { useLocation } from 'react-router-dom';

import ConnectedProfile from './Profile';

import { getAllMoves } from 'services/internalApi';
import { customerRoutes } from 'constants/routes';
import { MockProviders } from 'testUtils';

jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useLocation: jest.fn(),
}));

jest.mock('services/internalApi', () => ({
  ...jest.requireActual('services/internalApi'),
  getAllMoves: jest.fn().mockImplementation(() => Promise.resolve()),
}));

const serviceMemberMoves = {
  currentMove: [
    {
      lockExpiresAt: '2099-04-07T17:21:30.450Z',
      createdAt: '2024-02-23T19:30:11.374Z',
      eTag: 'MjAyNC0wMi0yM1QxOTozMDoxMS4zNzQxN1o=',
      id: 'testMoveId',
      moveCode: '44649B',
      orders: {
        authorizedWeight: 11000,
        created_at: '2024-02-23T19:30:11.369Z',
        entitlement: {
          proGear: 2000,
          proGearSpouse: 500,
        },
        grade: 'E_7',
        has_dependents: false,
        id: 'testOrders1',
        issue_date: '2024-02-29',
        new_duty_location: {
          address: {
            city: 'Fort Irwin',
            country: 'United States',
            id: '77dca457-d0d6-4718-9ca4-a630b4614cf8',
            postalCode: '92310',
            state: 'CA',
            streetAddress1: 'n/a',
          },
          address_id: '77dca457-d0d6-4718-9ca4-a630b4614cf8',
          affiliation: 'ARMY',
          created_at: '2024-02-22T21:34:21.449Z',
          id: '12421bcb-2ded-4165-b0ac-05f76301082a',
          name: 'Fort Irwin, CA 92310',
          transportation_office: {
            address: {
              city: 'Fort Irwin',
              country: 'United States',
              id: '65a97b21-cf6a-47c1-a4b6-e3f885dacba5',
              postalCode: '92310',
              state: 'CA',
              streetAddress1: 'Langford Lake Rd',
              streetAddress2: 'Bldg 105',
            },
            created_at: '2018-05-28T14:27:37.312Z',
            gbloc: 'LKNQ',
            id: 'd00e3ee8-baba-4991-8f3b-86c2e370d1be',
            name: 'PPPO Fort Irwin - USA',
            phone_lines: [],
            updated_at: '2018-05-28T14:27:37.312Z',
          },
          transportation_office_id: 'd00e3ee8-baba-4991-8f3b-86c2e370d1be',
          updated_at: '2024-02-22T21:34:21.449Z',
        },
        originDutyLocationGbloc: 'BGAC',
        origin_duty_location: {
          address: {
            city: 'Fort Gregg-Adams',
            country: 'United States',
            id: '12270b68-01cf-4416-8b19-125d11bc8340',
            postalCode: '23801',
            state: 'VA',
            streetAddress1: 'n/a',
          },
          address_id: '12270b68-01cf-4416-8b19-125d11bc8340',
          affiliation: 'ARMY',
          created_at: '2024-02-22T21:34:26.430Z',
          id: '9cf15b8d-985b-4ca3-9f27-4ba32a263908',
          name: 'Fort Gregg-Adams, VA 23801',
          transportation_office: {
            address: {
              city: 'Fort Gregg-Adams',
              country: 'United States',
              id: '10dc88f5-d76a-427f-89a0-bf85587b0570',
              postalCode: '23801',
              state: 'VA',
              streetAddress1: '1401 B Ave',
              streetAddress2: 'Bldg 3400, Room 119',
            },
            created_at: '2018-05-28T14:27:42.125Z',
            gbloc: 'BGAC',
            id: '4cc26e01-f0ea-4048-8081-1d179426a6d9',
            name: 'PPPO Fort Gregg-Adams - USA',
            phone_lines: [],
            updated_at: '2018-05-28T14:27:42.125Z',
          },
          transportation_office_id: '4cc26e01-f0ea-4048-8081-1d179426a6d9',
          updated_at: '2024-02-22T21:34:26.430Z',
        },
        report_by_date: '2024-02-29',
        service_member_id: '81aeac60-80f3-44d1-9b74-ba6d405ee2da',
        spouse_has_pro_gear: false,
        status: 'DRAFT',
        updated_at: '2024-02-23T19:30:11.369Z',
        uploaded_orders: {
          id: 'bd35c4c2-41c6-44a1-bf54-9098c68d87cc',
          service_member_id: '81aeac60-80f3-44d1-9b74-ba6d405ee2da',
          uploads: [
            {
              bytes: 92797,
              contentType: 'image/png',
              createdAt: '2024-02-26T18:43:58.515Z',
              filename: 'Screenshot 2024-02-08 at 12.57.43 PM.png',
              id: '786237dc-c240-449d-8859-3f37583b3406',
              status: 'PROCESSING',
              updatedAt: '2024-02-26T18:43:58.515Z',
              url: '/storage/user/5fe4d948-aa1c-4823-8967-b1fb40cf6679/uploads/786237dc-c240-449d-8859-3f37583b3406?contentType=image%2Fpng',
            },
          ],
        },
      },
      status: 'DRAFT',
      submittedAt: '0001-01-01T00:00:00.000Z',
      updatedAt: '0001-01-01T00:00:00.000Z',
    },
  ],
  previousMoves: [],
};

describe('Profile component', () => {
  const testProps = {};
  const multiMove = process.env.FEATURE_FLAG_MULTI_MOVE;

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

    if (multiMove) {
      render(
        <MockProviders initialState={mockState} path={customerRoutes.MOVE_HOME_PATH} params={{ moveId: 'testMoveId' }}>
          <ConnectedProfile {...testProps} />
        </MockProviders>,
      );
    } else {
      render(
        <MockProviders initialState={mockState}>
          <ConnectedProfile {...testProps} />
        </MockProviders>,
      );
    }

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

  it('renders the Profile Page with disabled edit buttons when the move has been locked by an office user', async () => {
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
    getAllMoves.mockResolvedValue(serviceMemberMoves);
    useLocation.mockReturnValue({ state: { moveId: 'testMoveId' } });

    if (multiMove) {
      render(
        <MockProviders initialState={mockState} path={customerRoutes.MOVE_HOME_PATH} params={{ moveId: 'testMoveId' }}>
          <ConnectedProfile {...testProps} />
        </MockProviders>,
      );
    } else {
      render(
        <MockProviders initialState={mockState}>
          <ConnectedProfile {...testProps} />
        </MockProviders>,
      );
    }

    const mainHeader = await screen.findByRole('heading', { name: 'Profile', level: 1 });

    expect(mainHeader).toBeInTheDocument();

    const contactInfoHeader = await screen.findByRole('heading', { name: 'Contact info', level: 2 });

    expect(contactInfoHeader).toBeInTheDocument();

    const serviceInfoHeader = await screen.findByRole('heading', { name: 'Service info', level: 2 });

    expect(serviceInfoHeader).toBeInTheDocument();

    const editLinks = screen.getAllByText('Edit');

    expect(editLinks.length).toBe(1);

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

    if (multiMove) {
      render(
        <MockProviders initialState={mockState} path={customerRoutes.MOVE_HOME_PATH} params={{ moveId: 'testMoveId' }}>
          <ConnectedProfile {...testProps} />
        </MockProviders>,
      );
    } else {
      render(
        <MockProviders initialState={mockState}>
          <ConnectedProfile {...testProps} />
        </MockProviders>,
      );
    }

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

    if (multiMove) {
      render(
        <MockProviders initialState={mockState} path={customerRoutes.MOVE_HOME_PATH} params={{ moveId: 'testMoveId' }}>
          <ConnectedProfile {...testProps} />
        </MockProviders>,
      );
    } else {
      render(
        <MockProviders initialState={mockState}>
          <ConnectedProfile {...testProps} />
        </MockProviders>,
      );
    }

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

    if (multiMove) {
      render(
        <MockProviders initialState={mockState} path={customerRoutes.MOVE_HOME_PATH} params={{ moveId: 'testMoveId' }}>
          <ConnectedProfile {...testProps} />
        </MockProviders>,
      );
    } else {
      render(
        <MockProviders initialState={mockState}>
          <ConnectedProfile {...testProps} />
        </MockProviders>,
      );
    }

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
