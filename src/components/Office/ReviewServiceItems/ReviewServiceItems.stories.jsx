import React from 'react';
import { action } from '@storybook/addon-actions';

import testParams from '../ServiceItemCalculations/serviceItemTestParams';

import ReviewServiceItems from './ReviewServiceItems';

import {
  SHIPMENT_OPTIONS,
  SERVICE_ITEM_STATUS,
  PAYMENT_SERVICE_ITEM_STATUS,
  PAYMENT_REQUEST_STATUS,
  LOA_TYPE,
} from 'shared/constants';
import { serviceItemCodes } from 'content/serviceItems';

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
        mtoServiceItemName: serviceItemCodes.CS,
        amount: 1234.0,
        createdAt: '2020-01-01T00:08:00.999Z',
      },
    ]}
    handleClose={action('clicked')}
    onCompleteReview={action('clicked')}
    patchPaymentServiceItem={action('patchPaymentServiceItem')}
    TACs={{ HHG: '1234', NTS: '5678' }}
    SACs={{ HHG: 'AB12', NTS: 'CD34' }}
  />
);

export const HHG = () => (
  <ReviewServiceItems
    paymentRequest={pendingPaymentRequest}
    serviceItemCards={[
      {
        id: '1',
        mtoShipmentID: '10',
        mtoShipmentDepartureDate: '2020-04-29',
        mtoShipmentPickupAddress: 'Fairfield, CA 94535',
        mtoShipmentDestinationAddress: 'Beverly Hills, CA 90210',
        mtoShipmentType: SHIPMENT_OPTIONS.HHG,
        mtoShipmentTacType: LOA_TYPE.HHG,
        mtoShipmentSacType: LOA_TYPE.HHG,
        mtoServiceItemName: serviceItemCodes.DLH,
        mtoServiceItemCode: 'DLH',
        paymentServiceItemParams: testParams.DomesticLongHaul,
        amount: 5678.05,
        createdAt: '2020-01-01T00:08:00.999Z',
      },
    ]}
    handleClose={action('clicked')}
    onCompleteReview={action('clicked')}
    patchPaymentServiceItem={action('patchPaymentServiceItem')}
    TACs={{ HHG: '1234', NTS: '5678' }}
    SACs={{ HHG: 'AB12', NTS: 'CD34' }}
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
        mtoShipmentTacType: LOA_TYPE.NTS,
        mtoShipmentSacType: LOA_TYPE.NTS,
        mtoServiceItemName: serviceItemCodes.DLH,
        mtoServiceItemCode: 'DLH',
        paymentServiceItemParams: testParams.DomesticLongHaul,
        amount: 6423.51,
        createdAt: '2020-01-01T00:08:00.999Z',
      },
    ]}
    handleClose={action('clicked')}
    onCompleteReview={action('clicked')}
    patchPaymentServiceItem={action('patchPaymentServiceItem')}
    TACs={{ HHG: '1234', NTS: '5678' }}
    SACs={{ HHG: 'AB12', NTS: 'CD34' }}
  />
);

export const WithStatusAndReason = () => (
  <ReviewServiceItems
    paymentRequest={pendingPaymentRequest}
    serviceItemCards={[
      {
        id: '1',
        mtoServiceItemName: serviceItemCodes.CS,
        amount: 0.01,
        createdAt: '2020-01-01T00:09:00.999Z',
      },
      {
        id: '2',
        mtoServiceItemName: serviceItemCodes.MS,
        amount: 1234.0,
        status: SERVICE_ITEM_STATUS.REJECTED,
        rejectionReason: 'Amount exceeds limit',
        createdAt: '2020-01-01T00:06:00.999Z',
      },
      {
        id: '3',
        mtoShipmentID: '20',
        mtoShipmentDepartureDate: '04 May 2021',
        mtoShipmentPickupAddress: 'Fairfield, CA 94535',
        mtoShipmentDestinationAddress: 'Beverly Hills, CA 90210',
        mtoShipmentType: SHIPMENT_OPTIONS.HHG,
        mtoShipmentTacType: LOA_TYPE.HHG,
        mtoShipmentSacType: LOA_TYPE.HHG,
        mtoServiceItemName: serviceItemCodes.DLH,
        mtoServiceItemCode: 'DLH',
        paymentServiceItemParams: testParams.DomesticLongHaul,
        amount: 5678.05,
        status: SERVICE_ITEM_STATUS.APPROVED,
        createdAt: '2020-01-01T00:08:00.999Z',
      },
      {
        id: '4',
        mtoShipmentID: '30',
        mtoShipmentType: SHIPMENT_OPTIONS.NTSR,
        mtoShipmentTacType: LOA_TYPE.NTS,
        mtoShipmentSacType: LOA_TYPE.NTS,
        mtoServiceItemName: serviceItemCodes.DLH,
        mtoServiceItemCode: 'DLH',
        paymentServiceItemParams: testParams.DomesticLongHaul,
        amount: 6423.51,
        status: SERVICE_ITEM_STATUS.APPROVED,
        createdAt: '2020-01-01T00:07:30.999Z',
      },
      {
        id: '5',
        mtoShipmentID: '30',
        mtoShipmentType: SHIPMENT_OPTIONS.NTSR,
        mtoShipmentTacType: LOA_TYPE.NTS,
        mtoShipmentSacType: LOA_TYPE.NTS,
        mtoServiceItemName: serviceItemCodes.FSC,
        mtoServiceItemCode: 'FSC',
        paymentServiceItemParams: testParams.FuelSurchage,
        amount: 100000000000000,
        createdAt: '2020-01-01T00:07:00.999Z',
      },
    ]}
    handleClose={action('clicked')}
    onCompleteReview={action('clicked')}
    patchPaymentServiceItem={action('patchPaymentServiceItem')}
    TACs={{ HHG: '1234', NTS: '5678' }}
    SACs={{ HHG: 'AB12', NTS: 'CD34' }}
  />
);

export const WithNeedsReview = () => (
  <ReviewServiceItems
    paymentRequest={pendingPaymentRequest}
    serviceItemCards={[
      {
        id: '1',
        mtoServiceItemName: serviceItemCodes.CS,
        status: PAYMENT_SERVICE_ITEM_STATUS.REQUESTED,
        amount: 0.01,
        createdAt: '2020-01-01T00:09:00.999Z',
      },
    ]}
    handleClose={action('clicked')}
    onCompleteReview={action('clicked')}
    patchPaymentServiceItem={action('patchPaymentServiceItem')}
    TACs={{ HHG: '1234', NTS: '5678' }}
    SACs={{ HHG: 'AB12', NTS: 'CD34' }}
  />
);

export const WithRejectRequest = () => (
  <ReviewServiceItems
    paymentRequest={pendingPaymentRequest}
    serviceItemCards={[
      {
        id: '1',
        mtoServiceItemName: serviceItemCodes.CS,
        status: PAYMENT_SERVICE_ITEM_STATUS.DENIED,
        amount: 0.01,
        createdAt: '2020-01-01T00:09:00.999Z',
      },
    ]}
    handleClose={action('clicked')}
    onCompleteReview={action('clicked')}
    patchPaymentServiceItem={action('patchPaymentServiceItem')}
    TACs={{ HHG: '1234', NTS: '5678' }}
    SACs={{ HHG: 'AB12', NTS: 'CD34' }}
  />
);

export const WithAuthorizePayment = () => (
  <ReviewServiceItems
    paymentRequest={pendingPaymentRequest}
    serviceItemCards={[
      {
        id: '1',
        mtoServiceItemName: serviceItemCodes.CS,
        status: PAYMENT_SERVICE_ITEM_STATUS.APPROVED,
        amount: 0.01,
        createdAt: '2020-01-01T00:09:00.999Z',
      },
    ]}
    handleClose={action('clicked')}
    onCompleteReview={action('clicked')}
    patchPaymentServiceItem={action('patchPaymentServiceItem')}
    TACs={{ HHG: '1234', NTS: '5678' }}
    SACs={{ HHG: 'AB12', NTS: 'CD34' }}
  />
);

export const WithPaymentReviewedApproved = () => (
  <ReviewServiceItems
    paymentRequest={{
      status: PAYMENT_REQUEST_STATUS.REVIEWED,
      reviewedAt: '2020-08-31T20:30:59.000Z',
    }}
    serviceItemCards={[
      {
        id: '1',
        status: PAYMENT_SERVICE_ITEM_STATUS.APPROVED,
        mtoShipmentID: '10',
        mtoShipmentDepartureDate: '2020-04-29',
        mtoShipmentPickupAddress: 'Fairfield, CA 94535',
        mtoShipmentDestinationAddress: 'Beverly Hills, CA 90210',
        mtoShipmentType: SHIPMENT_OPTIONS.HHG,
        mtoShipmentTacType: LOA_TYPE.HHG,
        mtoShipmentSacType: LOA_TYPE.HHG,
        mtoServiceItemName: serviceItemCodes.DLH,
        mtoServiceItemCode: 'DLH',
        paymentServiceItemParams: testParams.DomesticLongHaul,
        amount: 5678.05,
        createdAt: '2020-01-01T00:08:00.999Z',
      },
      {
        amount: 54.89,
        status: PAYMENT_SERVICE_ITEM_STATUS.APPROVED,
        mtoServiceItemName: serviceItemCodes.MS,
      },
      {
        amount: 61.32,
        status: PAYMENT_SERVICE_ITEM_STATUS.APPROVED,
        mtoServiceItemName: serviceItemCodes.CS,
      },
    ]}
    curCardIndex={3}
    requestReviewed
    handleClose={action('clicked')}
    onCompleteReview={action('clicked')}
    patchPaymentServiceItem={action('patchPaymentServiceItem')}
    TACs={{ HHG: '1234', NTS: '5678' }}
    SACs={{ HHG: 'AB12', NTS: 'CD34' }}
  />
);

export const WithPaymentReviewedRejected = () => (
  <ReviewServiceItems
    paymentRequest={{
      status: PAYMENT_REQUEST_STATUS.REVIEWED,
    }}
    serviceItemCards={[
      {
        id: '1',
        mtoServiceItemName: serviceItemCodes.CS,
        status: PAYMENT_SERVICE_ITEM_STATUS.DENIED,
        rejectionReason: 'Service member already counseled',
        amount: 0.01,
        createdAt: '2020-01-01T00:09:00.999Z',
      },
    ]}
    curCardIndex={1}
    requestReviewed
    handleClose={action('clicked')}
    onCompleteReview={action('clicked')}
    patchPaymentServiceItem={action('patchPaymentServiceItem')}
    TACs={{ HHG: '1234', NTS: '5678' }}
    SACs={{ HHG: 'AB12', NTS: 'CD34' }}
  />
);
