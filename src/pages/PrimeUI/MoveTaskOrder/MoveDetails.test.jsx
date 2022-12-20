import React from 'react';
import { screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import MoveDetails from './MoveDetails';

import { usePrimeSimulatorGetMove } from 'hooks/queries';
import { MockProviders, renderWithRouter } from 'testUtils';
import { completeCounseling, deleteShipment } from 'services/primeApi';

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
  deleteShipment: jest.fn(),
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
    {
      id: '4',
      approvedDate: '2022-05-24',
      createdAt: '2022-05-24T21:06:35.888Z',
      eTag: 'MjAyMi0wNS0yNFQyMTowNzoyMS4wNjc0MzJa',
      moveTaskOrderID: '9c7b255c-2981-4bf8-839f-61c7458e2b4d',
      ppmShipment: {
        advance: 598700,
        advanceRequested: true,
        createdAt: '2022-05-24T21:06:35.901Z',
        destinationPostalCode: '30813',
        eTag: 'MjAyMi0wNS0yNFQyMTowNjozNS45MDEwMjNa',
        estimatedIncentive: 1000000,
        estimatedWeight: 4000,
        expectedDepartureDate: '2020-03-15',
        hasProGear: false,
        id: '5b21b808-6933-43ea-8f6f-02fc0a639835',
        pickupPostalCode: '90210',
        shipmentId: '88ececed-eaf1-42e2-b060-cd90d11ad080',
        status: 'SUBMITTED',
        submittedAt: '2022-05-24T21:06:35.890Z',
        updatedAt: '2022-05-24T21:06:35.901Z',
      },
      shipmentType: 'PPM',
      status: 'APPROVED',
      updatedAt: '2022-05-24T21:07:21.067Z',
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

const moveTaskOrderCounselingCompleted = {
  ...moveTaskOrder,
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
      usePrimeSimulatorGetMove.mockReturnValue(moveReturnValue);
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
      usePrimeSimulatorGetMove.mockReturnValue(moveReturnValue);
      renderWithRouter(
        <MockProviders>
          <MoveDetails />
        </MockProviders>,
      );

      const completeCounselingButton = screen.getByText(/Complete Counseling/, { selector: 'button' });
      expect(completeCounselingButton).toBeInTheDocument();
      await userEvent.click(completeCounselingButton);

      await waitFor(() => {
        expect(screen.getByText('Successfully completed counseling')).toBeInTheDocument();
      });
    });

    it('error when completing counseling', async () => {
      usePrimeSimulatorGetMove.mockReturnValue(moveReturnValue);
      completeCounseling.mockRejectedValue({
        response: { body: { title: 'Error title', detail: 'Error detail' } },
      });

      renderWithRouter(
        <MockProviders>
          <MoveDetails />
        </MockProviders>,
      );

      const completeCounselingButton = screen.getByText(/Complete Counseling/, { selector: 'button' });
      await userEvent.click(completeCounselingButton);

      await waitFor(() => {
        expect(screen.getByText(/Error title/)).toBeInTheDocument();
        expect(screen.getByText('Error detail')).toBeInTheDocument();
      });
    });

    it('success when deleting PPM', async () => {
      usePrimeSimulatorGetMove.mockReturnValue(moveReturnValue);
      renderWithRouter(
        <MockProviders>
          <MoveDetails />
        </MockProviders>,
      );

      const deleteShipmentButton = screen.getByText(/Delete Shipment/, { selector: 'button' });
      expect(deleteShipmentButton).toBeInTheDocument();
      await userEvent.click(deleteShipmentButton);

      const modalDeleteButton = screen.getByText('Delete shipment', { selector: 'button.usa-button--destructive' });
      await userEvent.click(modalDeleteButton);

      await waitFor(() => {
        expect(screen.getByText('Successfully deleted shipment')).toBeInTheDocument();
      });
    });

    it('error when deleting PPM', async () => {
      usePrimeSimulatorGetMove.mockReturnValue(moveReturnValue);
      deleteShipment.mockRejectedValue({
        response: { body: { title: 'Error title', detail: 'Error detail' } },
      });

      renderWithRouter(
        <MockProviders>
          <MoveDetails />
        </MockProviders>,
      );

      const deleteShipmentButton = screen.getByText(/Delete Shipment/, { selector: 'button' });
      expect(deleteShipmentButton).toBeInTheDocument();
      await userEvent.click(deleteShipmentButton);

      const modalDeleteButton = screen.getByText('Delete shipment', { selector: 'button.usa-button--destructive' });
      await userEvent.click(modalDeleteButton);

      await waitFor(() => {
        expect(screen.getByText(/Error title/)).toBeInTheDocument();
        expect(screen.getByText('Error detail')).toBeInTheDocument();
      });
    });
  });
});
