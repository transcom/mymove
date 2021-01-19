import React from 'react';
import { mount } from 'enzyme';

import PaymentRequestCard from './PaymentRequestCard';

import { MockProviders } from 'testUtils';

jest.mock('hooks/queries', () => ({
  useMovePaymentRequestsQueries: () => {
    const order = {
      sac: '1234456',
      tac: '1213',
    };

    const contractor = {
      contractNumber: 'HTC-123-3456',
    };

    const move = {
      contractor,
      orders: order,
    };

    return {
      paymentRequests: [
        {
          id: '09474c6a-69b6-4501-8e08-670a12512e5f',
          createdAt: '2020-12-01T00:00:00.000Z',
          moveTaskOrderID: 'f8c2f97f-99e7-4fb1-9cc4-473debd04dbc',
          paymentRequestNumber: '1843-9061-1',
          status: 'REVIEWED',
          reviewedAt: '2020-12-01T00:00:00.000Z',
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
    };
  },
}));

const testMoveLocator = 'AF7K1P';

describe('PaymentRequestCard', () => {
  const order = {
    sac: '1234456',
    tac: '1213',
  };

  const contractor = {
    contractNumber: 'HTC-123-3456',
  };

  const move = {
    contractor,
    orders: order,
  };
  describe('pending payment request', () => {
    const pendingPaymentRequest = {
      id: '29474c6a-69b6-4501-8e08-670a12512e5f',
      createdAt: '2020-12-01T00:00:00.000Z',
      moveTaskOrderID: 'f8c2f97f-99e7-4fb1-9cc4-473debd04dbc',
      paymentRequestNumber: '1843-9061-2',
      moveTaskOrder: move,
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
      expect(wrapper.find({ 'data-testid': 'tag' }).contains('Needs review')).toBe(true);
    });

    it('sums the service items total', () => {
      expect(wrapper.find('.amountRequested').contains('$60,000.02')).toBe(true);
    });

    it('displays the payment request details ', () => {
      const prDetails = wrapper.find('.footer dd');
      expect(prDetails.contains(order.sac)).toBe(true);
      expect(prDetails.contains(order.tac)).toBe(true);
      expect(prDetails.contains(contractor.contractNumber)).toBe(true);
    });

    it('renders the view orders link', () => {
      const viewLink = wrapper.find('.footer a');

      expect(viewLink.contains('View orders')).toBe(true);
      expect(viewLink.prop('href')).toBe('orders');
    });

    it('renders request details toggle drawer by default and hides button', () => {
      const showRequestDetailsButton = wrapper.find('button[data-testid="showRequestDetailsButton"]');

      expect(showRequestDetailsButton.length).toBe(0);
      expect(wrapper.find('[data-testid="toggleDrawer"]').length).toBe(1);
    });
  });

  describe('reviewed payment request', () => {
    const reviewedPaymentRequest = {
      id: '29474c6a-69b6-4501-8e08-670a12512e5f',
      createdAt: '2020-12-01T00:00:00.000Z',
      moveTaskOrderID: 'f8c2f97f-99e7-4fb1-9cc4-473debd04dbc',
      paymentRequestNumber: '1843-9061-2',
      status: 'REVIEWED_AND_ALL_SERVICE_ITEMS_REJECTED',
      moveTaskOrder: move,
      reviewedAt: '2020-12-01T00:00:00.000Z',
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
      reviewedAt: '2020-12-01T00:00:00.000Z',
      moveTaskOrder: move,
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

    it('renders the rejected status tag', () => {
      expect(wrapper.find({ 'data-testid': 'tag' }).contains('Rejected')).toBe(true);
    });

    it('sums the approved service items total', () => {
      expect(wrapper.find('.amountAccepted h2').contains('$20,000.01')).toBe(true);
    });

    it('displays the reviewed at date', () => {
      expect(wrapper.find('.amountAccepted span').at(1).text().includes('01 Dec 2020')).toBe(true);
    });

    it('sums the rejected service items total', () => {
      expect(wrapper.find('.amountRejected h2').contains('$40,000.01')).toBe(true);
    });

    it('displays the reviewed at date', () => {
      expect(wrapper.find('.amountRejected span').at(1).text().includes('01 Dec 2020')).toBe(true);
    });

    it('displays the payment request details ', () => {
      const prDetails = wrapper.find('.footer dd');
      expect(prDetails.contains(order.sac)).toBe(true);
      expect(prDetails.contains(order.tac)).toBe(true);
      expect(prDetails.contains(contractor.contractNumber)).toBe(true);
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

    it('renders request details toggle drawer after click', () => {
      const showRequestDetailsButton = wrapper.find('button[data-testid="showRequestDetailsButton"]');
      showRequestDetailsButton.simulate('click');

      expect(wrapper.find('[data-testid="toggleDrawer"]').length).toBe(1);
    });
  });

  describe('payment request gex statuses', () => {
    it('renders the reviewed status tag for sent_to_gex', () => {
      const sentToGexPaymentRequest = {
        id: '29474c6a-69b6-4501-8e08-670a12512e5f',
        createdAt: '2020-12-01T00:00:00.000Z',
        moveTaskOrderID: 'f8c2f97f-99e7-4fb1-9cc4-473debd04dbc',
        paymentRequestNumber: '1843-9061-2',
        status: 'SENT_TO_GEX',
        moveTaskOrder: move,
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
      const sentToGex = mount(
        <MockProviders initialEntries={[`/moves/${testMoveLocator}/payment-requests`]}>
          <PaymentRequestCard paymentRequest={sentToGexPaymentRequest} />
        </MockProviders>,
      );
      expect(sentToGex.find({ 'data-testid': 'tag' }).contains('Reviewed')).toBe(true);
    });

    it('renders the reviewed status tag for received_by_gex', () => {
      const receivedByGexPaymentRequest = {
        id: '29474c6a-69b6-4501-8e08-670a12512e5f',
        createdAt: '2020-12-01T00:00:00.000Z',
        moveTaskOrderID: 'f8c2f97f-99e7-4fb1-9cc4-473debd04dbc',
        paymentRequestNumber: '1843-9061-2',
        status: 'RECEIVED_BY_GEX',
        moveTaskOrder: move,
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
      const receivedByGex = mount(
        <MockProviders initialEntries={[`/moves/${testMoveLocator}/payment-requests`]}>
          <PaymentRequestCard paymentRequest={receivedByGexPaymentRequest} />
        </MockProviders>,
      );
      expect(receivedByGex.find({ 'data-testid': 'tag' }).contains('Reviewed')).toBe(true);
    });

    it('renders the paid status tag for paid request', () => {
      const paidPaymentRequest = {
        id: '29474c6a-69b6-4501-8e08-670a12512e5f',
        createdAt: '2020-12-01T00:00:00.000Z',
        moveTaskOrderID: 'f8c2f97f-99e7-4fb1-9cc4-473debd04dbc',
        paymentRequestNumber: '1843-9061-2',
        status: 'PAID',
        moveTaskOrder: move,
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
      const paid = mount(
        <MockProviders initialEntries={[`/moves/${testMoveLocator}/payment-requests`]}>
          <PaymentRequestCard paymentRequest={paidPaymentRequest} />
        </MockProviders>,
      );
      expect(paid.find({ 'data-testid': 'tag' }).contains('Paid')).toBe(true);
    });
  });
});
