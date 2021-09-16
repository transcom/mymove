import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import ReviewBillableWeight from './ReviewBillableWeight';

import { formatWeight, formatDateFromIso } from 'shared/formatters';
import { useOrdersDocumentQueries, useMovePaymentRequestsQueries } from 'hooks/queries';
import { shipmentStatuses } from 'constants/shipments';

// Mock the document viewer since we're not really testing that aspect here.
// Document Viewer tests should be covered in the component itself.
jest.mock('components/DocumentViewer/DocumentViewer', () => {
  const MockDocumentViewer = () => <div>Document viewer text</div>;
  return MockDocumentViewer;
});

jest.mock('hooks/queries', () => ({
  useOrdersDocumentQueries: jest.fn(),
  useMovePaymentRequestsQueries: jest.fn(),
}));

const mockPush = jest.fn();

jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useParams: jest.fn().mockReturnValue({ moveCode: 'testMoveCode' }),
  useHistory: () => ({
    push: mockPush,
  }),
}));

const mockOriginDutyStation = {
  address: {
    city: 'Des Moines',
    country: 'US',
    eTag: 'MjAyMC0wOS0xNFQxNzo0MTozOC42OTg1OTha',
    id: '2e26b066-aaca-4563-b284-d7f3f978fb3c',
    postal_code: '50309',
    state: 'IA',
    street_address_1: '987 Other Avenue',
    street_address_2: 'P.O. Box 1234',
    street_address_3: 'c/o Another Person',
  },
  address_id: '2e26b066-aaca-4563-b284-d7f3f978fb3c',
  eTag: 'MjAyMC0wOS0xNFQxNzo0MTozOC43MDcxOTVa',
  id: 'a3ec2bdd-aa0a-434a-ba58-34c85f047704',
  name: 'XBc1KNi3pA',
};

const mockDestinationDutyStation = {
  address: {
    city: 'Augusta',
    country: 'United States',
    eTag: 'MjAyMC0wOS0xNFQxNzo0MDo0OC44OTM3MDVa',
    id: '5ac95be8-0230-47ea-90b4-b0f6f60de364',
    postal_code: '30813',
    state: 'GA',
    street_address_1: 'Fort Gordon',
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
    department_indicator: 'AIR_FORCE',
    destinationDutyStation: mockDestinationDutyStation,
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
    order_number: 'ORDER3',
    order_type: 'PERMANENT_CHANGE_OF_STATION',
    order_type_detail: 'HHG_PERMITTED',
    originDutyStation: mockOriginDutyStation,
    report_by_date: '2018-08-01',
    tac: 'F8E1',
    sac: 'E2P3',
  },
};

const mockMtoShipments = [
  {
    id: 1,
    status: shipmentStatuses.APPROVED,
    calculatedBillableWeight: 3000,
    billableWeightCap: 1000,
    primeEstimatedWeight: 1000,
    primeActualWeight: 300,
    reweigh: { verificationReason: 'reweigh required', requestedAt: '2021-09-01' },
    pickupAddress: { city: 'Las Vegas', state: 'NV', postal_code: '90210' },
    destinationAddress: { city: 'Miami', state: 'FL', postal_code: '33607' },
    actualPickupDate: '2021-08-31',
  },
  {
    id: 2,
    status: shipmentStatuses.APPROVED,
    calculatedBillableWeightCap: 2000,
    billableWeightCap: 2000,
    primeEstimatedWeight: 2000,
    primeActualWeight: 400,
    reweigh: { weight: 1000, verificationReason: 'reweigh required', requestedAt: '2021-09-01' },
    pickupAddress: { city: 'Las Vegas', state: 'NV', postal_code: '90210' },
    destinationAddress: { city: 'Miami', state: 'FL', postal_code: '33607' },
    actualPickupDate: '2021-08-31',
  },
  {
    id: 3,
    status: shipmentStatuses.DIVERSION_REQUESTED,
    calculatedBillableWeight: 3000,
    billableWeightCap: 3000,
    primeEstimatedWeight: 7000,
    primeActualWeight: 300,
    reweigh: { weight: 200, verificationReason: 'reweigh required', requestedAt: '2021-09-01' },
    pickupAddress: { city: 'Las Vegas', state: 'NV', postal_code: '90210' },
    destinationAddress: { city: 'Miami', state: 'FL', postal_code: '33607' },
    actualPickupDate: '2021-08-31',
  },
];

const useOrdersDocumentQueriesReturnValue = {
  orders: mockOrders,
  upload: {
    z: {
      id: 'z',
      filename: 'test.pdf',
      contentType: 'application/pdf',
      url: '/storage/user/1/uploads/2?contentType=application%2Fpdf',
    },
  },
};

const useMovePaymentRequestsReturnValue = {
  order: mockOrders['1'],
  mtoShipments: mockMtoShipments,
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

describe('ReviewBillableWeight', () => {
  describe('check loading and error component states', () => {
    it('renders the loading placeholder when the query is still loading', async () => {
      useOrdersDocumentQueries.mockReturnValue(loadingReturnValue);
      useMovePaymentRequestsQueries.mockReturnValue(useMovePaymentRequestsReturnValue);

      render(<ReviewBillableWeight />);

      const h2 = await screen.getByRole('heading', { name: 'Loading, please wait...', level: 2 });
      expect(h2).toBeInTheDocument();
    });

    it('renders the Something Went Wrong component when the query errors', async () => {
      useOrdersDocumentQueries.mockReturnValue(errorReturnValue);
      useMovePaymentRequestsQueries.mockReturnValue(useMovePaymentRequestsReturnValue);

      render(<ReviewBillableWeight />);

      const errorMessage = await screen.getByText(/Something went wrong./);
      expect(errorMessage).toBeInTheDocument();
    });
  });

  it('renders the component', () => {
    useOrdersDocumentQueries.mockReturnValue(useOrdersDocumentQueriesReturnValue);
    useMovePaymentRequestsQueries.mockReturnValue(useMovePaymentRequestsReturnValue);

    render(<ReviewBillableWeight />);
    expect(screen.getByText('Review weights')).toBeInTheDocument();
    expect(screen.getByText('Document viewer text')).toBeInTheDocument();
  });

  it('takes the user back to the payment requests page when x is clicked', async () => {
    useOrdersDocumentQueries.mockReturnValue(useOrdersDocumentQueriesReturnValue);

    render(<ReviewBillableWeight />);

    const xButton = screen.getByTestId('closeSidebar');

    userEvent.click(xButton);

    await waitFor(() => {
      expect(mockPush).toHaveBeenCalledWith('/moves/testMoveCode/payment-requests');
    });
  });

  it('takes the user to review the shipment weights when the review weights button is clicked', async () => {
    useOrdersDocumentQueries.mockReturnValue(useOrdersDocumentQueriesReturnValue);
    useMovePaymentRequestsQueries.mockReturnValue(useMovePaymentRequestsReturnValue);

    render(<ReviewBillableWeight />);

    const reviewShipmentWeights = screen.getByRole('button', { name: 'Review shipment weights' });

    userEvent.click(reviewShipmentWeights);

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
    useOrdersDocumentQueries.mockReturnValue(useOrdersDocumentQueriesReturnValue);
    useMovePaymentRequestsQueries.mockReturnValue(useMovePaymentRequestsReturnValue);

    render(<ReviewBillableWeight />);

    const reviewShipmentWeights = screen.getByRole('button', { name: 'Review shipment weights' });

    userEvent.click(reviewShipmentWeights);

    expect(screen.getByText('Shipment 1 of 3')).toBeInTheDocument();
    expect(screen.getByTestId('estimatedWeight').textContent).toBe(
      formatWeight(mockMtoShipments[0].primeEstimatedWeight),
    );
    expect(screen.getByTestId('originalWeight').textContent).toBe(formatWeight(mockMtoShipments[0].primeActualWeight));
    expect(screen.getByTestId('reweighWeight').textContent).toBe('Missing');
    expect(screen.getByTestId('dateReweighRequested').textContent).toBe(
      formatDateFromIso(mockMtoShipments[0].reweigh.requestedAt, 'DD MMM YYYY'),
    );
    expect(screen.getByTestId('reweighRemarks').textContent).toBe(mockMtoShipments[0].reweigh.verificationReason);

    const nextShipment = screen.getByRole('button', { name: 'Next Shipment' });
    userEvent.click(nextShipment);
    await waitFor(() => {
      expect(screen.getByText('Shipment 2 of 3')).toBeInTheDocument();
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
    });
    userEvent.click(nextShipment);
    await waitFor(() => {
      expect(screen.getByText('Shipment 3 of 3')).toBeInTheDocument();
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
  });

  it('takes the user to the previous shipment when the Back button is clicked', async () => {
    useOrdersDocumentQueries.mockReturnValue(useOrdersDocumentQueriesReturnValue);
    useMovePaymentRequestsQueries.mockReturnValue(useMovePaymentRequestsReturnValue);

    render(<ReviewBillableWeight />);

    const reviewShipmentWeights = screen.getByRole('button', { name: 'Review shipment weights' });

    userEvent.click(reviewShipmentWeights);

    const nextShipment = screen.getByRole('button', { name: 'Next Shipment' });
    userEvent.click(nextShipment);
    userEvent.click(nextShipment);
    await waitFor(() => {
      expect(screen.getByText('Shipment 3 of 3')).toBeInTheDocument();
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
    });

    const back = screen.getByRole('button', { name: 'Back' });
    userEvent.click(back);
    await waitFor(() => {
      expect(screen.getByText('Shipment 2 of 3')).toBeInTheDocument();
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
    });

    userEvent.click(back);
    await waitFor(() => {
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
    });

    userEvent.click(back);
    await waitFor(() => {
      expect(screen.getByText('Edit max billable weight')).toBeInTheDocument();
      expect(screen.queryByRole('button', { name: 'Next Shipment' })).not.toBeInTheDocument();
      expect(screen.queryByRole('button', { name: 'Back' })).not.toBeInTheDocument();
      expect(screen.queryByRole('button', { name: 'Review shipment weights' })).toBeInTheDocument();
    });
  });

  it('renders weight summary', () => {
    useOrdersDocumentQueries.mockReturnValue(useOrdersDocumentQueriesReturnValue);
    useMovePaymentRequestsQueries.mockReturnValue(useMovePaymentRequestsReturnValue);
    render(<ReviewBillableWeight />);
    expect(screen.getByTestId('maxBillableWeight').textContent).toBe(
      formatWeight(useMovePaymentRequestsReturnValue.order.entitlement.authorizedWeight),
    );
    expect(screen.getByTestId('weightAllowance').textContent).toBe(
      formatWeight(useMovePaymentRequestsReturnValue.order.entitlement.totalWeight),
    );
    expect(screen.getByTestId('weightRequested').textContent).toBe('900 lbs');
    expect(screen.getByTestId('totalBillableWeight').textContent).toBe('6,000 lbs');
  });

  it('renders max billable weight and edit view', () => {
    useOrdersDocumentQueries.mockReturnValue(useOrdersDocumentQueriesReturnValue);
    useMovePaymentRequestsQueries.mockReturnValue(useMovePaymentRequestsReturnValue);
    const weightAllowance = formatWeight(useMovePaymentRequestsReturnValue.order.entitlement.totalWeight);

    render(<ReviewBillableWeight />);

    userEvent.click(screen.getByText('Edit'));
    expect(screen.getByTestId('maxWeight-weightAllowance').textContent).toBe(weightAllowance);
    expect(screen.getByTestId('maxWeight-estimatedWeight').textContent).toBe('11,000 lbs');
  });
});
