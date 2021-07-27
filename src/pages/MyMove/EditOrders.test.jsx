import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import { EditOrders } from './EditOrders';

// import { getOrdersForServiceMember, patchOrders } from 'services/internalApi';

const mockPush = jest.fn();
const mockGoBack = jest.fn();

jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useLocation: () => ({
    pathname: 'localhost:3000/',
  }),
  useHistory: () => ({
    push: mockPush,
    goBack: mockGoBack,
  }),
}));

jest.mock('services/internalApi', () => ({
  ...jest.requireActual('services/internalApi'),
  getOrdersForServiceMember: jest.fn().mockImplementation(() => Promise.resolve()),
  patchOrders: jest.fn().mockImplementation(() => Promise.resolve()),
}));

describe('EditOrders Page', () => {
  const testProps = {
    moveIsApproved: false,
    serviceMemberId: 'id123',
    setFlashMessage: jest.fn(),
    updateOrders: jest.fn(),
    currentOrders: {
      moves: ['testMove'],
    },
    entitlement: {
      authorizedWeight: 5000,
      dependentsAuthorized: true,
      eTag: 'MjAyMC0wOS0xNFQxNzo0MTozOC42ODAwOVo=',
      id: '0dbc9029-dfc5-4368-bc6b-dfc95f5fe317',
      nonTemporaryStorage: true,
      privatelyOwnedVehicle: true,
      proGearWeight: 2000,
      proGearWeightSpouse: 500,
      storageInTransit: 2,
      totalDependents: 1,
      totalWeight: 5000,
    },
    existingUploads: [],
    schema: {},
    spouseHasProGear: false,
    context: { flags: { allOrdersTypes: true } },
  };

  it('renders the edit orders form', async () => {
    render(<EditOrders {...testProps} />);

    const h1 = await screen.findByRole('heading', { name: 'Orders', level: 1 });
    expect(h1).toBeInTheDocument();

    const editOrdersHeader = screen.getByRole('heading', { name: 'Edit Orders:', level: 2 });
    expect(editOrdersHeader).toBeInTheDocument();
  });

  it('goes back to the previous page when the cancel button is clicked', async () => {
    render(<EditOrders {...testProps} />);

    const cancel = await screen.findByText('Cancel');

    expect(cancel).toBeInTheDocument();

    userEvent.click(cancel);

    await waitFor(() => {
      expect(mockGoBack).toHaveBeenCalled();
    });
  });

  it('shows an error if the API returns an error', async () => {});

  it('next button patches the orders and goes to the previous page', async () => {
    const currentOrders = {
      currentOrders: {
        id: 'testOrdersId',
        orders_type: 'PERMANENT_CHANGE_OF_STATION',
        issue_date: '2020-11-08',
        report_by_date: '2020-11-26',
        has_dependents: false,
        new_duty_station: {
          address: {
            city: 'Des Moines',
            country: 'US',
            id: 'a4b30b99-4e82-48a6-b736-01662b499d6a',
            postal_code: '50309',
            state: 'IA',
            street_address_1: '987 Other Avenue',
            street_address_2: 'P.O. Box 1234',
            street_address_3: 'c/o Another Person',
          },
          address_id: 'a4b30b99-4e82-48a6-b736-01662b499d6a',
          affiliation: 'AIR_FORCE',
          created_at: '2020-10-19T17:01:16.114Z',
          id: 'f9299768-16d2-4a13-ae39-7087a58b1f62',
          name: 'Yuma AFB',
          updated_at: '2020-10-19T17:01:16.114Z',
        },
      },
    };

    render(<EditOrders {...testProps} {...currentOrders} />);
  });
});
