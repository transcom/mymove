import React from 'react';
import { isHappoRun } from 'happo-plugin-storybook/register';

import { SERVICE_ITEM_STATUS } from '../../../shared/constants';

import RequestedServiceItemsTable from './RequestedServiceItemsTable';

export default {
  title: 'Office Components/RequestedServiceItemsTable',
  component: RequestedServiceItemsTable,
};

const serviceItems = [
  {
    id: 'abc-123',
    createdAt: '2020-11-20T00:00:00',
    approvedAt: '2020-11-20T00:00:00',
    rejectedAt: '2020-11-20T00:00:00',
    serviceItem: 'Domestic origin 1st day SIT',
    code: 'DOFSIT',
    status: 'SUBMITTED',
    details: {
      pickupPostalCode: '60612',
      reason: "here's the reason",
    },
  },
  {
    id: 'abc-1234',
    createdAt: '2020-11-20T00:00:00',
    approvedAt: '2020-11-20T00:00:00',
    rejectedAt: '2020-11-20T00:00:00',
    serviceItem: 'Domestic destination 1st day SIT',
    code: 'DDFSIT',
    status: 'SUBMITTED',
    details: {
      firstCustomerContact: { timeMilitary: '', firstAvailableDeliveryDate: '2020-11-22' },
      secondCustomerContact: { timeMilitary: '', firstAvailableDeliveryDate: '2020-11-23' },
    },
  },
  {
    id: 'cba-123',
    createdAt: '2020-11-20T00:00:00',
    approvedAt: '2020-11-20T00:00:00',
    rejectedAt: '2020-11-20T00:00:00',
    serviceItem: 'Domestic origin shuttle service',
    code: 'DOSHUT',
    status: 'SUBMITTED',
    details: {
      reason: "Here's the reason",
    },
  },
  {
    id: 'cba-1234',
    createdAt: '2020-11-20T00:00:00',
    approvedAt: '2020-11-20T00:00:00',
    rejectedAt: '2020-11-20T00:00:00',
    serviceItem: 'Domestic destination shuttle service',
    code: 'DDSHUT',
    status: 'SUBMITTED',
    details: {
      reason: "Here's the reason",
    },
  },
  {
    id: 'abc12345',
    createdAt: '2020-11-20T00:00:00',
    approvedAt: '2020-11-20T00:00:00',
    rejectedAt: '2020-11-20T00:00:00',
    serviceItem: 'Domestic crating',
    code: 'DCRT',
    status: 'SUBMITTED',
    details: {
      description: "Here's the description",
      itemDimensions: { length: 8400, width: 2600, height: 4200 },
      crateDimensions: { length: 110000, width: 36000, height: 54000 },
      imgURL: isHappoRun() ? null : 'https://live.staticflickr.com/4735/24289917967_27840ed1af_b.jpg',
    },
  },
];

const approvedServiceItems = serviceItems.map((serviceItem) => {
  return { ...serviceItem, status: 'APPROVED' };
});
const rejectedServiceItems = serviceItems.map((serviceItem) => {
  return { ...serviceItem, status: 'REJECTED' };
});

export const Default = () => (
  <RequestedServiceItemsTable serviceItems={serviceItems} statusForTableType={SERVICE_ITEM_STATUS.SUBMITTED} />
);
export const ApprovedServiceItems = () => (
  <RequestedServiceItemsTable serviceItems={approvedServiceItems} statusForTableType={SERVICE_ITEM_STATUS.APPROVED} />
);
export const RejectedServiceItems = () => (
  <RequestedServiceItemsTable serviceItems={rejectedServiceItems} statusForTableType={SERVICE_ITEM_STATUS.REJECTED} />
);
