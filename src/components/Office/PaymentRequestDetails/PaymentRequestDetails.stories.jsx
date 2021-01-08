import React from 'react';

import PaymentRequestDetails from './PaymentRequestDetails';

import { PAYMENT_SERVICE_ITEM_STATUS } from 'shared/constants';

const unreviewedPaymentRequest = {
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
      status: PAYMENT_SERVICE_ITEM_STATUS.REQUESTED,
      shipmentType: null,
      serviceItemName: 'Move management',
    },
    {
      id: '39474c6a-69b6-4501-8e08-670a12512a5f',
      createdAt: '2020-12-01T00:00:00.000Z',
      mtoServiceItemID: 'a8c2f97f-99e7-4fb1-9cc4-473debd24dbc',
      priceCents: 4000001,
      status: PAYMENT_SERVICE_ITEM_STATUS.REQUESTED,
      rejectionReason: 'duplicate charge',
      shipmentType: null,
      serviceItemName: 'Counseling',
    },
  ],
};
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
      status: PAYMENT_SERVICE_ITEM_STATUS.APPROVED,
      shipmentType: null,
      serviceItemName: 'Move management',
    },
    {
      id: '39474c6a-69b6-4501-8e08-670a12512a5f',
      createdAt: '2020-12-01T00:00:00.000Z',
      mtoServiceItemID: 'a8c2f97f-99e7-4fb1-9cc4-473debd24dbc',
      priceCents: 4000001,
      status: PAYMENT_SERVICE_ITEM_STATUS.DENIED,
      rejectionReason: 'duplicate charge',
      shipmentType: null,
      serviceItemName: 'Counseling',
    },
  ],
};
const singleServicePaymentRequest = {
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
      status: PAYMENT_SERVICE_ITEM_STATUS.APPROVED,
      shipmentType: null,
      serviceItemName: 'Move management',
    },
  ],
};

export default {
  title: 'Office Components/PaymentRequestDetails',
};

export const withUnreviewedBasicServiceItems = () => (
  <div style={{ padding: '20px' }}>
    <PaymentRequestDetails serviceItems={unreviewedPaymentRequest.serviceItems} />
  </div>
);
export const withReviewedBasicServiceItems = () => (
  <div style={{ padding: '20px' }}>
    <PaymentRequestDetails serviceItems={reviewedPaymentRequest.serviceItems} />
  </div>
);
export const withSingleBasicServiceItem = () => (
  <div style={{ padding: '20px' }}>
    <PaymentRequestDetails serviceItems={singleServicePaymentRequest.serviceItems} />
  </div>
);
