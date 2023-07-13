import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import UploadServiceRequestDocuments from './UploadServiceRequestDocuments';

import { usePrimeSimulatorGetMove } from 'hooks/queries';
import { MockProviders } from 'testUtils';
import { primeSimulatorRoutes } from 'constants/routes';

const mockNavigate = jest.fn();
jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useNavigate: () => mockNavigate,
}));
const routingParams = { moveCodeOrID: 'LN4T89', mtoServiceItemId: 'test-payment-id-123' };

jest.mock('services/primeApi', () => ({
  ...jest.requireActual('services/primeApi'),
}));

jest.mock('hooks/queries', () => ({
  usePrimeSimulatorGetMove: jest.fn(),
}));

const moveTaskOrder = {
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
};

const moveReturnValue = {
  moveTaskOrder,
  isLoading: false,
  isError: false,
};

describe('Upload Service Request Documents Page', () => {
  it('renders the page without errors', () => {
    render(
      <MockProviders>
        <UploadServiceRequestDocuments />
      </MockProviders>,
    );

    expect(screen.getByText('Upload Service Request Document')).toBeInTheDocument();
  });

  it('navigates the user to the move details page when the cancel button is clicked', async () => {
    usePrimeSimulatorGetMove.mockReturnValue(moveReturnValue);
    render(
      <MockProviders path={primeSimulatorRoutes.UPLOAD_SERVICE_REQUEST_DOCUMENTS_PATH} params={routingParams}>
        <UploadServiceRequestDocuments />
      </MockProviders>,
    );

    const cancel = screen.getByRole('button', { name: 'Cancel' });
    await userEvent.click(cancel);

    await waitFor(() => {
      expect(mockNavigate).toHaveBeenCalledWith('/simulator/moves/LN4T89/details');
    });
  });
});
