import React from 'react';
import { render, screen } from '@testing-library/react';

import { usePrimeSimulatorGetMove } from '../../../hooks/queries';

import PrimeUIUpdateServiceItems from './PrimeUIUpdateServiceItems';

import { ReactQueryWrapper, MockProviders } from 'testUtils';
import { primeSimulatorRoutes } from 'constants/routes';

jest.mock('hooks/queries', () => ({
  usePrimeSimulatorGetMove: jest.fn(),
}));

const routingParams = { moveCodeOrID: 'bf2fc98f-3cb5-40a0-a125-4c222096c35b' };

const moveTaskOrder = {
  id: '1',
  moveCode: 'LN4T89',
  mtoShipments: [
    {
      id: '2',
      shipmentType: 'HHG',
      requestedPickupDate: '2021-11-26',
      pickupAddress: { streetAddress1: '100 1st Avenue', city: 'New York', state: 'NY', postalCode: '10001' },
    },
  ],
  mtoServiceItems: [
    {
      reServiceCode: 'DDDSIT',
      reason: 'Holiday break',
      sitDestinationFinalAddress: {
        streetAddress1: '444 Main Ave',
        streetAddress2: 'Apartment 9000',
        streetAddress3: 'c/o Some Person',
        city: 'Anytown',
        state: 'AL',
        postalCode: '90210',
      },
      id: '45fe9475-d592-48f5-896a-45d4d6eb7e76',
    },
  ],
};

const moveReturnValue = {
  moveTaskOrder,
  isLoading: false,
  isError: false,
};

const renderComponent = () => {
  render(
    <ReactQueryWrapper>
      <MockProviders path={primeSimulatorRoutes.UPDATE_SERVICE_ITEMS_PATH} params={routingParams}>
        <PrimeUIUpdateServiceItems />
      </MockProviders>
    </ReactQueryWrapper>,
  );
};

describe('PrimeUIUpdateServiceItems page', () => {
  it('renders the update service items page', async () => {
    usePrimeSimulatorGetMove.mockReturnValue(moveReturnValue);

    renderComponent();

    expect(screen.getByRole('heading', { name: 'Update Service Items', level: 1 })).toBeInTheDocument();
  });
});
