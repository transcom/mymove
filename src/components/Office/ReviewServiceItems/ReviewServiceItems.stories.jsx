import React from 'react';
import { action } from '@storybook/addon-actions';

import ReviewServiceItems from './ReviewServiceItems';

import {
  SHIPMENT_OPTIONS,
  SERVICE_ITEM_STATUS,
  PAYMENT_SERVICE_ITEM_STATUS,
  PAYMENT_REQUEST_STATUS,
} from 'shared/constants';

export default {
  title: 'Office Components/ReviewServiceItems',
  component: ReviewServiceItems,
  decorators: [
    (storyFn) => (
      <div style={{ margin: '10px', height: '80vh', display: 'flex', flexDirection: 'column' }}>{storyFn()}</div>
    ),
  ],
};

const pendingPaymentRequest = { status: PAYMENT_REQUEST_STATUS.PENDING };
export const Basic = () => (
  <ReviewServiceItems
    paymentRequest={pendingPaymentRequest}
    serviceItemCards={[
      {
        id: '1',
        mtoServiceItemName: 'Counseling services',
        amount: 1234.0,
        createdAt: '2020-01-01T00:08:00.999Z',
      },
    ]}
    handleClose={action('clicked')}
    onCompleteReview={action('clicked')}
    patchPaymentServiceItem={action('patchPaymentServiceItem')}
  />
);

export const BasicWithTwoItems = () => {
  return (
    <ReviewServiceItems
      disableScrollIntoView
      paymentRequest={pendingPaymentRequest}
      serviceItemCards={[
        {
          id: '1',
          mtoServiceItemName: 'Counseling services',
          amount: 1234.0,
          createdAt: '2020-01-01T00:08:00.999Z',
        },
        {
          id: '2',
          mtoServiceItemName: 'Move management',
          amount: 1234.0,
          createdAt: '2020-01-01T00:08:00.999Z',
        },
      ]}
      handleClose={action('clicked')}
      onCompleteReview={action('clicked')}
      patchPaymentServiceItem={action('patchPaymentServiceItem')}
    />
  );
};

export const HHG = () => (
  <ReviewServiceItems
    paymentRequest={pendingPaymentRequest}
    serviceItemCards={[
      {
        id: '1',
        mtoShipmentID: '10',
        mtoShipmentType: SHIPMENT_OPTIONS.HHG,
        mtoServiceItemName: 'Domestic linehaul',
        amount: 5678.05,
        createdAt: '2020-01-01T00:08:00.999Z',
      },
    ]}
    handleClose={action('clicked')}
    onCompleteReview={action('clicked')}
    patchPaymentServiceItem={action('patchPaymentServiceItem')}
  />
);

export const NonTemporaryStorage = () => (
  <ReviewServiceItems
    paymentRequest={pendingPaymentRequest}
    serviceItemCards={[
      {
        id: '1',
        mtoShipmentID: '10',
        mtoShipmentType: SHIPMENT_OPTIONS.NTS,
        mtoServiceItemName: 'Domestic linehaul',
        amount: 6423.51,
        createdAt: '2020-01-01T00:08:00.999Z',
      },
    ]}
    handleClose={action('clicked')}
    onCompleteReview={action('clicked')}
    patchPaymentServiceItem={action('patchPaymentServiceItem')}
  />
);

export const MultipleShipmentsGroups = () => (
  <ReviewServiceItems
    paymentRequest={pendingPaymentRequest}
    serviceItemCards={[
      {
        id: '1',
        mtoServiceItemName: 'Counseling services',
        amount: 0.01,
        createdAt: '2020-01-01T00:09:00.999Z',
      },
      {
        id: '3',
        mtoShipmentID: '20',
        mtoShipmentType: SHIPMENT_OPTIONS.HHG,
        mtoServiceItemName: 'Domestic linehaul',
        amount: 5678.05,
        createdAt: '2020-01-01T00:08:00.999Z',
      },
      {
        id: '4',
        mtoShipmentID: '30',
        mtoShipmentType: SHIPMENT_OPTIONS.NTS,
        mtoServiceItemName: 'Domestic linehaul',
        amount: 6423.51,
        createdAt: '2020-01-01T00:07:30.999Z',
      },
      {
        id: '5',
        mtoShipmentID: '30',
        mtoShipmentType: SHIPMENT_OPTIONS.NTS,
        mtoServiceItemName: 'Fuel Surcharge',
        amount: 100000000000000,
        createdAt: '2020-01-01T00:07:00.999Z',
      },
    ]}
    handleClose={action('clicked')}
    onCompleteReview={action('clicked')}
    patchPaymentServiceItem={action('patchPaymentServiceItem')}
  />
);

export const WithStatusAndReason = () => (
  <ReviewServiceItems
    disableScrollIntoView
    paymentRequest={pendingPaymentRequest}
    serviceItemCards={[
      {
        id: '1',
        mtoServiceItemName: 'Counseling services',
        amount: 0.01,
        createdAt: '2020-01-01T00:09:00.999Z',
      },
      {
        id: '2',
        mtoServiceItemName: 'Move management',
        amount: 1234.0,
        status: SERVICE_ITEM_STATUS.REJECTED,
        rejectionReason: 'Amount exceeds limit',
        createdAt: '2020-01-01T00:06:00.999Z',
      },
      {
        id: '3',
        mtoShipmentID: '20',
        mtoShipmentType: SHIPMENT_OPTIONS.HHG,
        mtoServiceItemName: 'Domestic linehaul',
        amount: 5678.05,
        status: SERVICE_ITEM_STATUS.APPROVED,
        createdAt: '2020-01-01T00:08:00.999Z',
      },
      {
        id: '4',
        mtoShipmentID: '30',
        mtoShipmentType: SHIPMENT_OPTIONS.NTS,
        mtoServiceItemName: 'Domestic linehaul',
        amount: 6423.51,
        status: SERVICE_ITEM_STATUS.APPROVED,
        createdAt: '2020-01-01T00:07:30.999Z',
      },
      {
        id: '5',
        mtoShipmentID: '30',
        mtoShipmentType: SHIPMENT_OPTIONS.NTS,
        mtoServiceItemName: 'Fuel Surcharge',
        amount: 100000000000000,
        createdAt: '2020-01-01T00:07:00.999Z',
      },
    ]}
    handleClose={action('clicked')}
    onCompleteReview={action('clicked')}
    patchPaymentServiceItem={action('patchPaymentServiceItem')}
  />
);

export const WithNeedsReview = () => (
  <ReviewServiceItems
    disableScrollIntoView
    paymentRequest={pendingPaymentRequest}
    serviceItemCards={[
      {
        id: '1',
        mtoServiceItemName: 'Counseling services',
        status: PAYMENT_SERVICE_ITEM_STATUS.REQUESTED,
        amount: 0.01,
        createdAt: '2020-01-01T00:09:00.999Z',
      },
    ]}
    handleClose={action('clicked')}
    onCompleteReview={action('clicked')}
    patchPaymentServiceItem={action('patchPaymentServiceItem')}
  />
);

export const WithRejectRequest = () => (
  <ReviewServiceItems
    disableScrollIntoView
    paymentRequest={pendingPaymentRequest}
    serviceItemCards={[
      {
        id: '1',
        mtoServiceItemName: 'Counseling services',
        status: PAYMENT_SERVICE_ITEM_STATUS.DENIED,
        amount: 0.01,
        createdAt: '2020-01-01T00:09:00.999Z',
      },
    ]}
    handleClose={action('clicked')}
    onCompleteReview={action('clicked')}
    patchPaymentServiceItem={action('patchPaymentServiceItem')}
  />
);

export const WithAuthorizePayment = () => (
  <ReviewServiceItems
    disableScrollIntoView
    paymentRequest={pendingPaymentRequest}
    serviceItemCards={[
      {
        id: '1',
        mtoServiceItemName: 'Counseling services',
        status: PAYMENT_SERVICE_ITEM_STATUS.APPROVED,
        amount: 0.01,
        createdAt: '2020-01-01T00:09:00.999Z',
      },
    ]}
    handleClose={action('clicked')}
    onCompleteReview={action('clicked')}
    patchPaymentServiceItem={action('patchPaymentServiceItem')}
  />
);

export const WithPaymentReviewedApproved = () => (
  <ReviewServiceItems
    disableScrollIntoView
    paymentRequest={{
      status: PAYMENT_REQUEST_STATUS.REVIEWED,
      reviewedAt: '2020-08-31T20:30:59.000Z',
    }}
    serviceItemCards={[
      {
        id: '1',
        mtoServiceItemName: 'Counseling services',
        status: PAYMENT_SERVICE_ITEM_STATUS.APPROVED,
        amount: 0.01,
        createdAt: '2020-01-01T00:09:00.999Z',
      },
    ]}
    handleClose={action('clicked')}
    onCompleteReview={action('clicked')}
    patchPaymentServiceItem={action('patchPaymentServiceItem')}
  />
);

export const WithPaymentReviewedRejected = () => (
  <ReviewServiceItems
    disableScrollIntoView
    paymentRequest={{
      status: PAYMENT_REQUEST_STATUS.REVIEWED,
    }}
    serviceItemCards={[
      {
        id: '1',
        mtoServiceItemName: 'Counseling services',
        status: PAYMENT_SERVICE_ITEM_STATUS.DENIED,
        rejectionReason: 'Service member already counseled',
        amount: 0.01,
        createdAt: '2020-01-01T00:09:00.999Z',
      },
    ]}
    handleClose={action('clicked')}
    onCompleteReview={action('clicked')}
    patchPaymentServiceItem={action('patchPaymentServiceItem')}
  />
);
