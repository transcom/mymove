import React from 'react';
import { render, waitFor, screen, fireEvent } from '@testing-library/react';
import { act } from 'react-dom/test-utils';
import userEvent from '@testing-library/user-event';

import ReviewWeightTicket from './ReviewWeightTicket';

import { MockProviders } from 'testUtils';

beforeEach(() => {
  jest.clearAllMocks();
});

jest.setTimeout(60000);

const mockCallback = jest.fn();

const defaultProps = {
  order: {
    entitlement: {
      totalWeight: 2000,
    },
  },
  ppmShipmentInfo: {
    id: '32ecb311-edbe-4fd4-96ee-bd693113f3f3',
    expectedDepartureDate: '2022-12-02',
    actualMoveDate: '2022-12-06',
    actualPickupPostalCode: '90210',
    actualDestinationPostalCode: '94611',
    miles: 300,
    estimatedWeight: 3000,
    actualWeight: 3500,
  },
  tripNumber: 1,
  showAllFields: false,
  ppmNumber: 1,
  setCurrentMtoShipments: mockCallback,
};

const baseWeightTicketProps = {
  id: '32ecb311-edbe-4fd4-96ee-bd693113f3f3',
  ppmShipmentId: '343bb456-63af-4f76-89bd-7403094a5c4d',
  vehicleDescription: 'Kia Forte',
  emptyWeight: 400,
  fullWeight: 1200,
};

const missingWeightTicketProps = {
  weightTicket: {
    ...baseWeightTicketProps,
    ownsTrailer: false,
    missingEmptyWeightTicket: true,
    missingFullWeightTicket: true,
  },
};

const weightTicketRequiredProps = {
  weightTicket: {
    ...baseWeightTicketProps,
    ownsTrailer: false,
  },
};

const ownsTrailerProps = {
  weightTicket: {
    ...baseWeightTicketProps,
    ownsTrailer: true,
  },
};

const claimableTrailerProps = {
  weightTicket: {
    ...baseWeightTicketProps,
    ownsTrailer: true,
    trailerMeetsCriteria: true,
  },
};
const fullShipmentProps = {
  weightTicket: {
    adjustedNetWeight: null,
    createdAt: '2023-12-20T17:33:48.215Z',
    eTag: 'MjAyMy0xMi0yMFQxNzozMzo1MC4xMDIwNzJa',
    emptyDocument: {
      id: '5d86029c-8b02-4ceb-917f-2122b92a9582',
      service_member_id: 'f9ee61ef-d1af-407a-a970-8a83bec0f5f5',
      uploads: [
        {
          bytes: 2202009,
          contentType: 'application/pdf',
          createdAt: '2023-12-20T17:33:48.209Z',
          filename: 'testFile.pdf',
          id: 'b8425855-7472-408a-9150-5be0b6673200',
          status: 'PROCESSING',
          updatedAt: '2023-12-20T17:33:48.209Z',
          url: '/storage/USER/uploads/b8425855-7472-408a-9150-5be0b6673200?contentType=application%2Fpdf',
        },
      ],
    },
    emptyDocumentId: '5d86029c-8b02-4ceb-917f-2122b92a9582',
    emptyWeight: 1000,
    fullDocument: {
      id: '1a0000fa-b882-4383-a9c1-9fa8b8ba57d9',
      service_member_id: 'f9ee61ef-d1af-407a-a970-8a83bec0f5f5',
      uploads: [
        {
          bytes: 2202009,
          contentType: 'application/pdf',
          createdAt: '2023-12-20T17:33:48.211Z',
          filename: 'testFile.pdf',
          id: 'df9f7a51-1282-4929-8278-5d9249662860',
          status: 'PROCESSING',
          updatedAt: '2023-12-20T17:33:48.211Z',
          url: '/storage/USER/uploads/df9f7a51-1282-4929-8278-5d9249662860?contentType=application%2Fpdf',
        },
      ],
    },
    fullDocumentId: '1a0000fa-b882-4383-a9c1-9fa8b8ba57d9',
    fullWeight: 8000,
    id: '4702f07c-486f-4420-b5b6-3ebb843e2277',
    missingEmptyWeightTicket: false,
    missingFullWeightTicket: false,
    netWeightRemarks: null,
    ownsTrailer: false,
    ppmShipmentId: '8ac2411e-d3eb-43e1-b98f-2476cfe8c02e',
    proofOfTrailerOwnershipDocument: {
      id: '5ae4d700-4344-4648-8ab2-a424e3cbcec6',
      service_member_id: 'f9ee61ef-d1af-407a-a970-8a83bec0f5f5',
      uploads: [
        {
          bytes: 2202009,
          contentType: 'application/pdf',
          createdAt: '2023-12-20T17:33:48.213Z',
          filename: 'testFile.pdf',
          id: 'e56bc856-73a8-4006-b5b2-045142c1c71d',
          status: 'PROCESSING',
          updatedAt: '2023-12-20T17:33:48.213Z',
          url: '/storage/USER/uploads/e56bc856-73a8-4006-b5b2-045142c1c71d?contentType=application%2Fpdf',
        },
      ],
    },
    proofOfTrailerOwnershipDocumentId: '5ae4d700-4344-4648-8ab2-a424e3cbcec6',
    reason: null,
    status: 'APPROVED',
    trailerMeetsCriteria: false,
    updatedAt: '2023-12-20T17:33:50.102Z',
    vehicleDescription: '2022 Honda CR-V Hybrid',
  },
  mtoShipment: {
    approvedDate: '2020-03-20T00:00:00.000Z',
    calculatedBillableWeight: 980,
    createdAt: '2023-12-20T17:33:48.205Z',
    customerRemarks: 'Please treat gently',
    eTag: 'MjAyMy0xMi0yMFQxNzozMzo1MC4wNzk0MjJa',
    hasSecondaryDeliveryAddress: null,
    hasSecondaryPickupAddress: null,
    id: '0ec4a5d2-3cca-4d1f-a2db-91f83d22644a',
    moveTaskOrderID: '0536abcf-15c3-4935-bcd8-353ecb03a385',
    ppmShipment: {
      actualDestinationPostalCode: '30813',
      actualMoveDate: '2020-03-16',
      actualPickupPostalCode: '42444',
      advanceAmountReceived: 340000,
      advanceAmountRequested: 598700,
      advanceStatus: 'APPROVED',
      approvedAt: '2022-04-15T12:30:00.000Z',
      createdAt: '2023-12-20T17:33:48.206Z',
      eTag: 'MjAyMy0xMi0yMFQxNzozMzo1MC4xMDExNDFa',
      estimatedIncentive: 1000000,
      estimatedWeight: 4000,
      allowableWeight: 7000,
      expectedDepartureDate: '2020-03-15',
      finalIncentive: null,
      hasProGear: true,
      hasReceivedAdvance: true,
      hasRequestedAdvance: true,
      id: '8ac2411e-d3eb-43e1-b98f-2476cfe8c02e',
      movingExpenses: null,
      proGearWeight: 1987,
      proGearWeightTickets: null,
      reviewedAt: null,
      shipmentId: '0ec4a5d2-3cca-4d1f-a2db-91f83d22644a',
      sitEstimatedCost: null,
      sitEstimatedDepartureDate: null,
      sitEstimatedEntryDate: null,
      sitEstimatedWeight: null,
      sitExpected: false,
      spouseProGearWeight: 498,
      status: 'NEEDS_CLOSEOUT',
      submittedAt: null,
      updatedAt: '2023-12-20T17:33:50.101Z',
      w2Address: {
        city: 'Beverly Hills',
        country: 'US',
        eTag: 'MjAyMy0xMi0yMFQxNzozMzo1MC4wOTk5NDha',
        id: '7d972704-d448-4d9e-97d4-5f5585edf3bb',
        postalCode: '90210',
        state: 'CA',
        streetAddress1: '123 Any Street',
        streetAddress2: 'P.O. Box 12345',
        streetAddress3: 'c/o Some Person',
      },
      weightTickets: [
        {
          adjustedNetWeight: null,
          createdAt: '2023-12-20T17:33:48.215Z',
          eTag: 'MjAyMy0xMi0yMFQxNzozMzo1MC4xMDIwNzJa',
          emptyDocument: {
            id: '5d86029c-8b02-4ceb-917f-2122b92a9582',
            service_member_id: 'f9ee61ef-d1af-407a-a970-8a83bec0f5f5',
            uploads: [
              {
                bytes: 2202009,
                contentType: 'application/pdf',
                createdAt: '2023-12-20T17:33:48.209Z',
                filename: 'testFile.pdf',
                id: 'b8425855-7472-408a-9150-5be0b6673200',
                status: 'PROCESSING',
                updatedAt: '2023-12-20T17:33:48.209Z',
                url: '/storage/USER/uploads/b8425855-7472-408a-9150-5be0b6673200?contentType=application%2Fpdf',
              },
            ],
          },
          emptyDocumentId: '5d86029c-8b02-4ceb-917f-2122b92a9582',
          emptyWeight: 1000,
          fullDocument: {
            id: '1a0000fa-b882-4383-a9c1-9fa8b8ba57d9',
            service_member_id: 'f9ee61ef-d1af-407a-a970-8a83bec0f5f5',
            uploads: [
              {
                bytes: 2202009,
                contentType: 'application/pdf',
                createdAt: '2023-12-20T17:33:48.211Z',
                filename: 'testFile.pdf',
                id: 'df9f7a51-1282-4929-8278-5d9249662860',
                status: 'PROCESSING',
                updatedAt: '2023-12-20T17:33:48.211Z',
                url: '/storage/USER/uploads/df9f7a51-1282-4929-8278-5d9249662860?contentType=application%2Fpdf',
              },
            ],
          },
          fullDocumentId: '1a0000fa-b882-4383-a9c1-9fa8b8ba57d9',
          fullWeight: 8000,
          id: '4702f07c-486f-4420-b5b6-3ebb843e2277',
          missingEmptyWeightTicket: false,
          missingFullWeightTicket: false,
          netWeightRemarks: null,
          ownsTrailer: false,
          ppmShipmentId: '8ac2411e-d3eb-43e1-b98f-2476cfe8c02e',
          proofOfTrailerOwnershipDocument: {
            id: '5ae4d700-4344-4648-8ab2-a424e3cbcec6',
            service_member_id: 'f9ee61ef-d1af-407a-a970-8a83bec0f5f5',
            uploads: [
              {
                bytes: 2202009,
                contentType: 'application/pdf',
                createdAt: '2023-12-20T17:33:48.213Z',
                filename: 'testFile.pdf',
                id: 'e56bc856-73a8-4006-b5b2-045142c1c71d',
                status: 'PROCESSING',
                updatedAt: '2023-12-20T17:33:48.213Z',
                url: '/storage/USER/uploads/e56bc856-73a8-4006-b5b2-045142c1c71d?contentType=application%2Fpdf',
              },
            ],
          },
          proofOfTrailerOwnershipDocumentId: '5ae4d700-4344-4648-8ab2-a424e3cbcec6',
          reason: null,
          status: 'APPROVED',
          trailerMeetsCriteria: false,
          updatedAt: '2023-12-20T17:33:50.102Z',
          vehicleDescription: '2022 Honda CR-V Hybrid',
        },
      ],
      reviewShipmentWeightsURL:
        '/counseling/moves/RVDGWV/shipments/0ec4a5d2-3cca-4d1f-a2db-91f83d22644a/document-review',
    },
    primeActualWeight: 980,
    shipmentType: 'PPM',
    status: 'APPROVED',
    updatedAt: '2023-12-20T17:33:50.079Z',
  },
  currentMtoShipments: [
    {
      approvedDate: '2020-03-20T00:00:00.000Z',
      calculatedBillableWeight: 980,
      createdAt: '2023-12-20T17:33:48.205Z',
      customerRemarks: 'Please treat gently',
      eTag: 'MjAyMy0xMi0yMFQxNzozMzo1MC4wNzk0MjJa',
      hasSecondaryDeliveryAddress: null,
      hasSecondaryPickupAddress: null,
      id: '0ec4a5d2-3cca-4d1f-a2db-91f83d22644a',
      moveTaskOrderID: '0536abcf-15c3-4935-bcd8-353ecb03a385',
      ppmShipment: {
        actualDestinationPostalCode: '30813',
        actualMoveDate: '2020-03-16',
        actualPickupPostalCode: '42444',
        advanceAmountReceived: 340000,
        advanceAmountRequested: 598700,
        advanceStatus: 'APPROVED',
        approvedAt: '2022-04-15T12:30:00.000Z',
        createdAt: '2023-12-20T17:33:48.206Z',
        eTag: 'MjAyMy0xMi0yMFQxNzozMzo1MC4xMDExNDFa',
        estimatedIncentive: 1000000,
        estimatedWeight: 4000,
        allowableWeight: 7000,
        expectedDepartureDate: '2020-03-15',
        finalIncentive: null,
        hasProGear: true,
        hasReceivedAdvance: true,
        hasRequestedAdvance: true,
        id: '8ac2411e-d3eb-43e1-b98f-2476cfe8c02e',
        movingExpenses: null,
        proGearWeight: 1987,
        proGearWeightTickets: null,
        reviewedAt: null,
        shipmentId: '0ec4a5d2-3cca-4d1f-a2db-91f83d22644a',
        sitEstimatedCost: null,
        sitEstimatedDepartureDate: null,
        sitEstimatedEntryDate: null,
        sitEstimatedWeight: null,
        sitExpected: false,
        spouseProGearWeight: 498,
        status: 'NEEDS_CLOSEOUT',
        submittedAt: null,
        updatedAt: '2023-12-20T17:33:50.101Z',
        w2Address: {
          city: 'Beverly Hills',
          country: 'US',
          eTag: 'MjAyMy0xMi0yMFQxNzozMzo1MC4wOTk5NDha',
          id: '7d972704-d448-4d9e-97d4-5f5585edf3bb',
          postalCode: '90210',
          state: 'CA',
          streetAddress1: '123 Any Street',
          streetAddress2: 'P.O. Box 12345',
          streetAddress3: 'c/o Some Person',
        },
        weightTickets: [
          {
            adjustedNetWeight: null,
            createdAt: '2023-12-20T17:33:48.215Z',
            eTag: 'MjAyMy0xMi0yMFQxNzozMzo1MC4xMDIwNzJa',
            emptyDocument: {
              id: '5d86029c-8b02-4ceb-917f-2122b92a9582',
              service_member_id: 'f9ee61ef-d1af-407a-a970-8a83bec0f5f5',
              uploads: [
                {
                  bytes: 2202009,
                  contentType: 'application/pdf',
                  createdAt: '2023-12-20T17:33:48.209Z',
                  filename: 'testFile.pdf',
                  id: 'b8425855-7472-408a-9150-5be0b6673200',
                  status: 'PROCESSING',
                  updatedAt: '2023-12-20T17:33:48.209Z',
                  url: '/storage/USER/uploads/b8425855-7472-408a-9150-5be0b6673200?contentType=application%2Fpdf',
                },
              ],
            },
            emptyDocumentId: '5d86029c-8b02-4ceb-917f-2122b92a9582',
            emptyWeight: 1000,
            fullDocument: {
              id: '1a0000fa-b882-4383-a9c1-9fa8b8ba57d9',
              service_member_id: 'f9ee61ef-d1af-407a-a970-8a83bec0f5f5',
              uploads: [
                {
                  bytes: 2202009,
                  contentType: 'application/pdf',
                  createdAt: '2023-12-20T17:33:48.211Z',
                  filename: 'testFile.pdf',
                  id: 'df9f7a51-1282-4929-8278-5d9249662860',
                  status: 'PROCESSING',
                  updatedAt: '2023-12-20T17:33:48.211Z',
                  url: '/storage/USER/uploads/df9f7a51-1282-4929-8278-5d9249662860?contentType=application%2Fpdf',
                },
              ],
            },
            fullDocumentId: '1a0000fa-b882-4383-a9c1-9fa8b8ba57d9',
            fullWeight: 8000,
            id: '4702f07c-486f-4420-b5b6-3ebb843e2277',
            missingEmptyWeightTicket: false,
            missingFullWeightTicket: false,
            netWeightRemarks: null,
            ownsTrailer: false,
            ppmShipmentId: '8ac2411e-d3eb-43e1-b98f-2476cfe8c02e',
            proofOfTrailerOwnershipDocument: {
              id: '5ae4d700-4344-4648-8ab2-a424e3cbcec6',
              service_member_id: 'f9ee61ef-d1af-407a-a970-8a83bec0f5f5',
              uploads: [
                {
                  bytes: 2202009,
                  contentType: 'application/pdf',
                  createdAt: '2023-12-20T17:33:48.213Z',
                  filename: 'testFile.pdf',
                  id: 'e56bc856-73a8-4006-b5b2-045142c1c71d',
                  status: 'PROCESSING',
                  updatedAt: '2023-12-20T17:33:48.213Z',
                  url: '/storage/USER/uploads/e56bc856-73a8-4006-b5b2-045142c1c71d?contentType=application%2Fpdf',
                },
              ],
            },
            proofOfTrailerOwnershipDocumentId: '5ae4d700-4344-4648-8ab2-a424e3cbcec6',
            reason: null,
            status: 'APPROVED',
            trailerMeetsCriteria: false,
            updatedAt: '2023-12-20T17:33:50.102Z',
            vehicleDescription: '2022 Honda CR-V Hybrid',
          },
        ],
        reviewShipmentWeightsURL:
          '/counseling/moves/RVDGWV/shipments/0ec4a5d2-3cca-4d1f-a2db-91f83d22644a/document-review',
      },
      primeActualWeight: 980,
      shipmentType: 'PPM',
      status: 'APPROVED',
      updatedAt: '2023-12-20T17:33:50.079Z',
    },
  ],
};

jest.mock('pages/Office/PPM/ReviewDocuments/ReviewDocuments', () => {
  const ReviewDocumentsMock = () => <div />;
  return ReviewDocumentsMock;
});

describe('ReviewWeightTicket component', () => {
  describe('displays form', () => {
    it('renders blank form on load with defaults', async () => {
      render(
        <MockProviders>
          <ReviewWeightTicket {...defaultProps} {...weightTicketRequiredProps} />
        </MockProviders>,
      );

      await waitFor(() => {
        expect(screen.getByRole('heading', { level: 3, name: 'Trip 1' })).toBeInTheDocument();
      });

      expect(screen.getByText('Vehicle description')).toBeInTheDocument();
      expect(screen.getByLabelText('Full weight')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByText('Net weight')).toBeInTheDocument();
      expect(screen.getByText('Did they use a trailer they owned?')).toBeInTheDocument();
      expect(screen.getByLabelText('Yes')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByLabelText('No')).toBeInstanceOf(HTMLInputElement);
      expect(screen.queryByText("Is the trailer's weight claimable?")).not.toBeInTheDocument();

      expect(screen.getByRole('heading', { level: 3, name: 'Review trip 1' })).toBeInTheDocument();
      expect(screen.getByLabelText('Accept')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByLabelText('Reject')).toBeInstanceOf(HTMLInputElement);
    });

    it('populates edit form with existing weight ticket values', async () => {
      render(
        <MockProviders>
          <ReviewWeightTicket {...defaultProps} {...weightTicketRequiredProps} />
        </MockProviders>,
      );

      await waitFor(() => {
        expect(screen.getByText('Kia Forte')).toBeInTheDocument();
      });
      expect(screen.getByLabelText('Empty weight', { description: 'Weight tickets' })).toHaveDisplayValue('400');
      expect(screen.getByLabelText('Full weight', { description: 'Weight tickets' })).toHaveDisplayValue('1,200');
      expect(screen.getAllByText('800 lbs')).toHaveLength(2);
      expect(screen.getByLabelText('No')).toBeChecked();
    });

    it('populates edit form when weight ticket is missing', async () => {
      render(
        <MockProviders>
          <ReviewWeightTicket {...defaultProps} {...missingWeightTicketProps} />
        </MockProviders>,
      );
      await waitFor(() => {
        expect(screen.getByLabelText('Empty weight', { description: 'Vehicle weight' })).toBeInTheDocument();
      });
      expect(screen.getByLabelText('Full weight', { description: 'Constructed weight' })).toBeInTheDocument();
    });

    it('toggles the reason field when Reject is selected', async () => {
      render(
        <MockProviders>
          <ReviewWeightTicket {...defaultProps} {...weightTicketRequiredProps} />
        </MockProviders>,
      );
      await waitFor(() => {
        expect(screen.getByLabelText('Reject')).toBeInstanceOf(HTMLInputElement);
      });
      await act(async () => {
        await fireEvent.click(screen.getByLabelText('Reject'));
      });
      expect(screen.getByLabelText('Reason')).toBeInstanceOf(HTMLTextAreaElement);
      await act(async () => {
        await fireEvent.click(screen.getByLabelText('Accept'));
      });
      expect(screen.queryByLabelText('Reason')).not.toBeInTheDocument();
    });

    it('populates edit form with existing weight ticket values and modifies them to alter the net weight calculation', async () => {
      render(
        <MockProviders>
          <ReviewWeightTicket updateTotalWeight={mockCallback} {...defaultProps} {...fullShipmentProps} />
        </MockProviders>,
      );

      await waitFor(() => {
        expect(screen.getByText('2022 Honda CR-V Hybrid')).toBeInTheDocument();
      });
      const emptyWeightInput = screen.getByTestId('emptyWeight');
      const fullWeightInput = screen.getByTestId('fullWeight');
      const netWeightDisplay = screen.getByTestId('net-weight-display');
      expect(emptyWeightInput).toHaveDisplayValue('1,000');
      expect(fullWeightInput).toHaveDisplayValue('8,000');
      expect(netWeightDisplay).toHaveTextContent('7,000');
      expect(screen.getByLabelText('No')).toBeChecked();

      await act(async () => {
        await userEvent.clear(fullWeightInput);
        await userEvent.type(fullWeightInput, '10,000');
        fullWeightInput.blur();
      });

      await waitFor(() => {
        expect(netWeightDisplay).toHaveTextContent('9,000');
      });
    });

    it('notifies the user when a trailer is claimable, and disables approval', async () => {
      render(
        <MockProviders>
          <ReviewWeightTicket {...defaultProps} {...ownsTrailerProps} />
        </MockProviders>,
      );
      await waitFor(() => {
        expect(screen.queryByText("Is the trailer's weight claimable?")).toBeInTheDocument();
      });
      const claimableYesButton = screen.getAllByRole('radio', { name: 'Yes' })[1];
      await act(async () => {
        await fireEvent.click(claimableYesButton);
      });
      expect(screen.queryByText('Proof of ownership is needed to accept this item.')).toBeInTheDocument();
      expect(screen.getByLabelText('Accept')).not.toBeChecked();
      expect(screen.getByLabelText('Accept')).toHaveAttribute('disabled');
    });

    it('notifies the user when a trailer is claimable after toggling ownership', async () => {
      render(
        <MockProviders>
          <ReviewWeightTicket {...defaultProps} {...claimableTrailerProps} />
        </MockProviders>,
      );
      await waitFor(() => {
        expect(screen.queryByText("Is the trailer's weight claimable?")).toBeInTheDocument();
      });
      const ownedNoButton = screen.getAllByRole('radio', { name: 'No' })[0];
      const ownedYesButton = screen.getAllByRole('radio', { name: 'Yes' })[0];
      await act(async () => {
        await fireEvent.click(ownedNoButton);
        await fireEvent.click(ownedYesButton);
      });
      expect(screen.queryByText('Proof of ownership is needed to accept this item.')).not.toBeInTheDocument();
    });

    it('reenables approval after disabling it and updating weight claimable field', async () => {
      render(
        <MockProviders>
          <ReviewWeightTicket {...defaultProps} {...claimableTrailerProps} />
        </MockProviders>,
      );
      await waitFor(() => {
        expect(screen.queryByText("Is the trailer's weight claimable?")).toBeInTheDocument();
      });
      const claimableNoButton = screen.getAllByRole('radio', { name: 'No' })[1];
      await act(async () => {
        await fireEvent.click(claimableNoButton);
      });
      expect(screen.queryByText('Proof of ownership is needed to accept this item.')).not.toBeInTheDocument();
      expect(screen.getByLabelText('Accept')).not.toHaveAttribute('disabled');
    });
    describe('shows an error when submitting', () => {
      it('without a status selected', async () => {
        render(
          <MockProviders>
            <ReviewWeightTicket {...defaultProps} {...claimableTrailerProps} />
          </MockProviders>,
        );
        await waitFor(async () => {
          const form = screen.getByRole('form');
          expect(form).toBeInTheDocument();
          await fireEvent.submit(form);
          expect(screen.getByText('Reviewing this weight ticket is required'));
        });
      });
      it('with Rejected but no reason selected', async () => {
        render(
          <MockProviders>
            <ReviewWeightTicket {...defaultProps} {...claimableTrailerProps} />
          </MockProviders>,
        );
        await waitFor(async () => {
          const form = screen.getByRole('form');
          expect(form).toBeInTheDocument();
          const rejectionButton = screen.getByTestId('rejectRadio');
          expect(rejectionButton).toBeInTheDocument();
          await fireEvent.click(rejectionButton);
          await fireEvent.submit(form);
          expect(screen.getByText('Add a reason why this weight ticket is rejected'));
        });
      });
    });
  });

  describe('displays disabled read only form', () => {
    it('renders disabled blank form on load with defaults', async () => {
      render(
        <MockProviders>
          <ReviewWeightTicket {...defaultProps} {...weightTicketRequiredProps} readOnly />
        </MockProviders>,
      );

      await waitFor(() => {
        expect(screen.getByRole('heading', { level: 3, name: 'Trip 1' })).toBeInTheDocument();
      });

      expect(screen.getByText('Vehicle description')).toBeInTheDocument();
      expect(screen.getByLabelText('Full weight')).toBeDisabled();
      expect(screen.getByText('Did they use a trailer they owned?')).toBeInTheDocument();
      expect(screen.getByLabelText('Yes')).toBeDisabled();
      expect(screen.getByLabelText('No')).toBeDisabled();
      expect(screen.queryByText("Is the trailer's weight claimable?")).not.toBeInTheDocument();

      expect(screen.getByRole('heading', { level: 3, name: 'Review trip 1' })).toBeInTheDocument();
      expect(screen.getByLabelText('Accept')).toBeDisabled();
      expect(screen.getByLabelText('Reject')).toBeDisabled();
    });

    it('populates disabled edit form with existing weight ticket values', async () => {
      render(
        <MockProviders>
          <ReviewWeightTicket {...defaultProps} {...weightTicketRequiredProps} readOnly />
        </MockProviders>,
      );

      await waitFor(() => {
        expect(screen.getByText('Kia Forte')).toBeInTheDocument();
      });
      expect(screen.getByLabelText('Empty weight', { description: 'Weight tickets' })).toHaveDisplayValue('400');
      expect(screen.getByLabelText('Empty weight', { description: 'Weight tickets' })).toBeDisabled();
      expect(screen.getByLabelText('Full weight', { description: 'Weight tickets' })).toHaveDisplayValue('1,200');
      expect(screen.getByLabelText('Full weight', { description: 'Weight tickets' })).toBeDisabled();
      expect(screen.getByLabelText('No')).toBeChecked();
      expect(screen.getByLabelText('No')).toBeDisabled();
    });

    it('populates disabled edit form when weight ticket is missing', async () => {
      render(
        <MockProviders>
          <ReviewWeightTicket {...defaultProps} {...missingWeightTicketProps} readOnly />
        </MockProviders>,
      );
      await waitFor(() => {
        expect(screen.getByLabelText('Empty weight', { description: 'Vehicle weight' })).toBeDisabled();
      });
      expect(screen.getByLabelText('Full weight', { description: 'Constructed weight' })).toBeDisabled();
    });
  });
});
