import React from 'react';
import { mount } from 'enzyme';

import PaymentRequestDetails from './PaymentRequestDetails';

import { MockProviders } from 'testUtils';

const basicPaymentRequest = {
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
      shipmentType: null,
      serviceItemName: 'Move management',
    },
    {
      id: '39474c6a-69b6-4501-8e08-670a12512a5f',
      createdAt: '2020-12-01T00:00:00.000Z',
      mtoServiceItemID: 'a8c2f97f-99e7-4fb1-9cc4-473debd24dbc',
      priceCents: 4000001,
      status: 'DENIED',
      rejectionReason: 'duplicate charge',
      shipmentType: null,
      serviceItemName: 'Counseling',
    },
    {
      id: '39474c6a-69b6-4501-8e08-670a12512a5f',
      createdAt: '2020-12-01T00:00:00.000Z',
      mtoServiceItemID: 'a8c2f97f-99e7-4fb1-9cc4-473debd24dbc',
      priceCents: 4000001,
      status: 'PENDING',
      rejectionReason: 'duplicate charge',
      shipmentType: null,
      serviceItemName: 'Counseling',
    },
    {
      id: '09474c6a-69b6-4501-8e08-670a12512a5f',
      createdAt: '2020-12-01T00:00:00.000Z',
      mtoServiceItemID: 'f8c2f97f-99e7-4fb1-9cc4-473debd24dbc',
      priceCents: 2000001,
      status: 'APPROVED',
      shipmentType: null,
      serviceItemName: 'Move management',
    },
    {
      id: '39474c6a-69b6-4501-8e08-670a12512a5f',
      createdAt: '2020-12-01T00:00:00.000Z',
      mtoServiceItemID: 'a8c2f97f-99e7-4fb1-9cc4-473debd24dbc',
      priceCents: 4000001,
      status: 'DENIED',
      rejectionReason: 'duplicate charge',
      shipmentType: null,
      serviceItemName: 'Counseling',
    },
    {
      id: '39474c6a-69b6-4501-8e08-670a12512a5f',
      createdAt: '2020-12-01T00:00:00.000Z',
      mtoServiceItemID: 'a8c2f97f-99e7-4fb1-9cc4-473debd24dbc',
      priceCents: 4000001,
      status: 'PENDING',
      rejectionReason: 'duplicate charge',
      shipmentType: null,
      serviceItemName: 'Counseling',
    },
  ],
};

const basicPaymentRequestOneServiceItem = {
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
      shipmentType: null,
      serviceItemName: 'Move management',
    },
  ],
};

const testMoveLocator = 'AF7K1P';

describe('PaymentRequestDetails', () => {
  describe('When given basic service items', () => {
    const wrapper = mount(
      <MockProviders initialEntries={[`/moves/${testMoveLocator}/payment-requests`]}>
        <PaymentRequestDetails serviceItems={basicPaymentRequest.serviceItems} />
      </MockProviders>,
    );

    it('renders the service items', async () => {
      expect(wrapper.find('td')).toBeTruthy();
    });

    it('renders the expected table title', () => {
      expect(wrapper.text().includes('Basic Service Items (6 items)')).toBeTruthy();
    });
  });

  describe('When given a single basic service item', () => {
    const wrapper = mount(
      <MockProviders initialEntries={[`/moves/${testMoveLocator}/payment-requests`]}>
        <PaymentRequestDetails serviceItems={basicPaymentRequestOneServiceItem.serviceItems} />
      </MockProviders>,
    );

    it('renders the expected table title', () => {
      expect(wrapper.text().includes('Basic Service Items (1 item)')).toBeTruthy();
    });
  });
});
