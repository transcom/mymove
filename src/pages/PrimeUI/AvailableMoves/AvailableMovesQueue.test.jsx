import React from 'react';
import { render, screen } from '@testing-library/react';

import PrimeSimulatorAvailableMoves from './AvailableMovesQueue';

import { usePrimeSimulatorAvailableMovesQueries } from 'hooks/queries';
import { MockProviders } from 'testUtils';

const mockNavigate = jest.fn();
jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useNavigate: () => mockNavigate,
}));

jest.mock('hooks/queries', () => ({
  usePrimeSimulatorAvailableMovesQueries: jest.fn(),
}));

jest.mock('services/primeApi', () => ({
  ...jest.requireActual('services/primeApi'),
}));

const moveTaskOrders = [
  {
    id: '1',
    moveCode: 'LN4T89',
    mtoShipments: [
      {
        id: '2',
        shipmentType: 'HHG',
        requestedPickupDate: '2021-11-26',
        pickupAddress: { streetAddress1: '100 1st Avenue', city: 'New York', state: 'NY', postalCode: '10001' },
        destinationAddress: { streetAddress1: '200 2nd Avenue', city: 'Buffalo', state: 'NY', postalCode: '1001' },
      },
      {
        id: '3',
        shipmentType: 'HHG_INTO_NTS_DOMESTIC',
        requestedPickupDate: '2021-12-01',
        pickupAddress: { streetAddress1: '800 Madison Avenue', city: 'New York', state: 'NY', postalCode: '10002' },
        destinationAddress: { streetAddress1: '200 2nd Avenue', city: 'Buffalo', state: 'NY', postalCode: '1001' },
      },
    ],
    mtoServiceItems: [
      { id: '4', reServiceCode: 'MS', reServiceName: 'Move management' },
      { id: '5', reServiceCode: 'DLH', mtoShipmentID: '2', reServiceName: 'Domestic linehaul' },
      { id: '6', reServiceCode: 'FSC', mtoShipmentID: '3', reServiceName: 'Fuel surcharge' },
    ],
  },
  {
    id: '2',
    moveCode: 'LN4T90',
    mtoShipments: [
      {
        id: '2',
        shipmentType: 'HHG',
        requestedPickupDate: '2021-11-26',
        pickupAddress: { streetAddress1: '100 1st Avenue', city: 'New York', state: 'NY', postalCode: '10001' },
        destinationAddress: { streetAddress1: '200 2nd Avenue', city: 'Buffalo', state: 'NY', postalCode: '1001' },
      },
      {
        id: '3',
        shipmentType: 'HHG_INTO_NTS_DOMESTIC',
        requestedPickupDate: '2021-12-01',
        pickupAddress: { streetAddress1: '800 Madison Avenue', city: 'New York', state: 'NY', postalCode: '10002' },
        destinationAddress: { streetAddress1: '200 2nd Avenue', city: 'Buffalo', state: 'NY', postalCode: '1001' },
      },
    ],
    mtoServiceItems: [
      { id: '4', reServiceCode: 'MS', reServiceName: 'Move management' },
      { id: '5', reServiceCode: 'DLH', mtoShipmentID: '2', reServiceName: 'Domestic linehaul' },
      { id: '6', reServiceCode: 'FSC', mtoShipmentID: '3', reServiceName: 'Fuel surcharge' },
    ],
  },
];

const renderWithProviders = () => {
  render(
    <MockProviders path="/simulator/moves/">
      <PrimeSimulatorAvailableMoves />
    </MockProviders>,
  );
};

describe('getPrimeAvailableMoves', () => {
  it('renders the loading text', async () => {
    renderWithProviders();

    expect(await screen.getByRole('heading', { name: 'Loading, please wait...', level: 2 }));
  });

  it('displays the prime simulator with table queue visible', async () => {
    usePrimeSimulatorAvailableMovesQueries.mockReturnValue(moveTaskOrders);
    renderWithProviders();

    const filterInput = screen.getByTestId('prime-date-filter-input');
    expect(filterInput).toBeInTheDocument();
  });
});
