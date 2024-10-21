import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import ReviewBillableWeight from './ReviewBillableWeight';

import { formatWeight, formatDateFromIso } from 'utils/formatters';
import { useMovePaymentRequestsQueries } from 'hooks/queries';
import { shipmentStatuses, ppmShipmentStatuses } from 'constants/shipments';
import { tioRoutes } from 'constants/routes';
import { MockProviders, ReactQueryWrapper } from 'testUtils';
import { createPPMShipmentWithFinalIncentive } from 'utils/test/factories/ppmShipment';
import createUpload from 'utils/test/factories/upload';

// Mock the document viewer since we're not really testing that aspect here.
// Document Viewer tests should be covered in the component itself.
jest.mock('components/DocumentViewer/DocumentViewer', () => {
  const MockDocumentViewer = () => <div>Document viewer text</div>;
  return MockDocumentViewer;
});

jest.mock('hooks/queries', () => ({
  useMovePaymentRequestsQueries: jest.fn(),
}));

const mockNavigate = jest.fn();
jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useNavigate: () => mockNavigate,
}));
const routingParams = { moveCode: 'testMoveCode' };

const mockOriginDutyLocation = {
  address: {
    city: 'Des Moines',
    country: 'US',
    eTag: 'MjAyMC0wOS0xNFQxNzo0MTozOC42OTg1OTha',
    id: '2e26b066-aaca-4563-b284-d7f3f978fb3c',
    postalCode: '50309',
    state: 'IA',
    streetAddress1: '987 Other Avenue',
    streetAddress2: 'P.O. Box 1234',
    streetAddress3: 'c/o Another Person',
  },
  address_id: '2e26b066-aaca-4563-b284-d7f3f978fb3c',
  eTag: 'MjAyMC0wOS0xNFQxNzo0MTozOC43MDcxOTVa',
  id: 'a3ec2bdd-aa0a-434a-ba58-34c85f047704',
  name: 'XBc1KNi3pA',
};

const mockDestinationDutyLocation = {
  address: {
    city: 'Augusta',
    country: 'United States',
    eTag: 'MjAyMC0wOS0xNFQxNzo0MDo0OC44OTM3MDVa',
    id: '5ac95be8-0230-47ea-90b4-b0f6f60de364',
    postalCode: '30813',
    state: 'GA',
    streetAddress1: 'Fort Gordon',
  },
  address_id: '5ac95be8-0230-47ea-90b4-b0f6f60de364',
  eTag: 'MjAyMC0wOS0xNFQxNzo0MDo0OC44OTM3MDVa',
  id: '2d5ada83-e09a-47f8-8de6-83ec51694a86',
  name: 'Fort Gordon',
};

const mockOrders = {
  1: {
    agency: 'ARMY',
    customerID: '6ac40a00-e762-4f5f-b08d-3ea72a8e4b63',
    date_issued: '2018-03-15',
    department_indicator: 'AIR_AND_SPACE_FORCE',
    destinationDutyLocation: mockDestinationDutyLocation,
    eTag: 'MjAyMC0wOS0xNFQxNzo0MTozOC43MTE0Nlo=',
    entitlement: {
      authorizedWeight: 5000,
      dependentsAuthorized: true,
      eTag: 'MjAyMC0wOS0xNFQxNzo0MTozOC42ODAwOVo=',
      id: '0dbc9029-dfc5-4368-bc6b-dfc95f5fe317',
      nonTemporaryStorage: true,
      privatelyOwnedVehicle: true,
      proGearWeight: 2000,
      proGearWeightSpouse: 500,
      requiredMedicalEquipmentWeight: 1000,
      organizationalClothingAndIndividualEquipment: true,
      storageInTransit: 2,
      totalDependents: 1,
      totalWeight: 5000,
    },
    first_name: 'Leo',
    grade: 'E_1',
    id: '1',
    last_name: 'Spacemen',
    moveTaskOrder: {},
    order_number: 'ORDER3',
    order_type: 'PERMANENT_CHANGE_OF_STATION',
    order_type_detail: 'HHG_PERMITTED',
    originDutyLocation: mockOriginDutyLocation,
    report_by_date: '2018-08-01',
    tac: 'F8E1',
    sac: 'E2P3',
  },
};

const mockMtoShipments = [
  {
    id: 1,
    status: shipmentStatuses.APPROVED,
    shipmentType: 'HHG',
    calculatedBillableWeight: 3000,
    billableWeightCap: 1000,
    primeEstimatedWeight: 1000,
    primeActualWeight: 300,
    reweigh: { verificationReason: 'reweigh required', requestedAt: '2021-09-01' },
    pickupAddress: { city: 'Las Vegas', state: 'NV', postalCode: '90210' },
    destinationAddress: { city: 'Miami', state: 'FL', postalCode: '33607' },
    actualPickupDate: '2021-08-31',
  },
  {
    id: 2,
    status: shipmentStatuses.APPROVED,
    shipmentType: 'HHG',
    calculatedBillableWeight: 2000,
    billableWeightCap: 2000,
    primeEstimatedWeight: 2000,
    primeActualWeight: 400,
    reweigh: { weight: 1000, verificationReason: 'reweigh required', requestedAt: '2021-09-01' },
    pickupAddress: { city: 'Las Vegas', state: 'NV', postalCode: '90210' },
    destinationAddress: { city: 'Miami', state: 'FL', postalCode: '33607' },
    actualPickupDate: '2021-08-31',
  },
  {
    id: 3,
    status: shipmentStatuses.DIVERSION_REQUESTED,
    shipmentType: 'HHG',
    calculatedBillableWeight: 3000,
    billableWeightCap: 3000,
    primeEstimatedWeight: 7000,
    primeActualWeight: 300,
    reweigh: { weight: 200, verificationReason: 'reweigh required', requestedAt: '2021-09-01' },
    pickupAddress: { city: 'Las Vegas', state: 'NV', postalCode: '90210' },
    destinationAddress: { city: 'Miami', state: 'FL', postalCode: '33607' },
    actualPickupDate: '2021-08-31',
  },
];

const mockMtoNTSReleaseShipments = [
  {
    id: 1,
    status: shipmentStatuses.APPROVED,
    shipmentType: 'HHG_OUTOF_NTS_DOMESTIC',
    calculatedBillableWeight: 3000,
    billableWeightCap: 1000,
    primeEstimatedWeight: 1000,
    primeActualWeight: 300,
    reweigh: { verificationReason: 'reweigh required', requestedAt: '2021-09-01' },
    pickupAddress: { city: 'Las Vegas', state: 'NV', postalCode: '90210' },
    destinationAddress: { city: 'Miami', state: 'FL', postalCode: '33607' },
    actualPickupDate: '2021-08-31',
  },
];

const mockHasAllInformationShipment = {
  id: 1,
  status: shipmentStatuses.DIVERSION_REQUESTED,
  shipmentType: 'HHG',
  calculatedBillableWeight: 3000,
  billableWeightCap: 3000,
  primeEstimatedWeight: 7000,
  primeActualWeight: 300,
  reweigh: { weight: 200, verificationReason: 'reweigh required', requestedAt: '2021-09-01' },
  pickupAddress: { city: 'Las Vegas', state: 'NV', postalCode: '90210' },
  destinationAddress: { city: 'Miami', state: 'FL', postalCode: '33607' },
  actualPickupDate: '2021-08-31',
};

const mockNoReweighWeightShipment = {
  id: 2,
  status: shipmentStatuses.DIVERSION_REQUESTED,
  shipmentType: 'HHG',
  calculatedBillableWeight: 3000,
  billableWeightCap: 3000,
  primeEstimatedWeight: 7000,
  primeActualWeight: 300,
  reweigh: { verificationReason: 'reweigh required', requestedAt: '2021-09-01' },
  pickupAddress: { city: 'Las Vegas', state: 'NV', postalCode: '90210' },
  destinationAddress: { city: 'Miami', state: 'FL', postalCode: '33607' },
  actualPickupDate: '2021-08-31',
};

const mockNoPrimeEstimatedWeightShipment = {
  id: 3,
  status: shipmentStatuses.DIVERSION_REQUESTED,
  shipmentType: 'HHG',
  calculatedBillableWeight: 3000,
  billableWeightCap: 3000,
  primeActualWeight: 300,
  reweigh: { weight: 200, verificationReason: 'reweigh required', requestedAt: '2021-09-01' },
  pickupAddress: { city: 'Las Vegas', state: 'NV', postalCode: '90210' },
  destinationAddress: { city: 'Miami', state: 'FL', postalCode: '33607' },
  actualPickupDate: '2021-08-31',
};

const mockDivertedMtoShipments = [
  {
    id: 1,
    status: shipmentStatuses.APPROVED,
    shipmentType: 'HHG',
    calculatedBillableWeight: 2000,
    billableWeightCap: 2000,
    primeEstimatedWeight: 1500,
    primeActualWeight: 1250,
    reweigh: {},
    pickupAddress: { city: 'Las Vegas', state: 'NV', postalCode: '88901' },
    destinationAddress: { city: 'Miami', state: 'FL', postalCode: '33607' },
    actualPickupDate: '2021-08-31',
    diversion: true,
  },
  {
    id: 2,
    status: shipmentStatuses.APPROVED,
    shipmentType: 'HHG',
    calculatedBillableWeight: 2000,
    billableWeightCap: 2000,
    primeEstimatedWeight: 1500,
    primeActualWeight: 1250,
    reweigh: {},
    pickupAddress: { city: 'Miami', state: 'FL', postalCode: '33101' },
    destinationAddress: { city: 'Portland', state: 'ME', postalCode: '04109' },
    actualPickupDate: '2021-09-01',
    diversion: true,
  },
  {
    id: 3,
    status: shipmentStatuses.APPROVED,
    shipmentType: 'HHG',
    calculatedBillableWeight: 2000,
    billableWeightCap: 2000,
    primeEstimatedWeight: 2000,
    primeActualWeight: 1800,
    reweigh: {},
    pickupAddress: { city: 'Las Vegas', state: 'NV', postalCode: '88901' },
    destinationAddress: { city: 'Portland', state: 'ME', postalCode: '04109' },
    actualPickupDate: '2021-08-31',
    diversion: false,
  },
];

const mtoShipment = createPPMShipmentWithFinalIncentive({
  ppmShipment: { status: ppmShipmentStatuses.NEEDS_CLOSEOUT },
  id: 'shipment123',
});

const weightTicketEmptyDocumentCreatedDate = new Date();
// The factory used above doesn't handle overrides for uploads correctly, so we need to do it manually.
const weightTicketEmptyDocumentUpload = createUpload({
  fileName: 'emptyWeightTicket.pdf',
  createdAtDate: weightTicketEmptyDocumentCreatedDate,
});

const weightTicketFullDocumentCreatedDate = new Date(weightTicketEmptyDocumentCreatedDate);
weightTicketFullDocumentCreatedDate.setDate(weightTicketFullDocumentCreatedDate.getDate() + 1);
const weightTicketFullDocumentUpload = createUpload(
  { fileName: 'fullWeightTicket.xls', createdAtDate: weightTicketFullDocumentCreatedDate },
  { contentType: 'application/vnd.ms-excel' },
);

const progearWeightTicketDocumentCreatedDate = new Date(weightTicketFullDocumentCreatedDate);
progearWeightTicketDocumentCreatedDate.setDate(progearWeightTicketDocumentCreatedDate.getDate() + 1);
const progearWeightTicketDocumentUpload = createUpload({
  fileName: 'progearWeightTicket.pdf',
  createdAtDate: progearWeightTicketDocumentCreatedDate,
});

const movingExpenseDocumentCreatedDate = new Date(progearWeightTicketDocumentCreatedDate);
movingExpenseDocumentCreatedDate.setDate(movingExpenseDocumentCreatedDate.getDate() + 1);
const movingExpenseDocumentUpload = createUpload(
  { fileName: 'movingExpense.jpg', createdAtDate: movingExpenseDocumentCreatedDate },
  { contentType: 'image/jpeg' },
);

mtoShipment.ppmShipment.weightTickets[0].emptyDocument.uploads = [weightTicketEmptyDocumentUpload];
mtoShipment.ppmShipment.weightTickets[0].fullDocument.uploads = [weightTicketFullDocumentUpload];
mtoShipment.ppmShipment.proGearWeightTickets[0].document.uploads = [progearWeightTicketDocumentUpload];
mtoShipment.ppmShipment.movingExpenses[0].document.uploads = [movingExpenseDocumentUpload];

const mockPaymentRequest = [
  {
    proofOfServiceDocs: [
      {
        uploads: [
          {
            id: 'z',
            filename: 'test.pdf',
            contentType: 'application/pdf',
            url: '/storage/user/1/uploads/2?contentType=application%2Fpdf',
          },
        ],
      },
    ],
  },
];

const move = {
  tioRemarks: 'the prime has already unloaded this move',
};

const useMovePaymentRequestsReturnValue = {
  order: mockOrders['1'],
  mtoShipments: mockMtoShipments,
  move,
  paymentRequests: mockPaymentRequest,
};

const useMovePaymentRequestsNTSReleaseReturnValue = {
  order: mockOrders['1'],
  mtoShipments: mockMtoNTSReleaseShipments,
  move,
  paymentRequests: mockPaymentRequest,
};

const useNonMaxBillableWeightExceededReturnValue = {
  order: mockOrders['1'],
  mtoShipments: [mockMtoShipments[0]],
  move,
  paymentRequests: mockPaymentRequest,
};

const useMissingShipmentWeightNoReweighReturnValue = {
  order: mockOrders['1'],
  mtoShipments: [mockNoReweighWeightShipment, mockHasAllInformationShipment],
  move,
  paymentRequests: mockPaymentRequest,
};

const useMissingShipmentWeightNoPrimeEstimatedWeightReturnValue = {
  order: mockOrders['1'],
  mtoShipments: [mockNoPrimeEstimatedWeightShipment, mockHasAllInformationShipment],
  move,
  paymentRequests: mockPaymentRequest,
};

const noAlertsReturnValue = {
  order: mockOrders['1'],
  mtoShipments: [mockHasAllInformationShipment],
  move,
  paymentRequests: mockPaymentRequest,
};

const useDivertedMovePaymentRequestsReturnValue = {
  order: mockOrders['1'],
  mtoShipments: mockDivertedMtoShipments,
  move,
  paymentRequests: mockPaymentRequest,
};

const useMovePaymentRequestQueriesReturnValueAllDocs = {
  paymentRequests: [],
  order: mockOrders['1'],
  mtoShipments: [mtoShipment],
  move,
  isLoading: false,
  isError: false,
};

const loadingReturnValue = {
  isLoading: true,
  isError: false,
  isSuccess: false,
};

const errorReturnValue = {
  isLoading: false,
  isError: true,
  isSuccess: false,
};

const renderWithProviders = (component) => {
  return render(
    <ReactQueryWrapper>
      <MockProviders path={tioRoutes.BASE_PAYMENT_REQUESTS_PATH} params={routingParams}>
        {component}
      </MockProviders>
    </ReactQueryWrapper>,
  );
};

describe('ReviewBillableWeight', () => {
  describe('check loading and error component states', () => {
    it('renders the loading placeholder when the query is still loading', async () => {
      useMovePaymentRequestsQueries.mockReturnValue(loadingReturnValue);

      renderWithProviders(<ReviewBillableWeight />);

      const h2 = await screen.getByRole('heading', { name: 'Loading, please wait...', level: 2 });
      expect(h2).toBeInTheDocument();
    });

    it('renders the Something Went Wrong component when the query errors', async () => {
      useMovePaymentRequestsQueries.mockReturnValue(errorReturnValue);

      renderWithProviders(<ReviewBillableWeight />);

      const errorMessage = await screen.getByText(/Something went wrong./);
      expect(errorMessage).toBeInTheDocument();
    });
  });

  describe('check that all the components render', () => {
    it('renders the component', () => {
      useMovePaymentRequestsQueries.mockReturnValue(useMovePaymentRequestsReturnValue);

      renderWithProviders(<ReviewBillableWeight />);
      expect(screen.getByText('Review weights')).toBeInTheDocument();
      expect(screen.getByText('Document viewer text')).toBeInTheDocument();
      expect(screen.getByText(move.tioRemarks)).toBeInTheDocument();
    });

    it('renders weight summary', () => {
      useMovePaymentRequestsQueries.mockReturnValue(useMovePaymentRequestsReturnValue);
      renderWithProviders(<ReviewBillableWeight />);
      expect(screen.getByTestId('maxBillableWeight').textContent).toBe(
        formatWeight(useMovePaymentRequestsReturnValue.order.entitlement.authorizedWeight),
      );
      expect(screen.getByTestId('weightAllowance').textContent).toBe(
        formatWeight(useMovePaymentRequestsReturnValue.order.entitlement.totalWeight),
      );
      expect(screen.getByTestId('weightRequested').textContent).toBe('900 lbs');
      expect(screen.getByTestId('totalBillableWeight').textContent).toBe('6,100 lbs');
    });

    it('renders 110% estimated weight and edit view', async () => {
      useMovePaymentRequestsQueries.mockReturnValue(useMovePaymentRequestsReturnValue);
      renderWithProviders(<ReviewBillableWeight />);

      await userEvent.click(screen.getByText('Edit'));

      expect(screen.getByTestId('maxWeight-estimatedWeight').textContent).toBe('11,000 lbs');
      expect(screen.getByText(move.tioRemarks)).toBeInTheDocument();
    });
  });

  describe('check the nagivation', () => {
    it('takes the user back to the payment requests page when x is clicked', async () => {
      renderWithProviders(<ReviewBillableWeight />);

      const xButton = screen.getByTestId('closeSidebar');

      await userEvent.click(xButton);

      await waitFor(() => {
        expect(mockNavigate).toHaveBeenCalledWith('../payment-requests', {
          state: {
            from: 'review-billable-weights',
          },
        });
      });
    });

    it('takes the user to review the shipment weights when the review weights button is clicked', async () => {
      useMovePaymentRequestsQueries.mockReturnValue(useMovePaymentRequestsReturnValue);

      renderWithProviders(<ReviewBillableWeight />);

      const reviewShipmentWeights = screen.getByRole('button', { name: 'Review shipment weights' });

      await userEvent.click(reviewShipmentWeights);

      await waitFor(() => {
        expect(screen.getByText('Review weights')).toBeInTheDocument();
        expect(screen.getByText('Shipment weights')).toBeInTheDocument();

        expect(screen.getByTestId('maxBillableWeight').textContent).toBe(
          formatWeight(useMovePaymentRequestsReturnValue.order.entitlement.authorizedWeight),
        );
        expect(screen.getByTestId('weightAllowance').textContent).toBe(
          formatWeight(useMovePaymentRequestsReturnValue.order.entitlement.totalWeight),
        );
        expect(screen.getByTestId('ShipmentContainer')).toBeInTheDocument();
        // shipment container labels
        expect(screen.getByText('Departed')).toBeInTheDocument();
        expect(screen.getByText('From')).toBeInTheDocument();
        expect(screen.getByText('To')).toBeInTheDocument();
        expect(screen.getByText('Estimated weight')).toBeInTheDocument();
        expect(screen.getByText('Original weight')).toBeInTheDocument();
        expect(screen.getByText('Reweigh weight')).toBeInTheDocument();
        expect(screen.getByText('Date reweigh requested')).toBeInTheDocument();
        expect(screen.getByText('Reweigh remarks')).toBeInTheDocument();
        expect(screen.getByText('Billable weight')).toBeInTheDocument();
        expect(screen.getByText('reweigh required')).toBeInTheDocument();
      });
    });

    it('takes the user to the next shipment when the Next Shipment button is clicked', async () => {
      useMovePaymentRequestsQueries.mockReturnValue(useMovePaymentRequestsReturnValue);

      renderWithProviders(<ReviewBillableWeight />);

      const reviewShipmentWeights = screen.getByRole('button', { name: 'Review shipment weights' });

      await userEvent.click(reviewShipmentWeights);

      expect(screen.getByText('Shipment 1 of 3')).toBeInTheDocument();
      expect(screen.getByTestId('estimatedWeight').textContent).toBe(
        formatWeight(mockMtoShipments[0].primeEstimatedWeight),
      );
      expect(screen.getByTestId('originalWeight').textContent).toBe(
        formatWeight(mockMtoShipments[0].primeActualWeight),
      );
      expect(screen.getByTestId('reweighWeight').textContent).toBe('Missing');
      expect(screen.getByTestId('dateReweighRequested').textContent).toBe(
        formatDateFromIso(mockMtoShipments[0].reweigh.requestedAt, 'DD MMM YYYY'),
      );
      expect(screen.getByTestId('reweighRemarks').textContent).toBe(mockMtoShipments[0].reweigh.verificationReason);

      const nextShipment = screen.getByRole('button', { name: 'Next Shipment' });
      await userEvent.click(nextShipment);
      await waitFor(() => {
        expect(screen.getByText('Shipment 2 of 3')).toBeInTheDocument();
      });
      expect(screen.getByTestId('estimatedWeight').textContent).toBe(
        formatWeight(mockMtoShipments[1].primeEstimatedWeight),
      );
      expect(screen.getByTestId('originalWeight').textContent).toBe(
        formatWeight(mockMtoShipments[1].primeActualWeight),
      );
      expect(screen.getByTestId('reweighWeight').textContent).toBe(formatWeight(mockMtoShipments[1].reweigh.weight));
      expect(screen.getByTestId('dateReweighRequested').textContent).toBe(
        formatDateFromIso(mockMtoShipments[1].reweigh.requestedAt, 'DD MMM YYYY'),
      );
      expect(screen.getByTestId('reweighRemarks').textContent).toBe(mockMtoShipments[1].reweigh.verificationReason);
      await userEvent.click(nextShipment);
      await waitFor(() => {
        expect(screen.getByText('Shipment 3 of 3')).toBeInTheDocument();
      });
      expect(screen.getByTestId('estimatedWeight').textContent).toBe(
        formatWeight(mockMtoShipments[2].primeEstimatedWeight),
      );
      expect(screen.getByTestId('originalWeight').textContent).toBe(
        formatWeight(mockMtoShipments[2].primeActualWeight),
      );
      expect(screen.getByTestId('reweighWeight').textContent).toBe(formatWeight(mockMtoShipments[2].reweigh.weight));
      expect(screen.getByTestId('dateReweighRequested').textContent).toBe(
        formatDateFromIso(mockMtoShipments[2].reweigh.requestedAt, 'DD MMM YYYY'),
      );
      expect(screen.getByTestId('reweighRemarks').textContent).toBe(mockMtoShipments[2].reweigh.verificationReason);
      expect(screen.queryByRole('button', { name: 'Next Shipment' })).not.toBeInTheDocument();
    });

    it('takes the user to the previous shipment when the Back button is clicked', async () => {
      useMovePaymentRequestsQueries.mockReturnValue(useMovePaymentRequestsReturnValue);

      renderWithProviders(<ReviewBillableWeight />);

      const reviewShipmentWeights = screen.getByRole('button', { name: 'Review shipment weights' });

      await userEvent.click(reviewShipmentWeights);

      const nextShipment = screen.getByRole('button', { name: 'Next Shipment' });
      await userEvent.click(nextShipment);
      await userEvent.click(nextShipment);
      await waitFor(() => {
        expect(screen.getByText('Shipment 3 of 3')).toBeInTheDocument();
      });
      expect(screen.getByTestId('estimatedWeight').textContent).toBe(
        formatWeight(mockMtoShipments[2].primeEstimatedWeight),
      );
      expect(screen.getByTestId('originalWeight').textContent).toBe(
        formatWeight(mockMtoShipments[2].primeActualWeight),
      );
      expect(screen.getByTestId('reweighWeight').textContent).toBe(formatWeight(mockMtoShipments[2].reweigh.weight));
      expect(screen.getByTestId('dateReweighRequested').textContent).toBe(
        formatDateFromIso(mockMtoShipments[2].reweigh.requestedAt, 'DD MMM YYYY'),
      );
      expect(screen.getByTestId('reweighRemarks').textContent).toBe(mockMtoShipments[2].reweigh.verificationReason);

      const back = screen.getByRole('button', { name: 'Back' });
      await userEvent.click(back);
      await waitFor(() => {
        expect(screen.getByText('Shipment 2 of 3')).toBeInTheDocument();
      });
      expect(screen.getByTestId('estimatedWeight').textContent).toBe(
        formatWeight(mockMtoShipments[1].primeEstimatedWeight),
      );
      expect(screen.getByTestId('originalWeight').textContent).toBe(
        formatWeight(mockMtoShipments[1].primeActualWeight),
      );
      expect(screen.getByTestId('reweighWeight').textContent).toBe(formatWeight(mockMtoShipments[1].reweigh.weight));
      expect(screen.getByTestId('dateReweighRequested').textContent).toBe(
        formatDateFromIso(mockMtoShipments[1].reweigh.requestedAt, 'DD MMM YYYY'),
      );
      expect(screen.getByTestId('reweighRemarks').textContent).toBe(mockMtoShipments[1].reweigh.verificationReason);

      await userEvent.click(back);
      await waitFor(() => {
        expect(screen.getByText('Shipment 1 of 3')).toBeInTheDocument();
      });
      expect(screen.getByTestId('estimatedWeight').textContent).toBe(
        formatWeight(mockMtoShipments[0].primeEstimatedWeight),
      );
      expect(screen.getByTestId('originalWeight').textContent).toBe(
        formatWeight(mockMtoShipments[0].primeActualWeight),
      );
      expect(screen.getByTestId('reweighWeight').textContent).toBe('Missing');
      expect(screen.getByTestId('dateReweighRequested').textContent).toBe(
        formatDateFromIso(mockMtoShipments[0].reweigh.requestedAt, 'DD MMM YYYY'),
      );
      expect(screen.getByTestId('reweighRemarks').textContent).toBe(mockMtoShipments[0].reweigh.verificationReason);

      await userEvent.click(back);
      await waitFor(() => {
        expect(screen.getByText('Edit max billable weight')).toBeInTheDocument();
      });
      expect(screen.queryByRole('button', { name: 'Next Shipment' })).not.toBeInTheDocument();
      expect(screen.queryByRole('button', { name: 'Back' })).not.toBeInTheDocument();
      expect(screen.queryByRole('button', { name: 'Review shipment weights' })).toBeInTheDocument();
    });
  });

  describe('check that the various alerts show up when expected', () => {
    describe('max billable weight alert', () => {
      it('does not render in edit view when billable weight is not exceeded', async () => {
        useMovePaymentRequestsQueries.mockReturnValue(useNonMaxBillableWeightExceededReturnValue);

        renderWithProviders(<ReviewBillableWeight />);

        await userEvent.click(screen.getByText('Edit'));
        await waitFor(() => {
          expect(screen.queryByTestId('maxBillableWeightAlert')).not.toBeInTheDocument();
        });
      });

      it('does not render in shipment view when billable weight is not exceeded', async () => {
        useMovePaymentRequestsQueries.mockReturnValue(useNonMaxBillableWeightExceededReturnValue);

        renderWithProviders(<ReviewBillableWeight />);

        await userEvent.click(screen.getByText('Edit'));
        await userEvent.click(screen.getByText('Review shipment weights'));
        expect(screen.queryByTestId('maxBillableWeightAlert')).not.toBeInTheDocument();
      });
    });

    describe('missing shipment weights may impact max billable weight', () => {
      it('renders when a shipment is missing a reweigh weight', () => {
        useMovePaymentRequestsQueries.mockReturnValue(useMissingShipmentWeightNoReweighReturnValue);

        renderWithProviders(<ReviewBillableWeight />);

        expect(screen.getByTestId('maxBillableWeightMissingShipmentWeightAlert')).toBeInTheDocument();
      });

      it('renders when a shipment is missing a prime estimated weight', () => {
        useMovePaymentRequestsQueries.mockReturnValue(useMissingShipmentWeightNoPrimeEstimatedWeightReturnValue);

        renderWithProviders(<ReviewBillableWeight />);

        expect(screen.getByTestId('maxBillableWeightMissingShipmentWeightAlert')).toBeInTheDocument();
      });

      it('does not render when none of the shipments are missing reweigh or prime estimated weight information', () => {
        useMovePaymentRequestsQueries.mockReturnValue(noAlertsReturnValue);

        renderWithProviders(<ReviewBillableWeight />);

        expect(screen.queryByTestId('maxBillableWeightMissingShipmentWeightAlert')).not.toBeInTheDocument();
      });
    });

    describe('shipment missing information', () => {
      it('renders the alert when the shipment is missing a reweigh weight', async () => {
        useMovePaymentRequestsQueries.mockReturnValue(useMissingShipmentWeightNoReweighReturnValue);

        renderWithProviders(<ReviewBillableWeight />);

        const reviewShipmentWeights = screen.getByRole('button', { name: 'Review shipment weights' });
        await userEvent.click(reviewShipmentWeights);

        expect(screen.getByTestId('shipmentMissingInformation')).toBeInTheDocument();
      });

      it('renders the alert when the shipment is missing a prime estimated weight', async () => {
        useMovePaymentRequestsQueries.mockReturnValue(useMissingShipmentWeightNoPrimeEstimatedWeightReturnValue);

        renderWithProviders(<ReviewBillableWeight />);

        const reviewShipmentWeights = screen.getByRole('button', { name: 'Review shipment weights' });
        await userEvent.click(reviewShipmentWeights);

        expect(screen.getByTestId('shipmentMissingInformation')).toBeInTheDocument();
      });

      it('does not render when the shipment is not missing a reweigh or a prime estimated ewight', async () => {
        useMovePaymentRequestsQueries.mockReturnValue(noAlertsReturnValue);

        renderWithProviders(<ReviewBillableWeight />);

        const reviewShipmentWeights = screen.getByRole('button', { name: 'Review shipment weights' });
        await userEvent.click(reviewShipmentWeights);

        expect(screen.queryByTestId('shipmentMissingInformation')).not.toBeInTheDocument();
      });
    });

    describe('shipment exceeds 110% of estimated weight', () => {
      it('renders the alert when the shipment is overweight - the billable weight is greater than the estimated weight * 110%', async () => {
        useMovePaymentRequestsQueries.mockReturnValue(useMovePaymentRequestsReturnValue);

        renderWithProviders(<ReviewBillableWeight />);

        const reviewShipmentWeights = screen.getByRole('button', { name: 'Review shipment weights' });
        await userEvent.click(reviewShipmentWeights);

        expect(screen.getByTestId('shipmentBillableWeightExceeds110OfEstimated')).toBeInTheDocument();
      });

      it('does not render the alert when the shipment is not overweight - the billable weight is less than the estimated weight * 110%', async () => {
        useMovePaymentRequestsQueries.mockReturnValue(noAlertsReturnValue);

        renderWithProviders(<ReviewBillableWeight />);

        const reviewShipmentWeights = screen.getByRole('button', { name: 'Review shipment weights' });
        await userEvent.click(reviewShipmentWeights);

        expect(screen.queryByTestId('shipmentBillableWeightExceeds110OfEstimated')).not.toBeInTheDocument();
      });

      it('does not render the alert when the shipment an NTS-release', async () => {
        useMovePaymentRequestsQueries.mockReturnValue(useMovePaymentRequestsNTSReleaseReturnValue);

        renderWithProviders(<ReviewBillableWeight />);

        const reviewShipmentWeights = screen.getByRole('button', { name: 'Review shipment weights' });
        await userEvent.click(reviewShipmentWeights);

        expect(screen.queryByTestId('shipmentBillableWeightExceeds110OfEstimated')).not.toBeInTheDocument();
      });
    });
  });

  describe('handles diverted shipments', () => {
    it('displays diversion tags where appropriate', async () => {
      useMovePaymentRequestsQueries.mockReturnValue(useDivertedMovePaymentRequestsReturnValue);

      renderWithProviders(<ReviewBillableWeight />);

      expect(screen.getByText('Review weights')).toBeInTheDocument();
      expect(screen.queryByTestId('tag', { name: 'DIVERSION' })).toBeInTheDocument();
      const sidebarTitle = screen.queryAllByText('Review weights');
      expect(sidebarTitle[0].lastChild.textContent).toEqual('DIVERSION');

      const reviewShipmentWeights = screen.getByRole('button', { name: 'Review shipment weights' });
      await userEvent.click(reviewShipmentWeights);

      expect(screen.getByText('Shipment 1 of 3')).toBeInTheDocument();
      let sidebarSubTitle = screen.queryAllByText('Shipment weights');
      expect(sidebarSubTitle[0].lastChild.textContent).toEqual('DIVERSION');

      const nextShipment = screen.getByRole('button', { name: 'Next Shipment' });
      await userEvent.click(nextShipment);

      expect(screen.getByText('Shipment 2 of 3')).toBeInTheDocument();
      sidebarSubTitle = screen.queryAllByText('Shipment weights');
      expect(sidebarSubTitle[0].lastChild.textContent).toEqual('DIVERSION');

      await userEvent.click(nextShipment);

      expect(screen.getByText('Shipment 3 of 3')).toBeInTheDocument();
      expect(screen.queryByTestId('tag', { name: 'DIVERSION' })).not.toBeInTheDocument();
    });
  });

  describe('handles PPM shipments', () => {
    it('displays PPM information', async () => {
      useMovePaymentRequestsQueries.mockReturnValue(useMovePaymentRequestQueriesReturnValueAllDocs);

      renderWithProviders(<ReviewBillableWeight />);

      expect(screen.getByText('Review weights')).toBeInTheDocument();
      expect(screen.getByText('Max billable weight')).toBeInTheDocument();
      expect(screen.getByText('Actual weight')).toBeInTheDocument();
      expect(screen.getByTestId('maxBillableWeight').textContent).toBe(
        formatWeight(useMovePaymentRequestQueriesReturnValueAllDocs.order.entitlement.authorizedWeight),
      );
      expect(screen.getByTestId('weightAllowance').textContent).toBe(
        formatWeight(useMovePaymentRequestQueriesReturnValueAllDocs.order.entitlement.totalWeight),
      );

      expect(screen.getByText('Weight allowance')).toBeInTheDocument();
      expect(screen.getByText('Actual billable weight')).toBeInTheDocument();
    });
  });
});
