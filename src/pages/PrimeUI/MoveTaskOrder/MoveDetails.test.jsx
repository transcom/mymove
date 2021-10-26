import React from 'react';
import { screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import { usePrimeSimulatorGetMove } from '../../../hooks/queries';

import MoveDetails from './MoveDetails';

import { renderWithRouter } from 'testUtils';

const mockUseHistoryPush = jest.fn();
const mockRequestedMoveCode = 'LN4T89';

jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useParams: jest.fn().mockReturnValue({ moveCodeOrID: 'LN4T89' }),
  useHistory: () => ({
    push: mockUseHistoryPush,
  }),
}));

jest.mock('hooks/queries', () => ({
  usePrimeSimulatorGetMove: jest.fn(),
}));

jest.mock('services/primeApi', () => ({
  ...jest.requireActual('services/primeApi'),
  moveDetails: jest.fn(),
}));

const moveTaskOrder = {
  id: '1',
  moveCode: mockRequestedMoveCode,
  mtoShipments: [
    {
      id: '2',
      shipmentType: 'HHG',
      requestedPickupDate: '2021-11-26',
      pickupAddress: { streetAddress1: '100 1st Avenue', city: 'New York', state: 'NY', postalCode: '10001' },
      destinationAddress: {
        streetAddress1: '800 Madison Avenue',
        streetAddress2: '900 Madison Avenue',
        city: 'New York',
        state: 'NY',
        postalCode: '10002',
      },
    },
    {
      id: '3',
      shipmentType: 'HHG_INTO_NTS_DOMESTIC',
      requestedPickupDate: '2021-12-01',
      pickupAddress: { streetAddress1: '800 Madison Avenue', city: 'New York', state: 'NY', postalCode: '10002' },
      destinationAddress: {
        streetAddress1: '800 Madison Avenue',
        streetAddress2: '900 Madison Avenue',
        city: 'New York',
        state: 'NY',
        postalCode: '10002',
      },
    },
  ],
  paymentRequests: [
    {
      id: '4a1b0048-ffe7-11eb-9a03-0242ac130003',
      paymentRequestNumber: '5924-0164-1',
    },
  ],
};

const moveReturnValue = {
  moveTaskOrder,
  isLoading: false,
  isError: false,
};

describe('PrimeUI MoveDetails page', () => {
  describe('check move details page load', () => {
    it('displays payment requests information', async () => {
      usePrimeSimulatorGetMove.mockReturnValue(moveReturnValue);
      renderWithRouter(<MoveDetails />);

      const paymentRequestsHeading = screen.getByRole('heading', { name: 'Payment Requests', level: 2 });
      expect(paymentRequestsHeading).toBeInTheDocument();

      const uploadButton = screen.getByText(/Upload Document/, { selector: 'a.usa-button' });
      expect(uploadButton).toBeInTheDocument();
    });
  });
  describe('details button works', () => {
    it('can go to uploads page when payment request button clicked', async () => {
      usePrimeSimulatorGetMove.mockReturnValue(moveReturnValue);
      const { history } = renderWithRouter(<MoveDetails />);

      const uploadButton = screen.getByText(/Upload Document/, { selector: 'a.usa-button' });
      await userEvent.click(uploadButton);

      await waitFor(() => {
        expect(history.location.pathname).toEqual(`/payment-requests/${moveTaskOrder.paymentRequests[0].id}/upload`);
      });
    });
  });
});
