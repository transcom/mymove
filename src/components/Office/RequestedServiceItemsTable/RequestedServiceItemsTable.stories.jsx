import React from 'react';

import { SERVICE_ITEM_STATUS } from '../../../shared/constants';

import RequestedServiceItemsTable from './RequestedServiceItemsTable';

import { MockProviders } from 'testUtils';
import { permissionTypes } from 'constants/permissions';

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
      SITPostalCode: '22030',
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
      estimatedWeight: 4999,
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
      estimatedWeight: 4999,
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
    },
  },
];

const approvedServiceItems = serviceItems.map((serviceItem) => {
  return {
    ...serviceItem,
    status: 'APPROVED',
  };
});
const rejectedServiceItems = serviceItems.map((serviceItem) => {
  return {
    ...serviceItem,
    status: 'REJECTED',
    details: { ...serviceItem.details, rejectionReason: 'Here is a reason for rejection' },
  };
});

export const Default = () => (
  <MockProviders permissions={[permissionTypes.updateMTOServiceItem]}>
    <RequestedServiceItemsTable serviceItems={serviceItems} statusForTableType={SERVICE_ITEM_STATUS.SUBMITTED} />
  </MockProviders>
);
export const ApprovedServiceItems = () => (
  <MockProviders permissions={[permissionTypes.updateMTOServiceItem]}>
    <RequestedServiceItemsTable serviceItems={approvedServiceItems} statusForTableType={SERVICE_ITEM_STATUS.APPROVED} />
  </MockProviders>
);
export const RejectedServiceItems = () => (
  <MockProviders permissions={[permissionTypes.updateMTOServiceItem]}>
    <RequestedServiceItemsTable serviceItems={rejectedServiceItems} statusForTableType={SERVICE_ITEM_STATUS.REJECTED} />
  </MockProviders>
);
