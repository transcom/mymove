import React from 'react';

import { SERVICE_ITEM_STATUS } from '../../../shared/constants';

import RequestedServiceItemsTable from './RequestedServiceItemsTable';

import { MockProviders } from 'testUtils';
import { permissionTypes } from 'constants/permissions';

export default {
  title: 'Office Components/RequestedServiceItemsTable',
  component: RequestedServiceItemsTable,
};

const domesticOriginSitDetails = {
  pickupPostalCode: '60612',
  SITPostalCode: '22030',
  reason: 'Housing is not ready',
};

const domesticOriginSitDocuments = [
  {
    uploads: [
      {
        filename: '/mto-service-item/ae1c6472-5e03-4f9c-bef5-55605dbeb31e/20230630161854-receipt.png',
        url: '/storage//mto-service-item/ae1c6472-5e03-4f9c-bef5-55605dbeb31e/20230630161854-receipt.png?contentType=image%2Fpng',
      },
    ],
  },
];

const DOFSIT = {
  id: 'dosit-123',
  createdAt: '2020-11-20T00:00:00',
  approvedAt: '2020-12-20T00:00:00',
  rejectedAt: '2020-13-20T00:00:00',
  serviceItem: 'Domestic origin 1st day SIT',
  code: 'DOFSIT',
  status: 'SUBMITTED',
  details: domesticOriginSitDetails,
  serviceRequestDocuments: domesticOriginSitDocuments,
};

const DOASIT = {
  id: 'dosit-234',
  createdAt: '2020-11-20T00:00:00',
  approvedAt: '2020-12-20T00:00:00',
  rejectedAt: '2020-13-20T00:00:00',
  serviceItem: 'Domestic origin Additional day SIT',
  code: 'DOASIT',
  status: 'SUBMITTED',
  details: domesticOriginSitDetails,
  serviceRequestDocuments: domesticOriginSitDocuments,
};

const DOPSIT = {
  id: 'dosit-345',
  createdAt: '2020-11-20T00:00:00',
  approvedAt: '2020-12-20T00:00:00',
  rejectedAt: '2020-13-20T00:00:00',
  serviceItem: 'Domestic origin SIT pickup',
  code: 'DOPSIT',
  status: 'SUBMITTED',
  details: domesticOriginSitDetails,
  serviceRequestDocuments: domesticOriginSitDocuments,
};

const DOSFSC = {
  id: 'abc-456',
  createdAt: '2020-11-20T00:00:00',
  approvedAt: '2020-12-20T00:00:00',
  rejectedAt: '2020-13-20T00:00:00',
  serviceItem: 'Domestic origin SIT fuel surcharge',
  code: 'DOSFSC',
  status: 'SUBMITTED',
  details: domesticOriginSitDetails,
  serviceRequestDocuments: domesticOriginSitDocuments,
};

const domesticDestinationSitDetails = {
  reason: "Customer's housing at base is not ready",
  sitEntryDate: '2022-11-01',
  sitDepartureDate: '2022-12-10',
  customerContacts: [
    { dateOfContact: '2022-11-10', timeMilitary: '1400Z', firstAvailableDeliveryDate: '2022-11-14' },
    { dateOfContact: '2022-11-01', timeMilitary: '1400Z', firstAvailableDeliveryDate: '2022-11-05' },
  ],
};

const domesticDestinationSitDocuments = [
  {
    uploads: [
      {
        filename: '/mto-service-item/ae1c6472-5e03-4f9c-bef5-55605dbeb31e/ImageOfItems.png',
        url: '/storage//mto-service-item/ae1c6472-5e03-4f9c-bef5-55605dbeb31e/ImageOfItems.png?contentType=image%2Fpng',
      },
    ],
  },
];

const DDFSIT = {
  id: 'ddsit-123',
  createdAt: '2020-11-20T00:00:00',
  approvedAt: '2020-12-20T00:00:00',
  rejectedAt: '2020-13-20T00:00:00',
  serviceItem: 'Domestic destination 1st day SIT',
  code: 'DDFSIT',
  status: 'SUBMITTED',
  details: domesticDestinationSitDetails,
  serviceRequestDocuments: domesticDestinationSitDocuments,
};

const DDASIT = {
  id: 'ddsit-234',
  createdAt: '2020-11-20T00:00:00',
  approvedAt: '2020-12-20T00:00:00',
  rejectedAt: '2020-13-20T00:00:00',
  serviceItem: 'Domestic destination Additional day SIT',
  code: 'DDASIT',
  status: 'SUBMITTED',
  details: domesticDestinationSitDetails,
  serviceRequestDocuments: domesticDestinationSitDocuments,
};

const DDDSIT = {
  id: 'ddsit-345',
  createdAt: '2020-11-20T00:00:00',
  approvedAt: '2020-12-20T00:00:00',
  rejectedAt: '2020-13-20T00:00:00',
  serviceItem: 'Domestic destination SIT delivery',
  code: 'DDDSIT',
  status: 'SUBMITTED',
  details: domesticDestinationSitDetails,
  serviceRequestDocuments: domesticDestinationSitDocuments,
};

const DDSFSC = {
  id: 'ddsit-456',
  createdAt: '2020-11-20T00:00:00',
  approvedAt: '2020-12-20T00:00:00',
  rejectedAt: '2020-13-20T00:00:00',
  serviceItem: 'Domestic destination SIT fuel surcharge',
  code: 'DDSFSC',
  status: 'SUBMITTED',
  details: domesticDestinationSitDetails,
  serviceRequestDocuments: domesticDestinationSitDocuments,
};

const DOSHUT = {
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
};

const DDSHUT = {
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
};

const DCRT = {
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
};

const serviceItems = [
  // Domestic Origin SIT Service Items
  DOFSIT,
  DOASIT,
  DOPSIT,
  DOSFSC,

  // Domestic Destination SIT Service Items
  DDFSIT,
  DDASIT,
  DDDSIT,
  DDSFSC,

  // Other
  DOSHUT,
  DDSHUT,
  DCRT,
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
      handleShowEditSitAddressModal={() => {}}
      handleShowRejectionDialog={() => {}}
      handleUpdateMTOServiceItemStatus={() => {}}
    />
  </MockProviders>
);
export const ApprovedServiceItems = () => (
  <MockProviders permissions={[permissionTypes.updateMTOServiceItem]}>
    <RequestedServiceItemsTable
      serviceItems={approvedServiceItems}
      statusForTableType={SERVICE_ITEM_STATUS.APPROVED}
      serviceItemAddressUpdateAlert={serviceItemAddressUpdateAlert}
      handleShowEditSitAddressModal={() => {}}
      handleShowRejectionDialog={() => {}}
      handleUpdateMTOServiceItemStatus={() => {}}
    />
  </MockProviders>
);
export const RejectedServiceItems = () => (
  <MockProviders permissions={[permissionTypes.updateMTOServiceItem]}>
    <RequestedServiceItemsTable
      serviceItems={rejectedServiceItems}
      statusForTableType={SERVICE_ITEM_STATUS.REJECTED}
      serviceItemAddressUpdateAlert={serviceItemAddressUpdateAlert}
      handleShowEditSitAddressModal={() => {}}
      handleShowRejectionDialog={() => {}}
      handleUpdateMTOServiceItemStatus={() => {}}
    />
  </MockProviders>
);
