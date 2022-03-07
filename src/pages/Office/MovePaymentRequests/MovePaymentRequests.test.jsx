/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { render, waitFor, screen, within } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import MovePaymentRequests from './MovePaymentRequests';

import MOVE_STATUSES from 'constants/moves';
import { MockProviders } from 'testUtils';
import { useMovePaymentRequestsQueries } from 'hooks/queries';
import { shipmentStatuses } from 'constants/shipments';
import SERVICE_ITEM_STATUSES from 'constants/serviceItems';

jest.mock('hooks/queries', () => ({
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

const testProps = {
  setUnapprovedShipmentCount: jest.fn(),
  setUnapprovedServiceItemCount: jest.fn(),
  setPendingPaymentRequestCount: jest.fn(),
};

const move = {
  id: '1',
  contractor: {
    contractNumber: 'HTC-123-3456',
  },
  orders: {
    sac: '1234456',
    tac: '1213',
  },
  billableWeightsReviewedAt: '2021-06-01',
};

const order = {
  orders_type: 'PERMANENT_CHANGE_OF_STATION',
  has_dependents: false,
  issue_date: '2020-08-11',
  grade: 'RANK',
  moves: ['123'],
  origin_duty_location: {
    name: 'Test Duty Station',
    address: {
      postalCode: '123456',
    },
  },
  new_duty_location: {
    name: 'New Test Duty Station',
    address: {
      postalCode: '123456',
    },
  },
  report_by_date: '2020-08-31',
  service_member_id: '666',
  spouse_has_pro_gear: false,
  status: MOVE_STATUSES.SUBMITTED,
  uploaded_orders: {
    uploads: [],
  },
  entitlement: {
    authorizedWeight: 8000,
    dependentsAuthorized: true,
    eTag: 'MjAyMS0wOC0yNFQxODoyNDo0MC45NzIzMTha',
    id: '188842d1-cf88-49ec-bd2f-dfa98da44bb2',
    nonTemporaryStorage: true,
    organizationalClothingAndIndividualEquipment: true,
    privatelyOwnedVehicle: true,
    proGearWeight: 2000,
    proGearWeightSpouse: 500,
    requiredMedicalEquipmentWeight: 1000,
    storageInTransit: 2,
    totalDependents: 1,
    totalWeight: 8000,
  },
};

const multiplePaymentRequests = {
  paymentRequests: [
    {
      id: '09474c6a-69b6-4501-8e08-670a12512e5f',
      createdAt: '2020-12-01T00:00:00.000Z',
      paymentRequestNumber: '1843-9061-1',
      status: 'REVIEWED',
      moveTaskOrderID: '1',
      moveTaskOrder: move,
      serviceItems: [
        {
          id: '09474c6a-69b6-4501-8e08-670a12512a5f',
          createdAt: '2020-12-01T00:00:00.000Z',
          mtoServiceItemID: 'f8c2f97f-99e7-4fb1-9cc4-473debd24dbc',
          priceCents: 2000001,
          status: 'APPROVED',
        },
        {
          id: '39474c6a-69b6-4501-8e08-670a12512a5f',
          createdAt: '2020-12-01T00:00:00.000Z',
          mtoServiceItemID: 'a8c2f97f-99e7-4fb1-9cc4-473debd24dbc',
          priceCents: 4000001,
          status: 'DENIED',
          rejectionReason: 'Requested amount exceeds guideline',
        },
      ],
      reviewedAt: '2020-12-01T00:00:00.000Z',
    },
    {
      id: '29474c6a-69b6-4501-8e08-670a12512e5f',
      createdAt: '2020-12-01T00:00:00.000Z',
      paymentRequestNumber: '1843-9061-2',
      status: 'PENDING',
      moveTaskOrderID: '1',
      moveTaskOrder: move,
      serviceItems: [
        {
          id: '09474c6a-69b6-4501-8e08-670a12512a5f',
          createdAt: '2020-12-01T00:00:00.000Z',
          mtoServiceItemID: 'f8c2f97f-99e7-4fb1-9cc4-473debd24dbc',
          priceCents: 2000001,
          status: 'REQUESTED',
        },
        {
          id: '39474c6a-69b6-4501-8e08-670a12512a5f',
          createdAt: '2020-12-01T00:00:00.000Z',
          mtoServiceItemID: 'a8c2f97f-99e7-4fb1-9cc4-473debd24dbc',
          priceCents: 4000001,
          status: 'REQUESTED',
        },
      ],
    },
  ],
  mtoShipments: [
    {
      shipmentType: 'HHG',
      id: '2',
      moveTaskOrderID: '1',
      status: shipmentStatuses.APPROVED,
      scheduledPickupDate: '2020-01-09T00:00:00.000Z',
      destinationAddress: { city: 'Princeton', state: 'NJ', postalCode: '08540' },
      pickupAddress: { city: 'Boston', state: 'MA', postalCode: '02101' },
      calculatedBillableWeight: 9000,
      primeActualWeight: 500,
      reweigh: {
        id: 'reweighID1',
        weight: 100,
      },
      mtoServiceItems: [
        {
          id: '5',
          mtoShipmentID: '2',
          status: SERVICE_ITEM_STATUSES.APPROVED,
        },
        {
          id: '6',
          status: SERVICE_ITEM_STATUSES.REJECTED,
          mtoShipmentID: '2',
        },
        {
          id: '7',
          status: SERVICE_ITEM_STATUSES.SUBMITTED,
          mtoShipmentID: '2',
        },
      ],
    },
    {
      shipmentType: 'HHG',
      id: '3',
      moveTaskOrderID: '1',
      status: shipmentStatuses.APPROVED,
      scheduledPickupDate: '2020-01-10T00:00:00.000Z',
      destinationAddress: { city: 'Princeton', state: 'NJ', postalCode: '08540' },
      pickupAddress: { city: 'Boston', state: 'MA', postalCode: '02101' },
      calculatedBillableWeight: 1000,
      primeActualWeight: 5000,
      reweigh: {
        id: 'reweighID2',
        weight: 600,
      },
      mtoServiceItems: [
        {
          id: '9',
          mtoShipmentID: '3',
          status: SERVICE_ITEM_STATUSES.APPROVED,
        },
        {
          id: '10',
          status: SERVICE_ITEM_STATUSES.REJECTED,
          mtoShipmentID: '3',
        },
        {
          id: '11',
          status: SERVICE_ITEM_STATUSES.SUBMITTED,
          mtoShipmentID: '3',
        },
      ],
    },
    {
      shipmentType: 'HHG',
      id: '4',
      moveTaskOrderID: '1',
      status: shipmentStatuses.SUBMITTED,
      scheduledPickupDate: '2020-01-11T00:00:00.000Z',
      destinationAddress: { city: 'Princeton', state: 'NJ', postalCode: '08540' },
      pickupAddress: { city: 'Boston', state: 'MA', postalCode: '02101' },
      calculatedBillableWeight: 2000,
      primeActualWeight: 300,
      reweigh: {
        id: 'reweighID3',
        weight: 900,
      },
      mtoServiceItems: [
        {
          id: '12',
          mtoShipmentID: '4',
          status: SERVICE_ITEM_STATUSES.APPROVED,
        },
        {
          id: '13',
          status: SERVICE_ITEM_STATUSES.REJECTED,
          mtoShipmentID: '4',
        },
        {
          id: '14',
          status: SERVICE_ITEM_STATUSES.SUBMITTED,
          mtoShipmentID: '4',
        },
      ],
    },
  ],
  order,
};

const singleReviewedPaymentRequest = {
  paymentRequests: [
    {
      id: '09474c6a-69b6-4501-8e08-670a12512e5f',
      createdAt: '2020-12-01T00:00:00.000Z',
      paymentRequestNumber: '1843-9061-1',
      status: 'REVIEWED',
      moveTaskOrderID: '1',
      moveTaskOrder: move,
      serviceItems: [
        {
          id: '09474c6a-69b6-4501-8e08-670a12512a5f',
          createdAt: '2020-12-01T00:00:00.000Z',
          mtoServiceItemID: 'f8c2f97f-99e7-4fb1-9cc4-473debd24dbc',
          priceCents: 2000001,
          status: 'APPROVED',
        },
        {
          id: '39474c6a-69b6-4501-8e08-670a12512a5f',
          createdAt: '2020-12-01T00:00:00.000Z',
          mtoServiceItemID: 'a8c2f97f-99e7-4fb1-9cc4-473debd24dbc',
          priceCents: 4000001,
          status: 'DENIED',
          rejectionReason: 'Requested amount exceeds guideline',
        },
      ],
      reviewedAt: '2020-12-01T00:00:00.000Z',
    },
  ],
  mtoShipments: [
    {
      shipmentType: 'HHG',
      id: '2',
      moveTaskOrderID: '1',
      status: shipmentStatuses.APPROVED,
      scheduledPickupDate: '2020-01-11T00:00:00.000Z',
      destinationAddress: { city: 'Princeton', state: 'NJ', postalCode: '08540' },
      pickupAddress: { city: 'Boston', state: 'MA', postalCode: '02101' },
      calculatedBillableWeight: 2000,
      primeActualWeight: 300,
      reweigh: {
        id: 'reweighID',
        weight: 900,
      },
      mtoServiceItems: [
        {
          id: '3',
          mtoShipmentID: '2',
          status: SERVICE_ITEM_STATUSES.APPROVED,
        },
      ],
    },
  ],
  order,
};

const emptyPaymentRequests = {
  paymentRequests: [],
  mtoShipments: [
    {
      shipmentType: 'HHG',
      id: '2',
      moveTaskOrderID: '1',
      status: shipmentStatuses.APPROVED,
      scheduledPickupDate: '2020-01-11T00:00:00.000Z',
      destinationAddress: { city: 'Princeton', state: 'NJ', postalCode: '08540' },
      pickupAddress: { city: 'Boston', state: 'MA', postalCode: '02101' },
      calculatedBillableWeight: 2000,
      primeActualWeight: 300,
      reweigh: {
        id: 'reweighID',
        weight: 900,
      },
      mtoServiceItems: [
        {
          id: '3',
          mtoShipmentID: '2',
          status: SERVICE_ITEM_STATUSES.APPROVED,
        },
      ],
    },
  ],
  order,
};

const moveShipmentOverweight = {
  paymentRequests: [],
  mtoShipments: [
    {
      shipmentType: 'HHG',
      id: '2',
      moveTaskOrderID: '1',
      status: shipmentStatuses.APPROVED,
      scheduledPickupDate: '2020-01-11T00:00:00.000Z',
      destinationAddress: { city: 'Princeton', state: 'NJ', postalCode: '08540' },
      pickupAddress: { city: 'Boston', state: 'MA', postalCode: '02101' },
      calculatedBillableWeight: 5000,
      primeActualWeight: 7000,
      primeEstimatedWeight: 3000,
      mtoServiceItems: [
        {
          id: '3',
          mtoShipmentID: '2',
          status: SERVICE_ITEM_STATUSES.APPROVED,
        },
      ],
    },
  ],
  order,
};

const moveShipmentMissingReweighWeight = {
  paymentRequests: [],
  mtoShipments: [
    {
      shipmentType: 'HHG',
      id: '2',
      moveTaskOrderID: '1',
      status: shipmentStatuses.APPROVED,
      scheduledPickupDate: '2020-01-11T00:00:00.000Z',
      destinationAddress: { city: 'Princeton', state: 'NJ', postalCode: '08540' },
      pickupAddress: { city: 'Boston', state: 'MA', postalCode: '02101' },
      calculatedBillableWeight: 2000,
      primeActualWeight: 8000,
      primeEstimatedWeight: 3000,
      reweigh: {
        id: '123',
      },
      mtoServiceItems: [
        {
          id: '3',
          mtoShipmentID: '2',
          status: SERVICE_ITEM_STATUSES.APPROVED,
        },
      ],
    },
  ],
  order,
};

const returnWithBillableWeightsReviewed = {
  paymentRequests: [],
  mtoShipments: [
    {
      shipmentType: 'HHG',
      id: '2',
      moveTaskOrderID: '1',
      status: shipmentStatuses.APPROVED,
      scheduledPickupDate: '2020-01-11T00:00:00.000Z',
      destinationAddress: { city: 'Princeton', state: 'NJ', postalCode: '08540' },
      pickupAddress: { city: 'Boston', state: 'MA', postalCode: '02101' },
      calculatedBillableWeight: 2000,
      primeActualWeight: 8000,
      primeEstimatedWeight: 3000,
      reweigh: {
        id: '123',
      },
      mtoServiceItems: [
        {
          id: '3',
          mtoShipmentID: '2',
          status: SERVICE_ITEM_STATUSES.APPROVED,
        },
      ],
    },
  ],
  order,
  move,
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

function renderMovePaymentRequests(props) {
  return render(
    <MockProviders initialEntries={[`/moves/L2BKD6/payment-requests`]}>
      <MovePaymentRequests {...props} />
    </MockProviders>,
  );
}

describe('MovePaymentRequests', () => {
  describe('check loading and error component states', () => {
    it('renders the Loading Placeholder when the query is still loading', async () => {
      useMovePaymentRequestsQueries.mockReturnValue(loadingReturnValue);

      renderMovePaymentRequests(testProps);

      const h2 = await screen.getByRole('heading', { name: 'Loading, please wait...', level: 2 });
      expect(h2).toBeInTheDocument();
    });

    it('renders the Something Went Wrong component when the query errors', async () => {
      useMovePaymentRequestsQueries.mockReturnValue(errorReturnValue);

      renderMovePaymentRequests(testProps);

      const errorMessage = await screen.getByText(/Something went wrong./);
      expect(errorMessage).toBeInTheDocument();
    });
  });

  describe('with multiple payment requests', () => {
    beforeEach(() => {
      useMovePaymentRequestsQueries.mockReturnValue(multiplePaymentRequests);
    });

    it('renders without errors', () => {
      renderMovePaymentRequests(testProps);
      expect(screen.getByTestId('MovePaymentRequests')).toBeInTheDocument();
    });

    it('renders multiple payment requests', async () => {
      renderMovePaymentRequests(testProps);
      await waitFor(() => {
        multiplePaymentRequests.paymentRequests.forEach((pr) => {
          expect(screen.getByText(`Payment Request ${pr.paymentRequestNumber}`)).toBeInTheDocument();
        });
      });
    });

    it('updates the pending payment request count callback', async () => {
      renderMovePaymentRequests(testProps);
      await waitFor(() => {
        expect(testProps.setPendingPaymentRequestCount).toHaveBeenCalledWith(1);
      });
    });

    it('updates the unapproved shipments tag callback', async () => {
      renderMovePaymentRequests(testProps);
      await waitFor(() => {
        expect(testProps.setUnapprovedShipmentCount).toHaveBeenCalledWith(1);
      });
    });

    it('updates the unapproved service items tag callback', async () => {
      renderMovePaymentRequests(testProps);
      await waitFor(() => {
        expect(testProps.setUnapprovedServiceItemCount).toHaveBeenCalledWith(2);
      });
    });

    it('displays the number of payment request on leftnav sidebar', async () => {
      renderMovePaymentRequests(testProps);
      await waitFor(() => {
        expect(screen.getByTestId('numOfPaymentRequestsTag').textContent).toEqual('2');
      });
    });
  });

  describe('renders side navigation for each section', () => {
    beforeEach(() => {
      useMovePaymentRequestsQueries.mockReturnValue(singleReviewedPaymentRequest);
    });

    it.each([
      ['Payment requests', '#payment-requests'],
      ['Billable weights', '#billable-weights'],
    ])('renders the %s side navigation', (name, tag) => {
      renderMovePaymentRequests(testProps);
      const leftNav = screen.getByRole('navigation');
      expect(leftNav).toBeInTheDocument();

      const paymentRequstNavLink = within(leftNav).getByText(name);

      expect(paymentRequstNavLink.href).toContain(tag);
      expect(paymentRequstNavLink.text).toContain(name);
    });

    it('displays the number of payment request on leftnav sidebar', async () => {
      renderMovePaymentRequests(testProps);
      await waitFor(() => {
        expect(screen.getByTestId('numOfPaymentRequestsTag').textContent).toEqual('1');
      });
    });
  });

  describe('with one reviewed payment request', () => {
    it('updates the pending payment request count callback', async () => {
      renderMovePaymentRequests(testProps);
      await waitFor(() => {
        expect(testProps.setPendingPaymentRequestCount).toHaveBeenCalledWith(0);
      });
    });

    it('updates the unapproved shipment count callback', async () => {
      renderMovePaymentRequests(testProps);
      await waitFor(() => {
        expect(testProps.setUnapprovedShipmentCount).toHaveBeenCalledWith(0);
      });
    });

    it('updates the unapproved service item count callback', async () => {
      renderMovePaymentRequests(testProps);
      await waitFor(() => {
        expect(testProps.setUnapprovedServiceItemCount).toHaveBeenCalledWith(0);
      });
    });
  });

  describe('with no payment requests for move', () => {
    beforeEach(() => {
      useMovePaymentRequestsQueries.mockReturnValue(emptyPaymentRequests);
    });

    it('renders side navigation for payment request section', () => {
      renderMovePaymentRequests(testProps);
      const leftNav = screen.getByRole('navigation');
      expect(leftNav).toBeInTheDocument();

      const paymentRequstNavLink = within(leftNav).queryByText('Payment requests');

      expect(paymentRequstNavLink).toBeInTheDocument();
    });

    it('renders with empty message when no payment requests exist', async () => {
      renderMovePaymentRequests(testProps);
      await waitFor(() => {
        expect(screen.getByText('No payment requests have been submitted for this move yet.')).toBeInTheDocument();
      });
    });

    it('updates the pending payment request count callback', async () => {
      renderMovePaymentRequests(testProps);
      await waitFor(() => {
        expect(testProps.setPendingPaymentRequestCount).toHaveBeenCalledWith(0);
      });
    });

    it('updates the unapproved shipment count callback', async () => {
      renderMovePaymentRequests(testProps);
      await waitFor(() => {
        expect(testProps.setUnapprovedShipmentCount).toHaveBeenCalledWith(0);
      });
    });

    it('updates the unapproved service item count callback', async () => {
      renderMovePaymentRequests(testProps);
      await waitFor(() => {
        expect(testProps.setUnapprovedServiceItemCount).toHaveBeenCalledWith(0);
      });
    });
  });

  describe('a billable weight that does not exceed the max billable weight', () => {
    beforeEach(() => {
      useMovePaymentRequestsQueries.mockReturnValue(emptyPaymentRequests);
    });

    it('does not show the max billable weight tag in sidebar', async () => {
      renderMovePaymentRequests(testProps);
      await waitFor(() => {
        expect(screen.queryByTestId('maxBillableWeightErrorTag')).not.toBeInTheDocument();
      });
    });

    it('does not show the max billable weight error text in the billable weight card', async () => {
      renderMovePaymentRequests(testProps);
      await waitFor(() => {
        expect(screen.queryByText('Move exceeds max billable weight')).not.toBeInTheDocument();
      });
    });

    it('navigates the user to the reivew billable weight page', async () => {
      jest.spyOn(console, 'error').mockImplementation(() => {});
      renderMovePaymentRequests(testProps);

      const reviewWeights = screen.getByRole('button', { name: 'Review weights' });

      userEvent.click(reviewWeights);

      await waitFor(() => {
        expect(mockPush).toHaveBeenCalledWith('/moves/testMoveCode/billable-weight');
      });
    });
  });

  describe('a billable weight that exceeds the max billable weight', () => {
    beforeEach(() => {
      useMovePaymentRequestsQueries.mockReturnValue(multiplePaymentRequests);
    });

    it('shows the max billable weight tag in sidebar', async () => {
      renderMovePaymentRequests(testProps);
      await waitFor(() => {
        expect(screen.getByTestId('maxBillableWeightErrorTag')).toBeInTheDocument();
      });
    });

    it('shows the max billable weight error text in the billable weight card', async () => {
      renderMovePaymentRequests(testProps);
      await waitFor(() => {
        expect(screen.getByText('Move exceeds max billable weight')).toBeInTheDocument();
      });
    });
  });

  describe('a move that has an overweight shipment displays a warning tag', () => {
    beforeEach(() => {
      useMovePaymentRequestsQueries.mockReturnValue(moveShipmentOverweight);
    });

    it('shows the max billable weight warning tag in sidebar', async () => {
      renderMovePaymentRequests(testProps);
      await waitFor(() => {
        expect(screen.getByTestId('maxBillableWeightWarningTag')).toBeInTheDocument();
      });
    });
  });

  describe('a move that has a missing shipment reweigh weight displays a warning tag', () => {
    beforeEach(() => {
      useMovePaymentRequestsQueries.mockReturnValue(moveShipmentMissingReweighWeight);
    });

    it('shows the max billable weight warning tag in sidebar', async () => {
      renderMovePaymentRequests(testProps);
      await waitFor(() => {
        expect(screen.getByTestId('maxBillableWeightWarningTag')).toBeInTheDocument();
      });
    });
  });

  describe('a move that does not have a billableWeightsReviewedAt timestamp displays a primary styled Review Weights btn', () => {
    beforeEach(() => {
      useMovePaymentRequestsQueries.mockReturnValue(moveShipmentMissingReweighWeight);
    });

    it('shows the max billable weight warning tag in sidebar', async () => {
      renderMovePaymentRequests(testProps);

      const reviewWeights = screen.getByRole('button', { name: 'Review weights' });
      expect(reviewWeights).not.toHaveClass('usa-button--secondary');
    });
  });

  describe('a move that has a billableWeightsReviewedAt timestamp displays a secondary styled Review Weights btn', () => {
    beforeEach(() => {
      useMovePaymentRequestsQueries.mockReturnValue(returnWithBillableWeightsReviewed);
    });

    it('shows the max billable weight warning tag in sidebar', async () => {
      renderMovePaymentRequests(testProps);

      const reviewWeights = screen.getByRole('button', { name: 'Review weights' });
      expect(reviewWeights).toHaveClass('usa-button--secondary');
    });
  });
});
