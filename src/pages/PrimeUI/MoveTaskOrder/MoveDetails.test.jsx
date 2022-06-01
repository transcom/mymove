import React from 'react';
import { screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import MoveDetails from './MoveDetails';

import { usePrimeSimulatorGetMove } from 'hooks/queries';
import { MockProviders, renderWithRouter } from 'testUtils';
import { completeCounseling } from 'services/primeApi';

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
  completeCounseling: jest.fn(),
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

const moveTaskOrderCounselingReady = {
  id: '1',
  moveCode: mockRequestedMoveCode,
  shipmentType: 'PPM',
};

const moveTaskOrderCounselingReadyReturnValue = {
  moveTaskOrder: moveTaskOrderCounselingReady,
  isLoading: false,
  isError: false,
};

const moveTaskOrderCounselingCompleted = {
  ...moveTaskOrderCounselingReady,
  primeCounselingCompletedAt: '2022-05-24T21:06:35.890Z',
};

const moveCounselingCompletedReturnValue = {
  moveTaskOrder: moveTaskOrderCounselingCompleted,
  isLoading: false,
  isError: false,
};

describe('PrimeUI MoveDetails page', () => {
  describe('check move details page load', () => {
    it('displays payment requests information', async () => {
      usePrimeSimulatorGetMove.mockReturnValue(moveReturnValue);
      renderWithRouter(
        <MockProviders>
          <MoveDetails />
        </MockProviders>,
      );

      const paymentRequestsHeading = screen.getByRole('heading', { name: 'Payment Requests', level: 2 });
      expect(paymentRequestsHeading).toBeInTheDocument();

      const uploadButton = screen.getByText(/Upload Document/, { selector: 'a.usa-button' });
      expect(uploadButton).toBeInTheDocument();
    });

    it('counseling ready to be completed', async () => {
      usePrimeSimulatorGetMove.mockReturnValue(moveTaskOrderCounselingReadyReturnValue);
      renderWithRouter(
        <MockProviders>
          <MoveDetails />
        </MockProviders>,
      );

      const completeCounselingButton = screen.getByText(/Complete Counseling/, { selector: 'button' });
      expect(completeCounselingButton).toBeInTheDocument();

      const field = screen.queryByText('Prime Counseling Completed At:');
      expect(field).not.toBeInTheDocument();
    });

    it('counseling already completed', async () => {
      usePrimeSimulatorGetMove.mockReturnValue(moveCounselingCompletedReturnValue);
      renderWithRouter(
        <MockProviders>
          <MoveDetails />
        </MockProviders>,
      );

      const completeCounselingButton = screen.queryByText(/Complete Counseling/, { selector: 'button' });
      expect(completeCounselingButton).not.toBeInTheDocument();

      const field = screen.getByText('Prime Counseling Completed At:');
      expect(field).toBeInTheDocument();
      expect(field.nextElementSibling.textContent).toBe(moveTaskOrderCounselingCompleted.primeCounselingCompletedAt);
    });

    it('success when completing counseling', async () => {
      usePrimeSimulatorGetMove.mockReturnValue(moveTaskOrderCounselingReadyReturnValue);
      renderWithRouter(
        <MockProviders>
          <MoveDetails />
        </MockProviders>,
      );

      const completeCounselingButton = screen.getByText(/Complete Counseling/, { selector: 'button' });
      expect(completeCounselingButton).toBeInTheDocument();
      userEvent.click(completeCounselingButton);

      await waitFor(() => {
        expect(screen.getByText('Successfully completed counseling')).toBeInTheDocument();
      });
    });

    it('error when completing counseling', async () => {
      usePrimeSimulatorGetMove.mockReturnValue(moveTaskOrderCounselingReadyReturnValue);
      completeCounseling.mockRejectedValue({
        response: { body: { title: 'Error title', detail: 'Error detail' } },
      });

      renderWithRouter(
        <MockProviders>
          <MoveDetails />
        </MockProviders>,
      );

      const completeCounselingButton = screen.getByText(/Complete Counseling/, { selector: 'button' });
      userEvent.click(completeCounselingButton);

      await waitFor(() => {
        expect(screen.getByText(/Error title/)).toBeInTheDocument();
        expect(screen.getByText('Error detail')).toBeInTheDocument();
      });
    });
  });
});
