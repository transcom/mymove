import React from 'react';
import { mount } from 'enzyme';

import PaymentRequestCard from './PaymentRequestCard';

import { MockProviders } from 'testUtils';

jest.mock('hooks/queries', () => ({
  useMovePaymentRequestsQueries: () => {
    return {
      paymentRequests: [
        {
          id: '09474c6a-69b6-4501-8e08-670a12512e5f',
          createdAt: '2020-12-01T00:00:00.000Z',
          moveTaskOrderID: 'f8c2f97f-99e7-4fb1-9cc4-473debd04dbc',
          paymentRequestNumber: '1843-9061-1',
          status: 'REVIEWED',
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
    };
  },
}));

const testMoveLocator = 'AF7K1P';

describe('PaymentRequestCard', () => {
  describe('pending payment request', () => {
    const pendingPaymentRequest = {
      id: '29474c6a-69b6-4501-8e08-670a12512e5f',
      createdAt: '2020-12-01T00:00:00.000Z',
      moveTaskOrderID: 'f8c2f97f-99e7-4fb1-9cc4-473debd04dbc',
      paymentRequestNumber: '1843-9061-2',
      status: 'PENDING',
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
    };
    const wrapper = mount(
      <MockProviders initialEntries={[`/moves/${testMoveLocator}/payment-requests`]}>
        <PaymentRequestCard paymentRequest={pendingPaymentRequest} />
      </MockProviders>,
    );

    it('renders the needs review status tag', () => {
      expect(wrapper.find({ 'data-testid': 'tag' }).contains('Needs Review')).toBe(true);
    });

    it('sums the service items total', () => {
      expect(wrapper.find('.amountRequested').contains('$60,000.02')).toBe(true);
    });

    it('renders the view orders link', () => {
      const viewLink = wrapper.find('.footer a');

      expect(viewLink.contains('View orders')).toBe(true);
      expect(viewLink.prop('href')).toBe('orders');
    });
  });

  describe('reviewed payment request', () => {
    const reviewedPaymentRequest = {
      id: '29474c6a-69b6-4501-8e08-670a12512e5f',
      createdAt: '2020-12-01T00:00:00.000Z',
      moveTaskOrderID: 'f8c2f97f-99e7-4fb1-9cc4-473debd04dbc',
      paymentRequestNumber: '1843-9061-2',
      status: 'REVIEWED',
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
          rejectionReason: 'duplicate charge',
        },
      ],
    };

    const rejectedPaymentRequest = {
      id: '29474c6a-69b6-4501-8e08-670a12512e5f',
      createdAt: '2020-12-01T00:00:00.000Z',
      moveTaskOrderID: 'f8c2f97f-99e7-4fb1-9cc4-473debd04dbc',
      paymentRequestNumber: '1843-9061-2',
      status: 'REVIEWED',
      serviceItems: [
        {
          id: '09474c6a-69b6-4501-8e08-670a12512a5f',
          createdAt: '2020-12-01T00:00:00.000Z',
          mtoServiceItemID: 'f8c2f97f-99e7-4fb1-9cc4-473debd24dbc',
          priceCents: 2000001,
          status: 'DENIED',
        },
        {
          id: '39474c6a-69b6-4501-8e08-670a12512a5f',
          createdAt: '2020-12-01T00:00:00.000Z',
          mtoServiceItemID: 'a8c2f97f-99e7-4fb1-9cc4-473debd24dbc',
          priceCents: 4000001,
          status: 'DENIED',
          rejectionReason: 'duplicate charge',
        },
      ],
    };

    const wrapper = mount(
      <MockProviders initialEntries={[`/moves/${testMoveLocator}/payment-requests`]}>
        <PaymentRequestCard paymentRequest={reviewedPaymentRequest} />
      </MockProviders>,
    );

    it('renders the reviewed status tag', () => {
      expect(wrapper.find({ 'data-testid': 'tag' }).contains('REVIEWED')).toBe(true);
    });

    it('sums the approved service items total', () => {
      expect(wrapper.find('.amountAccepted h2').contains('$20,000.01')).toBe(true);
    });

    it('sums the rejected service items total', () => {
      expect(wrapper.find('.amountRejected h2').contains('$40,000.01')).toBe(true);
    });

    it('renders the view documents link', () => {
      const viewLink = wrapper.find('.footer a');

      expect(viewLink.text()).toEqual('View documents');
      expect(viewLink.prop('href')).toBe(`payment-requests/${reviewedPaymentRequest.id}`);
    });

    it('shows only rejected if no service items are approved', () => {
      const rejected = mount(
        <MockProviders initialEntries={[`/moves/${testMoveLocator}/payment-requests`]}>
          <PaymentRequestCard paymentRequest={rejectedPaymentRequest} />
        </MockProviders>,
      );

      expect(rejected.find('.amountRejected h2').contains('$60,000.02')).toBe(true);
      expect(rejected.find('.amountAccepted').exists()).toBe(false);
    });
  });
});
