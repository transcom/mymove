import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import * as ReactRouterDom from 'react-router-dom';

import MoveDetails from './MoveDetails';

import { usePrimeSimulatorGetMove } from 'hooks/queries';
import { MockProviders } from 'testUtils';
import { completeCounseling, deleteShipment, downloadMoveOrder } from 'services/primeApi';
import { primeSimulatorRoutes } from 'constants/routes';
import { formatWeight } from 'utils/formatters';

// Labels
const payGradeLabelText = 'Pay Grade:';
const rankLabelText = 'Rank:';

const mockRequestedMoveCode = 'LN4T89';

jest.mock('hooks/queries', () => ({
  usePrimeSimulatorGetMove: jest.fn(),
}));

jest.mock('services/primeApi', () => ({
  completeCounseling: jest.fn(),
  deleteShipment: jest.fn(),
  downloadMoveOrder: jest.fn(),
}));

const mockRankValue = 'E-5';
const mockPayGradeValue = 'SGT';
const undefinedValue = 'undefined';

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
      shipmentType: 'HHG_INTO_NTS',
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
    {
      id: '5',
      shipmentType: 'HHG',
      requestedPickupDate: '2024-12-27',
      pickupAddress: { streetAddress1: '100 1st Avenue', city: 'New York', state: 'NY', postalCode: '10001' },
      destinationAddress: {
        streetAddress1: '123 East',
        streetAddress2: 'Apt 215H',
        city: 'Juneau',
        state: 'AK',
        postalCode: '99801',
      },
    },
    {
      id: '6',
      shipmentType: 'HHG',
      requestedPickupDate: '2024-12-27',
      pickupAddress: { streetAddress1: '123 East', city: 'Juneau', state: 'NY', postalCode: '99801' },
      destinationAddress: {
        streetAddress1: '100 1st Avenue',
        streetAddress2: 'Apt 215H',
        city: 'New York',
        state: 'NY',
        postalCode: '10001',
      },
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
    {
      reServiceCode: 'PODFSC',
      eTag: 'MjAyMy0xMS0yOVQxNToyMjoxMy45NjAwMTha',
      id: 'serviceItemPOEFSC',
      moveTaskOrderID: 'aa8dfe13-266a-4956-ac60-01c2355c06d3',
      mtoShipmentID: '6',
      reServiceName: 'International POD fuel surcharge',
      status: 'APPROVED',
    },
    {
      reServiceCode: 'POEFSC',
      eTag: 'MjAyMy0xMS0yOVQxNToyMjoxMy45NjAwMTha',
      id: 'serviceItemPOEFSC',
      moveTaskOrderID: 'aa8dfe13-266a-4956-ac60-01c2355c06d3',
      mtoShipmentID: '5',
      reServiceName: 'International POE fuel surcharge',
      status: 'APPROVED',
    },
  ],
  order: {
    entitlement: {
      gunSafe: true,
      weightRestriction: 500,
      ubWeightRestriction: 350,
    },
    rank: mockRankValue,
    grade: mockPayGradeValue,
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

const PayGradeTestEvaluate = (expectedValue) => {
  const payGradeLabel = screen.getByText(payGradeLabelText);
  const payGradeValue = payGradeLabel.nextElementSibling;
  expect(payGradeLabel).toBeInTheDocument();
  expect(payGradeValue).toBeInTheDocument();
  expect(payGradeValue.textContent).toBe(expectedValue);
};

const RankTestEvaluate = (expectedValue) => {
  const rankLabel = screen.getByText(rankLabelText);
  const rankValue = rankLabel.nextElementSibling;
  expect(rankLabel).toBeInTheDocument();
  expect(rankValue).toBeInTheDocument();
  expect(rankValue.textContent).toBe(expectedValue);
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
    it('renders move and entitlement detais on load', async () => {
      usePrimeSimulatorGetMove.mockReturnValue(moveReturnValue);

      renderWithProviders(<MoveDetails />);

      await waitFor(() => {
        expect(screen.getByText(/Move Code/)).toBeInTheDocument();
        expect(screen.getByText(/Move Id/)).toBeInTheDocument();
        const gunSafe = screen.getByText('Gun Safe:');
        expect(gunSafe).toBeInTheDocument();
        expect(gunSafe.nextElementSibling.textContent).toBe('yes');
        PayGradeTestEvaluate(mockPayGradeValue);
        RankTestEvaluate(mockRankValue);
        const adminRestrictedWeight = screen.getByText('Admin Restricted Weight:');
        expect(adminRestrictedWeight).toBeInTheDocument();
        expect(adminRestrictedWeight.nextElementSibling.textContent).toBe(
          formatWeight(moveTaskOrder.order.entitlement.weightRestriction),
        );
        const adminRestrictedUBWeight = screen.getByText('Admin Restricted UB Weight:');
        expect(adminRestrictedUBWeight).toBeInTheDocument();
        expect(adminRestrictedUBWeight.nextElementSibling.textContent).toBe(
          formatWeight(moveTaskOrder.order.entitlement.ubWeightRestriction),
        );
      });
    });

    it('displays payment requests information', async () => {
      usePrimeSimulatorGetMove.mockReturnValue(moveReturnValue);
      renderWithProviders(<MoveDetails />);

      const paymentRequestsHeading = screen.getByRole('heading', { name: 'Payment Requests', level: 2 });
      expect(paymentRequestsHeading).toBeInTheDocument();

      const uploadButton = screen.getByRole('link', { name: 'Upload Document' });
      expect(uploadButton).toBeInTheDocument();
    });

    it('displays the move acknowledge button', async () => {
      usePrimeSimulatorGetMove.mockReturnValue(moveReturnValue);
      renderWithProviders(<MoveDetails />);

      const acknowledgeButton = screen.getByLabelText('Acknowledge Move');
      expect(acknowledgeButton).toBeInTheDocument();
      expect(acknowledgeButton).toHaveAttribute('href');
      expect(acknowledgeButton.textContent).toBe('Acknowledge Move');
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

    it('shows edit button next to the right service items', async () => {
      usePrimeSimulatorGetMove.mockReturnValue(moveReturnValue);

      renderWithProviders(<MoveDetails />);

      // Check for Edit buttons
      const editButtons = screen.getAllByRole('link', { name: 'Edit' });
      expect(editButtons).toHaveLength(5);
    });
  });
  describe('MoveDetails component - Rank and Pay Grade', () => {
    usePrimeSimulatorGetMove.mockReturnValue(moveReturnValue);
    let useParamsSpy;

    beforeEach(() => {
      useParamsSpy = jest.spyOn(ReactRouterDom, 'useParams').mockReturnValue({ moveCodeOrID: 'MOCK123' });
    });

    afterEach(() => {
      useParamsSpy.mockRestore();
    });

    it('renders Rank information correctly', () => {
      usePrimeSimulatorGetMove.mockReturnValue(moveReturnValue);
      renderWithProviders(<MoveDetails />);

      RankTestEvaluate(mockRankValue);
    });

    it('renders Pay Grade information correctly', () => {
      usePrimeSimulatorGetMove.mockReturnValue(moveReturnValue);
      renderWithProviders(<MoveDetails />);

      PayGradeTestEvaluate(mockPayGradeValue);
    });

    it('handles missing Rank information', () => {
      usePrimeSimulatorGetMove.mockReturnValue({
        moveTaskOrder: {
          order: {
            grade: mockPayGradeValue,
            entitlement: {
              gunSafe: true,
            },
          },
        },
        isLoading: false,
        isError: false,
      });

      renderWithProviders(<MoveDetails />);

      RankTestEvaluate(undefinedValue);
    });

    it('handles missing Pay Grade information', () => {
      usePrimeSimulatorGetMove.mockReturnValue({
        moveTaskOrder: {
          order: {
            rank: mockRankValue,
            entitlement: {
              gunSafe: true,
            },
          },
        },
        isLoading: false,
        isError: false,
      });

      renderWithProviders(<MoveDetails />);

      PayGradeTestEvaluate(undefinedValue);
    });
  });
});
