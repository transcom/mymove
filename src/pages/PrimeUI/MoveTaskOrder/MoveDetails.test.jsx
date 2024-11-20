import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import MoveDetails from './MoveDetails';

import { INACCESSIBLE_RETURN_VALUE } from 'utils/test/api';
import { usePrimeSimulatorGetMove } from 'hooks/queries';
import { MockProviders } from 'testUtils';
import { completeCounseling, deleteShipment, downloadMoveOrder } from 'services/primeApi';
import { primeSimulatorRoutes } from 'constants/routes';

const mockRequestedMoveCode = 'LN4T89';

jest.mock('hooks/queries', () => ({
  usePrimeSimulatorGetMove: jest.fn(),
}));

jest.mock('services/primeApi', () => ({
  completeCounseling: jest.fn(),
  deleteShipment: jest.fn(),
  downloadMoveOrder: jest.fn(),
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
        eTag: 'MjAyMi0wNS0yNFQyMTowNjozNS45MDEwMjNa',
        estimatedIncentive: 1000000,
        estimatedWeight: 4000,
        expectedDepartureDate: '2020-03-15',
        hasProGear: false,
        id: '5b21b808-6933-43ea-8f6f-02fc0a639835',
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
  mtoServiceItems: [
    {
      reServiceCode: 'DDDSIT',
      reason: null,
      sitCustomerContacted: '2023-04-15',
      sitDestinationFinalAddress: {
        city: 'Beverly Hills',
        country: 'US',
        eTag: 'MjAyMy0xMS0yOVQxNToyMjoxMy43MDg2Nzla',
        id: '20d6218a-3fbc-4dbc-8258-d4b3ee009657',
        postalCode: '90210',
        state: 'CA',
        streetAddress1: '123 Any Street',
        streetAddress2: 'P.O. Box 12345',
        streetAddress3: 'c/o Some Person',
      },
      sitEntryDate: '2020-04-15',
      sitRequestedDelivery: '2023-04-15',
      eTag: 'MjAyMy0xMS0yOVQxNToyMjoxMy45Mjk0NzNa',
      id: 'serviceItemDDDSIT',
      modelType: 'MTOServiceItemDestSIT',
      moveTaskOrderID: 'aa8dfe13-266a-4956-ac60-01c2355c06d3',
      mtoShipmentID: '2',
      reServiceName: 'Domestic destination SIT delivery',
      status: 'APPROVED',
    },
    {
      reServiceCode: 'DDFSIT',
      reason: null,
      sitDepartureDate: '2020-04-15',
      sitDestinationFinalAddress: {
        city: 'Beverly Hills',
        country: 'US',
        eTag: 'MjAyMy0xMS0yOVQxNToyMjoxMy43MDg2Nzla',
        id: '20d6218a-3fbc-4dbc-8258-d4b3ee009657',
        postalCode: '90210',
        state: 'CA',
        streetAddress1: '123 Any Street',
        streetAddress2: 'P.O. Box 12345',
        streetAddress3: 'c/o Some Person',
      },
      sitEntryDate: '2020-04-15',
      eTag: 'MjAyMy0xMS0yOVQxNToyMjoxMy45NjAwMTha',
      id: 'serviceItemDDFSIT',
      modelType: 'MTOServiceItemDestSIT',
      moveTaskOrderID: 'aa8dfe13-266a-4956-ac60-01c2355c06d3',
      mtoShipmentID: '2',
      reServiceName: 'Domestic destination 1st day SIT',
      status: 'APPROVED',
    },
    {
      reServiceCode: 'DDASIT',
      reason: null,
      sitDepartureDate: '2020-04-15',
      sitEntryDate: '2020-04-15',
      eTag: 'MjAyMy0xMS0yOVQxNToyMjoxMy45NjAwMTha',
      id: 'serviceItemDDASIT',
      modelType: 'MTOServiceItemDestSIT',
      moveTaskOrderID: 'aa8dfe13-266a-4956-ac60-01c2355c06d3',
      mtoShipmentID: '2',
      reServiceName: "Domestic destination add'l SIT",
      status: 'APPROVED',
    },
  ],
  order: {
    entitlement: {
      gunSafe: true,
    },
  },
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

const renderWithProviders = (component) => {
  render(
    <MockProviders path={primeSimulatorRoutes.VIEW_MOVE_PATH} params={{ moveCodeOrID: mockRequestedMoveCode }}>
      {component}
    </MockProviders>,
  );
};
describe('PrimeUI MoveDetails page', () => {
  describe('check move details page load', () => {
    it('displays payment requests information', async () => {
      usePrimeSimulatorGetMove.mockReturnValue(moveReturnValue);
      renderWithProviders(<MoveDetails />);

      const paymentRequestsHeading = screen.getByRole('heading', { name: 'Payment Requests', level: 2 });
      expect(paymentRequestsHeading).toBeInTheDocument();

      const uploadButton = screen.getByRole('link', { name: 'Upload Document' });
      expect(uploadButton).toBeInTheDocument();
    });

    it('counseling ready to be completed', async () => {
      usePrimeSimulatorGetMove.mockReturnValue(moveReturnValue);
      renderWithProviders(<MoveDetails />);

      const completeCounselingButton = screen.getByText(/Complete Counseling/, { selector: 'button' });
      expect(completeCounselingButton).toBeInTheDocument();

      const field = screen.queryByText('Prime Counseling Completed At:');
      expect(field).not.toBeInTheDocument();
    });

    it('counseling already completed', async () => {
      usePrimeSimulatorGetMove.mockReturnValue(moveCounselingCompletedReturnValue);
      renderWithProviders(<MoveDetails />);

      const completeCounselingButton = screen.queryByText(/Complete Counseling/, { selector: 'button' });
      expect(completeCounselingButton).not.toBeInTheDocument();

      const field = screen.getByText('Prime Counseling Completed At:');
      expect(field).toBeInTheDocument();
      expect(field.nextElementSibling.textContent).toBe(moveTaskOrderCounselingCompleted.primeCounselingCompletedAt);
    });

    it('success when completing counseling', async () => {
      usePrimeSimulatorGetMove.mockReturnValue(moveReturnValue);
      renderWithProviders(<MoveDetails />);

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

      renderWithProviders(<MoveDetails />);

      const completeCounselingButton = screen.getByText(/Complete Counseling/, { selector: 'button' });
      await userEvent.click(completeCounselingButton);

      await waitFor(() => {
        expect(screen.getByText(/Error title/)).toBeInTheDocument();
        expect(screen.getByText('Error detail')).toBeInTheDocument();
      });
    });

    it('success when deleting PPM', async () => {
      usePrimeSimulatorGetMove.mockReturnValue(moveReturnValue);
      renderWithProviders(<MoveDetails />);

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

      renderWithProviders(<MoveDetails />);

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

    it('error when download move orders', async () => {
      usePrimeSimulatorGetMove.mockReturnValue(moveReturnValue);
      downloadMoveOrder.mockRejectedValue({
        response: { body: { title: 'Error title', detail: 'Error detail' } },
      });

      renderWithProviders(<MoveDetails />);

      const downloadMoveOrderButton = screen.getByText(/Download Move Orders/, { selector: 'button' });
      expect(downloadMoveOrderButton).toBeInTheDocument();
      await userEvent.click(downloadMoveOrderButton);

      await waitFor(() => {
        expect(screen.getByText(/Error title/)).toBeInTheDocument();
        expect(screen.getByText('Error detail')).toBeInTheDocument();
      });
    });

    it('success when downloading move orders', async () => {
      global.URL.createObjectURL = jest.fn();
      const mockResponse = {
        ok: true,
        headers: {
          'content-disposition': 'filename="test.pdf"',
        },
        status: 200,
        data: null,
      };
      usePrimeSimulatorGetMove.mockReturnValue(moveReturnValue);

      downloadMoveOrder.mockReturnValue(mockResponse);
      renderWithProviders(<MoveDetails />);

      const downloadMoveOrderButton = screen.getByText(/Download Move Orders/, { selector: 'button' });
      expect(downloadMoveOrderButton).toBeInTheDocument();

      jest.spyOn(document.body, 'appendChild');
      jest.spyOn(document, 'createElement');

      await userEvent.click(downloadMoveOrderButton);

      // verify hyperlink was created
      expect(document.createElement).toBeCalledWith('a');

      // verify hypelink element was created with correct
      // default file name from content-disposition
      expect(document.body.appendChild).toBeCalledWith(
        expect.objectContaining({
          download: 'test.pdf',
        }),
      );
    });

    it('shows edit button next to the right destination SIT service items', async () => {
      usePrimeSimulatorGetMove.mockReturnValue(moveReturnValue);

      renderWithProviders(<MoveDetails />);

      // Check for Edit buttons
      const editButtons = screen.getAllByRole('link', { name: 'Edit' });
      expect(editButtons).toHaveLength(3);
    });

    it("shows inaccessible error page when Prime user tries to access a safety move they don't have privileges to", async () => {
      usePrimeSimulatorGetMove.mockReturnValue(INACCESSIBLE_RETURN_VALUE);

      renderWithProviders(<MoveDetails />);

      const errorMessage = screen.getByText(/Page is not accessible./);
      expect(errorMessage).toBeInTheDocument();
    });
  });
});
