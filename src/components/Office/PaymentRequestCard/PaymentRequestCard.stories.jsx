import React from 'react';
import moment from 'moment';
import MockDate from 'mockdate';
import addons from '@storybook/addons';
import { isHappoRun } from 'happo-plugin-storybook/register';

import { SHIPMENT_OPTIONS } from '../../../shared/constants';

import PaymentRequestCard from './PaymentRequestCard';

import { MockProviders } from 'testUtils';
import { serviceItemCodes } from 'content/serviceItems';
import { formatPaymentRequestAddressString } from 'utils/shipmentDisplay';

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
        <div style={{ padding: '1em', backgroundColor: '#f9f9f9', minWidth: '1200px' }}>
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

const order = {
  tac: '1213',
  // sac: 'CD7890',
  ntsTac: '4546',
  ntsSac: 'AB7890',
};

const contractor = {
  contractNumber: 'HTC-123-3456',
};

const move = {
  contractor,
  orders: order,
  locator: '12345',
};

const shipmentAddress = formatPaymentRequestAddressString(
  { city: 'Boston', state: 'MA', postalCode: '02101' },
  { city: 'Princeton', state: 'NJ', postalCode: '08540' },
);

const shipmentsInfo = [
  {
    mtoShipmentID: 'd81175b7-e26d-4e1e-b1d1-47b17bf4b7f3',
    departureDate: '2020-01-09T00:00:00.000Z',
    address: shipmentAddress,
    tacType: 'HHG',
    sacType: 'HHG',
  },
  {
    mtoShipmentID: '9e8222e4-9cdb-4994-8294-6d918a4c684d',
    departureDate: '2020-01-09T00:00:00.000Z',
    address: shipmentAddress,
    tacType: 'NTS',
    sacType: 'NTS',
  },
];

const pendingPaymentRequest = {
  id: '09474c6a-69b6-4501-8e08-670a12512e5f',
  createdAt: isHappoRun() ? itsBeenOneWeek : '2020-12-01T00:00:00.000Z',
  moveTaskOrderID: 'f8c2f97f-99e7-4fb1-9cc4-473debd04dbc',
  paymentRequestNumber: '1843-9061-1',
  status: 'PENDING',
  moveTaskOrder: move,
  serviceItems: [
    {
      id: '09474c6a-69b6-4501-8e08-670a12512a5f',
      createdAt: '2020-12-01T00:00:00.000Z',
      mtoServiceItemName: serviceItemCodes.CS,
      mtoServiceItemID: 'f8c2f97f-99e7-4fb1-9cc4-473debd24dbc',
      priceCents: 1000001,
      status: 'REQUESTED',
    },
    {
      id: '39474c6a-69b6-4501-8e08-670a12512a5f',
      createdAt: '2020-12-01T00:00:00.000Z',
      mtoServiceItemName: serviceItemCodes.DLH,
      mtoServiceItemID: 'a8c2f97f-99e7-4fb1-9cc4-473debd24dbc',
      mtoShipmentID: 'd81175b7-e26d-4e1e-b1d1-47b17bf4b7f3',
      mtoShipmentType: SHIPMENT_OPTIONS.HHG_LONGHAUL_DOMESTIC,
      priceCents: 4000001,
      status: 'REQUESTED',
    },
    {
      id: 'ad8b97ed-bb8a-4efa-abb3-2b00c849f537',
      createdAt: '2020-12-01T00:00:00.000Z',
      mtoServiceItemName: serviceItemCodes.FSC,
      mtoServiceItemID: 'a8c2f97f-99e7-4fb1-9cc4-473debd24dbb',
      mtoShipmentID: '9e8222e4-9cdb-4994-8294-6d918a4c684d',
      mtoShipmentType: SHIPMENT_OPTIONS.NTSR,
      priceCents: 6000001,
      status: 'REQUESTED',
    },
  ],
  reviewedAt: null,
};

const reviewedPaymentRequest = {
  id: '09474c6a-69b6-4501-8e08-670a12512e5f',
  createdAt: isHappoRun() ? itsBeenOneWeek : '2020-12-01T00:00:00.000Z',
  moveTaskOrderID: 'f8c2f97f-99e7-4fb1-9cc4-473debd04dbc',
  moveTaskOrder: move,
  paymentRequestNumber: '1843-9061-1',
  status: 'REVIEWED',
  serviceItems: [
    {
      id: '09474c6a-69b6-4501-8e08-670a12512a5f',
      createdAt: '2020-12-01T00:00:00.000Z',
      mtoServiceItemName: serviceItemCodes.CS,
      mtoServiceItemID: 'f8c2f97f-99e7-4fb1-9cc4-473debd24dbc',
      priceCents: 2000001,
      status: 'APPROVED',
    },
    {
      id: '39474c6a-69b6-4501-8e08-670a12512a5f',
      createdAt: '2020-12-01T00:00:00.000Z',
      mtoServiceItemName: serviceItemCodes.DLH,
      mtoServiceItemID: 'a8c2f97f-99e7-4fb1-9cc4-473debd24dbc',
      mtoShipmentID: 'd81175b7-e26d-4e1e-b1d1-47b17bf4b7f3',
      mtoShipmentType: SHIPMENT_OPTIONS.HHG_LONGHAUL_DOMESTIC,
      priceCents: 4000001,
      status: 'DENIED',
      rejectionReason: 'Requested amount exceeds guideline',
    },
    {
      id: 'ad8b97ed-bb8a-4efa-abb3-2b00c849f537',
      createdAt: '2020-12-01T00:00:00.000Z',
      mtoServiceItemName: serviceItemCodes.FSC,
      mtoServiceItemID: 'a8c2f97f-99e7-4fb1-9cc4-473debd24dbb',
      mtoShipmentID: '9e8222e4-9cdb-4994-8294-6d918a4c684d',
      mtoShipmentType: SHIPMENT_OPTIONS.NTSR,
      priceCents: 6000001,
      status: 'APPROVED',
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
  moveTaskOrder: move,
  serviceItems: [
    {
      id: '09474c6a-69b6-4501-8e08-670a12512a5f',
      createdAt: '2020-12-01T00:00:00.000Z',
      mtoServiceItemName: serviceItemCodes.CS,
      mtoServiceItemID: 'f8c2f97f-99e7-4fb1-9cc4-473debd24dbc',
      priceCents: 2000001,
      status: 'DENIED',
    },
    {
      id: '39474c6a-69b6-4501-8e08-670a12512a5f',
      createdAt: '2020-12-01T00:00:00.000Z',
      mtoServiceItemName: serviceItemCodes.DLH,
      mtoServiceItemID: 'a8c2f97f-99e7-4fb1-9cc4-473debd24dbc',
      mtoShipmentID: 'd81175b7-e26d-4e1e-b1d1-47b17bf4b7f3',
      mtoShipmentType: SHIPMENT_OPTIONS.HHG_LONGHAUL_DOMESTIC,
      priceCents: 4000001,
      status: 'DENIED',
      rejectionReason: 'Requested amount exceeds guideline',
    },
    {
      id: 'ad8b97ed-bb8a-4efa-abb3-2b00c849f537',
      createdAt: '2020-12-01T00:00:00.000Z',
      mtoServiceItemName: serviceItemCodes.FSC,
      mtoServiceItemID: 'a8c2f97f-99e7-4fb1-9cc4-473debd24dbb',
      mtoShipmentID: '9e8222e4-9cdb-4994-8294-6d918a4c684d',
      mtoShipmentType: SHIPMENT_OPTIONS.NTSR,
      priceCents: 6000001,
      status: 'DENIED',
      rejectionReason: 'Duplicate charge',
    },
  ],
  reviewedAt: '2020-12-01T00:00:00.000Z',
};

const onEditAccountingCodes = () => {};

export const NeedsReview = () => {
  const [modifiedShipments, setModifiedShipments] = React.useState(shipmentsInfo);
  const handleEditAccountingCodes = (id, values) => {
    const updatedShipments = modifiedShipments.map((s) => {
      if (s.mtoShipmentID !== id) {
        return s;
      }

      return {
        ...s,
        tacType: values.tacType,
        sacType: values.sacType,
      };
    });

    setModifiedShipments(updatedShipments);
  };

  return (
    <PaymentRequestCard
      paymentRequest={pendingPaymentRequest}
      shipmentsInfo={modifiedShipments}
      onEditAccountingCodes={handleEditAccountingCodes}
    />
  );
};

export const Reviewed = () => (
  <PaymentRequestCard
    paymentRequest={reviewedPaymentRequest}
    shipmentsInfo={shipmentsInfo}
    onEditAccountingCodes={onEditAccountingCodes}
  />
);

export const Rejected = () => (
  <PaymentRequestCard
    paymentRequest={rejectedPaymentRequest}
    shipmentsInfo={shipmentsInfo}
    onEditAccountingCodes={onEditAccountingCodes}
  />
);
