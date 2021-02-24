import React from 'react';
import { mount } from 'enzyme';

import MovePaymentRequests from './MovePaymentRequests';

import { MockProviders } from 'testUtils';
import { useMovePaymentRequestsQueries } from 'hooks/queries';

jest.mock('hooks/queries', () => ({
  useMovePaymentRequestsQueries: jest.fn(),
}));

const move = {
  contractor: {
    contractNumber: 'HTC-123-3456',
  },
  orders: {
    sac: '1234456',
    tac: '1213',
  },
};

const multiplePaymentRequests = {
  paymentRequests: [
    {
      id: '09474c6a-69b6-4501-8e08-670a12512e5f',
      createdAt: '2020-12-01T00:00:00.000Z',
      moveTaskOrderID: 'f8c2f97f-99e7-4fb1-9cc4-473debd04dbc',
      paymentRequestNumber: '1843-9061-1',
      status: 'REVIEWED',
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
      moveTaskOrderID: 'f8c2f97f-99e7-4fb1-9cc4-473debd04dbc',
      paymentRequestNumber: '1843-9061-2',
      status: 'PENDING',
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
      destinationAddress: { city: 'Princeton', state: 'NJ', postal_code: '08540' },
      pickupAddress: { city: 'Boston', state: 'MA', postal_code: '02101' },
    },
  ],
};

const emptyPaymentRequests = {
  paymentRequests: [],
  mtoShipments: [],
};

describe('MovePaymentRequests', () => {
  describe('multiple payment requests for move', () => {
    useMovePaymentRequestsQueries.mockImplementation(() => multiplePaymentRequests);

    const testMoveCode = 'L2BKD6';
    const component = mount(
      <MockProviders initialEntries={[`/moves/${testMoveCode}/payment-requests`]}>
        <MovePaymentRequests />
      </MockProviders>,
    );

    it('renders without errors', () => {
      expect(component.find('h1').contains('Payment Requests')).toBe(true);
    });

    it('renders mulitple payment requests', () => {
      expect(component.find('PaymentRequestCard').length).toBe(2);
    });
  });

  describe('no payment requests for move', () => {
    useMovePaymentRequestsQueries.mockImplementation(() => emptyPaymentRequests);

    const component = mount(
      <MockProviders initialEntries={[`/moves/FG7W32/payment-requests`]}>
        <MovePaymentRequests />
      </MockProviders>,
    );

    it('renders with empty message when no payment requests exist', () => {
      expect(component.find('p').contains('No payment requests have been submitted for this move yet.')).toBe(true);
    });
  });
});
