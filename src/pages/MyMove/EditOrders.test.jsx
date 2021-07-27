import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import { EditOrders } from './EditOrders';

// import { getOrdersForServiceMember } from 'services/internalApi';

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
});
