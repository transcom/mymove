import React from 'react';

import PaymentRequestDetails from './PaymentRequestDetails';

import { PAYMENT_SERVICE_ITEM_STATUS, SHIPMENT_OPTIONS } from 'shared/constants';
import { formatPaymentRequestAddressString } from 'utils/shipmentDisplay';
import { shipmentModificationTypes } from 'constants/shipments';

export default {
  title: 'Office Components/PaymentRequestDetails',
  decorators: [
    (Story) => {
      return (
        <div style={{ padding: '20px' }}>
          <Story />
        </div>
      );
    },
  ],
};

const hhgAddress = formatPaymentRequestAddressString(
  {
    city: 'Beverly Hills',
    postal_code: '90210',
    state: 'CA',
  },
  {
    city: 'Fairfield',
    postal_code: '94535',
    state: 'CA',
  },
);

const ntsAddress = formatPaymentRequestAddressString(
  {
    city: 'Boston',
    postal_code: '02101',
    state: 'MA',
  },
  {
    city: 'Princeton',
    postal_code: '08540',
    state: 'NJ',
  },
);

const unreviewedPaymentRequestItems = [
  {
    id: '09474c6a-69b6-4501-8e08-670a12512a5f',
    createdAt: '2020-12-01T00:00:00.000Z',
    mtoServiceItemID: 'f8c2f97f-99e7-4fb1-9cc4-473debd24dbc',
    priceCents: 2000001,
    status: PAYMENT_SERVICE_ITEM_STATUS.REQUESTED,
    mtoShipmentType: null,
    mtoServiceItemName: 'Move management',
  },
  {
    id: '39474c6a-69b6-4501-8e08-670a12512a5f',
    createdAt: '2020-12-01T00:00:00.000Z',
    mtoServiceItemID: 'a8c2f97f-99e7-4fb1-9cc4-473debd24dbc',
    priceCents: 4000001,
    status: PAYMENT_SERVICE_ITEM_STATUS.REQUESTED,
    rejectionReason: 'duplicate charge',
    mtoShipmentType: null,
    mtoServiceItemName: 'Counseling',
  },
];

const reviewedPaymentRequestItems = [
  {
    id: '09474c6a-69b6-4501-8e08-670a12512a5f',
    createdAt: '2020-12-01T00:00:00.000Z',
    mtoServiceItemID: 'f8c2f97f-99e7-4fb1-9cc4-473debd24dbc',
    priceCents: 2000001,
    status: PAYMENT_SERVICE_ITEM_STATUS.APPROVED,
    mtoShipmentType: null,
    mtoServiceItemName: 'Move management',
  },
  {
    id: '39474c6a-69b6-4501-8e08-670a12512a5f',
    createdAt: '2020-12-01T00:00:00.000Z',
    mtoServiceItemID: 'a8c2f97f-99e7-4fb1-9cc4-473debd24dbc',
    priceCents: 4000001,
    status: PAYMENT_SERVICE_ITEM_STATUS.DENIED,
    rejectionReason: 'duplicate charge',
    mtoShipmentType: null,
    mtoServiceItemName: 'Counseling',
  },
];

const singleBasicServiceItem = [
  {
    id: '09474c6a-69b6-4501-8e08-670a12512a5f',
    createdAt: '2020-12-01T00:00:00.000Z',
    mtoServiceItemID: 'f8c2f97f-99e7-4fb1-9cc4-473debd24dbc',
    priceCents: 2000001,
    status: PAYMENT_SERVICE_ITEM_STATUS.APPROVED,
    mtoShipmentType: null,
    mtoServiceItemName: 'Move management',
  },
];

const hhgRequestedServiceItems = [
  {
    id: '09474c6a-69b6-4501-8e08-670a12512a5a',
    createdAt: '2020-12-01T00:04:00.000Z',
    mtoServiceItemID: 'f8c2f97f-99e7-4fb1-9cc4-473debd24dba',
    mtoShipmentID: 'a8c2f97f-99e7-4fb1-9cc4-473debd24dba',
    priceCents: 2000001,
    status: PAYMENT_SERVICE_ITEM_STATUS.REQUESTED,
    mtoShipmentType: SHIPMENT_OPTIONS.HHG,
    mtoServiceItemName: 'Domestic linehaul',
  },
  {
    id: '09474c6a-69b6-4501-8e08-670a12512a5b',
    createdAt: '2020-12-01T00:05:00.000Z',
    mtoServiceItemID: 'f8c2f97f-99e7-4fb1-9cc4-473debd24dbb',
    mtoShipmentID: 'a8c2f97f-99e7-4fb1-9cc4-473debd24dba',
    priceCents: 2000001,
    status: PAYMENT_SERVICE_ITEM_STATUS.REQUESTED,
    mtoShipmentType: SHIPMENT_OPTIONS.HHG,
    mtoServiceItemName: 'Fuel surcharge',
  },
  {
    id: '09474c6a-69b6-4501-8e08-670a12512a5c',
    createdAt: '2020-12-01T00:06:00.000Z',
    mtoServiceItemID: 'f8c2f97f-99e7-4fb1-9cc4-473debd24dbc',
    mtoShipmentID: 'a8c2f97f-99e7-4fb1-9cc4-473debd24dba',
    priceCents: 2000001,
    status: PAYMENT_SERVICE_ITEM_STATUS.APPROVED,
    mtoShipmentType: SHIPMENT_OPTIONS.HHG,
    mtoServiceItemName: 'Domestic origin price',
  },
  {
    id: '09474c6a-69b6-4501-8e08-670a12512a5d',
    createdAt: '2020-12-01T00:07:00.000Z',
    mtoServiceItemID: 'f8c2f97f-99e7-4fb1-9cc4-473debd24dbd',
    mtoShipmentID: 'a8c2f97f-99e7-4fb1-9cc4-473debd24dba',
    priceCents: 2000001,
    status: PAYMENT_SERVICE_ITEM_STATUS.APPROVED,
    mtoShipmentType: SHIPMENT_OPTIONS.HHG,
    mtoServiceItemName: 'Domestic destination price',
  },
  {
    id: '09474c6a-69b6-4501-8e08-670a12512a5e',
    createdAt: '2020-12-01T00:08:00.000Z',
    mtoServiceItemID: 'f8c2f97f-99e7-4fb1-9cc4-473debd24dbe',
    mtoShipmentID: 'a8c2f97f-99e7-4fb1-9cc4-473debd24dba',
    priceCents: 2000001,
    status: PAYMENT_SERVICE_ITEM_STATUS.DENIED,
    mtoShipmentType: SHIPMENT_OPTIONS.HHG,
    mtoServiceItemName: 'Domestic packing',
  },
  {
    id: '09474c6a-69b6-4501-8e08-670a12512a5f',
    createdAt: '2020-12-01T00:09:00.000Z',
    mtoServiceItemID: 'f8c2f97f-99e7-4fb1-9cc4-473debd24dbf',
    mtoShipmentID: 'a8c2f97f-99e7-4fb1-9cc4-473debd24dba',
    priceCents: 2000001,
    status: PAYMENT_SERVICE_ITEM_STATUS.DENIED,
    mtoShipmentType: SHIPMENT_OPTIONS.HHG,
    mtoServiceItemName: 'Domestic unpacking',
  },
];

const ntsrRequestedServiceItems = [
  {
    id: '09474c6a-69b6-4501-8e08-670a12512a5a',
    createdAt: '2020-12-01T00:04:00.000Z',
    mtoServiceItemID: 'f8c2f97f-99e7-4fb1-9cc4-473debd24dba',
    mtoShipmentID: 'a8c2f97f-99e7-4fb1-9cc4-473debd24dba',
    priceCents: 2000001,
    status: PAYMENT_SERVICE_ITEM_STATUS.APPROVED,
    mtoShipmentType: SHIPMENT_OPTIONS.NTSR,
    mtoServiceItemName: 'Domestic linehaul',
  },
  {
    id: '09474c6a-69b6-4501-8e08-670a12512a5b',
    createdAt: '2020-12-01T00:05:00.000Z',
    mtoServiceItemID: 'f8c2f97f-99e7-4fb1-9cc4-473debd24dbb',
    mtoShipmentID: 'a8c2f97f-99e7-4fb1-9cc4-473debd24dba',
    priceCents: 2000001,
    status: PAYMENT_SERVICE_ITEM_STATUS.APPROVED,
    mtoShipmentType: SHIPMENT_OPTIONS.NTSR,
    mtoServiceItemName: 'Fuel surcharge',
  },
  {
    id: '09474c6a-69b6-4501-8e08-670a12512a5c',
    createdAt: '2020-12-01T00:06:00.000Z',
    mtoServiceItemID: 'f8c2f97f-99e7-4fb1-9cc4-473debd24dbc',
    mtoShipmentID: 'a8c2f97f-99e7-4fb1-9cc4-473debd24dba',
    priceCents: 2000001,
    status: PAYMENT_SERVICE_ITEM_STATUS.REQUESTED,
    mtoShipmentType: SHIPMENT_OPTIONS.NTSR,
    mtoServiceItemName: 'Domestic origin price',
  },
  {
    id: '09474c6a-69b6-4501-8e08-670a12512a5d',
    createdAt: '2020-12-01T00:07:00.000Z',
    mtoServiceItemID: 'f8c2f97f-99e7-4fb1-9cc4-473debd24dbd',
    mtoShipmentID: 'a8c2f97f-99e7-4fb1-9cc4-473debd24dba',
    priceCents: 2000001,
    status: PAYMENT_SERVICE_ITEM_STATUS.DENIED,
    mtoShipmentType: SHIPMENT_OPTIONS.NTSR,
    mtoServiceItemName: 'Domestic destination price',
  },
  {
    id: '09474c6a-69b6-4501-8e08-670a12512a5f',
    createdAt: '2020-12-01T00:09:00.000Z',
    mtoServiceItemID: 'f8c2f97f-99e7-4fb1-9cc4-473debd24dbf',
    mtoShipmentID: 'a8c2f97f-99e7-4fb1-9cc4-473debd24dba',
    priceCents: 2000001,
    status: PAYMENT_SERVICE_ITEM_STATUS.REQUESTED,
    mtoShipmentType: SHIPMENT_OPTIONS.NTSR,
    mtoServiceItemName: 'Domestic unpacking',
  },
];

const hhgShipment = {
  address: hhgAddress,
  departureDate: '2020-12-01T00:00:00.000Z',
};

const hhgShipmentCanceled = {
  address: hhgAddress,
  departureDate: '2020-12-01T00:00:00.000Z',
  modificationType: shipmentModificationTypes.CANCELED,
};

const hhgShipmentDiversion = {
  address: hhgAddress,
  departureDate: '2020-12-01T00:00:00.000Z',
  modificationType: shipmentModificationTypes.DIVERSION,
};

const basicShipment = {
  address: '',
  departureDate: '',
};

const ntsShipment = {
  address: ntsAddress,
  departureDate: '020-12-01T00:00:00.000Z',
};

export const withUnreviewedBasicServiceItems = () => (
  <PaymentRequestDetails serviceItems={unreviewedPaymentRequestItems} />
);
export const withReviewedBasicServiceItems = () => (
  <PaymentRequestDetails serviceItems={reviewedPaymentRequestItems} shipment={basicShipment} />
);
export const withSingleBasicServiceItem = () => (
  <PaymentRequestDetails serviceItems={singleBasicServiceItem} shipment={basicShipment} />
);

export const withHHGShipmentServiceItems = () => (
  <PaymentRequestDetails
    shipmentDepartureDate="2021-04-20"
    serviceItems={hhgRequestedServiceItems}
    shipment={hhgShipment}
  />
);

export const withHHGShipmentServiceItemsWithACanceledShipment = () => (
  <PaymentRequestDetails
    shipmentDepartureDate="2021-04-20"
    serviceItems={hhgRequestedServiceItems}
    shipment={hhgShipmentCanceled}
  />
);

export const withHHGShipmentServiceItemsWithADivertedShipment = () => (
  <PaymentRequestDetails
    shipmentDepartureDate="2021-04-20"
    serviceItems={hhgRequestedServiceItems}
    shipment={hhgShipmentDiversion}
  />
);

export const withNTSRShipmentServiceItems = () => (
  <PaymentRequestDetails serviceItems={ntsrRequestedServiceItems} shipment={ntsShipment} />
);
