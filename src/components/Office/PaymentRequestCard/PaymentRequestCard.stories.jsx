import React from 'react';
import moment from 'moment';
import MockDate from 'mockdate';
import addons from '@storybook/addons';
import { isHappoRun } from 'happo-plugin-storybook/register';

import PaymentRequestCard from './PaymentRequestCard';

import { MockProviders } from 'testUtils';

const mockedDate = '2020-12-08T00:00:00.000Z';

export default {
  title: 'Office Components/PaymentRequestCard',
  component: PaymentRequestCard,
  decorators: [
    (Story) => {
      if (isHappoRun()) {
        MockDate.set(mockedDate);
        addons.getChannel().on('storyRendered', MockDate.reset);
      }
      return (
        <div style={{ padding: '1em', backgroundColor: '#f9f9f9' }}>
          <MockProviders initialEntries={['/moves/L0CATR/payment-requests']}>
            <Story />
          </MockProviders>
        </div>
      );
    },
  ],
};

// always show 7 days prior to mocked date time
const itsBeenOneWeek = moment(mockedDate).subtract(7, 'days').format('YYYY-MM-DDTHH:mm:ss.SSSZ');

const pendingPaymentRequest = {
  id: '09474c6a-69b6-4501-8e08-670a12512e5f',
  createdAt: isHappoRun() ? itsBeenOneWeek : '2020-12-01T00:00:00.000Z',
  moveTaskOrderID: 'f8c2f97f-99e7-4fb1-9cc4-473debd04dbc',
  paymentRequestNumber: '1843-9061-1',
  status: 'PENDING',
  serviceItems: [
    {
      id: '09474c6a-69b6-4501-8e08-670a12512a5f',
      createdAt: '2020-12-01T00:00:00.000Z',
      mtoServiceItemID: 'f8c2f97f-99e7-4fb1-9cc4-473debd24dbc',
      priceCents: 1000001,
      status: 'REQUESTED',
    },
  ],
  reviewedAt: null,
};

const reviewedPaymentRequest = {
  id: '09474c6a-69b6-4501-8e08-670a12512e5f',
  createdAt: isHappoRun() ? itsBeenOneWeek : '2020-12-01T00:00:00.000Z',
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
};

const rejectedPaymentRequest = {
  id: '09474c6a-69b6-4501-8e08-670a12512e5f',
  createdAt: isHappoRun() ? itsBeenOneWeek : '2020-12-01T00:00:00.000Z',
  moveTaskOrderID: 'f8c2f97f-99e7-4fb1-9cc4-473debd04dbc',
  paymentRequestNumber: '1843-9061-1',
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
      rejectionReason: 'Requested amount exceeds guideline',
    },
  ],
  reviewedAt: '2020-12-01T00:00:00.000Z',
};

export const NeedsReview = () => <PaymentRequestCard paymentRequest={pendingPaymentRequest} />;

export const Reviewed = () => <PaymentRequestCard paymentRequest={reviewedPaymentRequest} />;

export const Rejected = () => <PaymentRequestCard paymentRequest={rejectedPaymentRequest} />;
