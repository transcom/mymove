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
    serviceRequestDocuments: [
      {
        uploads: [
          {
            filename: '/mto-service-item/ae1c6472-5e03-4f9c-bef5-55605dbeb31e/20230630161854-receipt.png',
            url: '/storage//mto-service-item/ae1c6472-5e03-4f9c-bef5-55605dbeb31e/20230630161854-receipt.png?contentType=image%2Fpng',
          },
        ],
      },
    ],
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
      customerContacts: [
        { timeMilitary: '', firstAvailableDeliveryDate: '2020-11-22' },
        { timeMilitary: '', firstAvailableDeliveryDate: '2020-11-23' },
      ],
    },
    serviceRequestDocuments: [
      {
        uploads: [
          {
            filename: '/mto-service-item/ae1c6472-5e03-4f9c-bef5-55605dbeb31e/20230630161854-storage-document.pdf',
            url: '/storage//mto-service-item/ae1c6472-5e03-4f9c-bef5-55605dbeb31e/20230630161854-storage-document.pdf?contentType=application%2Fpdf',
          },
        ],
      },
    ],
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

const serviceItemAddressUpdateAlert = {
  makeVisible: false,
  alertMessage: '',
  alertType: '',
};

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
    <RequestedServiceItemsTable
      serviceItems={serviceItems}
      statusForTableType={SERVICE_ITEM_STATUS.SUBMITTED}
      serviceItemAddressUpdateAlert={serviceItemAddressUpdateAlert}
    />
  </MockProviders>
);
export const ApprovedServiceItems = () => (
  <MockProviders permissions={[permissionTypes.updateMTOServiceItem]}>
    <RequestedServiceItemsTable
      serviceItems={approvedServiceItems}
      statusForTableType={SERVICE_ITEM_STATUS.APPROVED}
      serviceItemAddressUpdateAlert={serviceItemAddressUpdateAlert}
    />
  </MockProviders>
);
export const RejectedServiceItems = () => (
  <MockProviders permissions={[permissionTypes.updateMTOServiceItem]}>
    <RequestedServiceItemsTable
      serviceItems={rejectedServiceItems}
      statusForTableType={SERVICE_ITEM_STATUS.REJECTED}
      serviceItemAddressUpdateAlert={serviceItemAddressUpdateAlert}
    />
  </MockProviders>
);
