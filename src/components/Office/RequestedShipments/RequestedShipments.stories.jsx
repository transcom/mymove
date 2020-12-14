import React from 'react';

import RequestedShipments from './RequestedShipments';

import { SHIPMENT_OPTIONS } from 'shared/constants';

const shipments = [
  {
    approvedDate: '0001-01-01',
    createdAt: '2020-06-10T15:58:02.404029Z',
    customerRemarks: 'please treat gently',
    destinationAddress: {
      city: 'Fairfield',
      country: 'US',
      eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi4zODk0MTJa',
      id: '672ff379-f6e3-48b4-a87d-796713f8f997',
      postal_code: '94535',
      state: 'CA',
      street_address_1: '987 Any Avenue',
      street_address_2: 'P.O. Box 9876',
      street_address_3: 'c/o Some Person',
    },
    eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi40MDQwMzFa',
    id: 'ce01a5b8-9b44-4511-8a8d-edb60f2a4aea',
    moveTaskOrderID: '9c7b255c-2981-4bf8-839f-61c7458e2b4d',
    pickupAddress: {
      city: 'Beverly Hills',
      country: 'US',
      eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi4zODQ3Njla',
      id: '1686751b-ab36-43cf-b3c9-c0f467d13c19',
      postal_code: '90210',
      state: 'CA',
      street_address_1: '123 Any Street',
      street_address_2: 'P.O. Box 12345',
      street_address_3: 'c/o Some Person',
    },
    rejectionReason: 'shipment not good enough',
    requestedPickupDate: '2018-03-15',
    scheduledPickupDate: '2018-03-16',
    secondaryDeliveryAddress: {
      city: 'Beverly Hills',
      country: 'US',
      eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi4zOTkzMlo=',
      id: '15e8f6cc-e1d7-44b2-b1e0-fcb3d6442831',
      postal_code: '90210',
      state: 'CA',
      street_address_1: '123 Any Street',
      street_address_2: 'P.O. Box 12345',
      street_address_3: 'c/o Some Person',
    },
    secondaryPickupAddress: {
      city: 'Beverly Hills',
      country: 'US',
      eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi4zOTM4OTZa',
      id: '9b79e0c3-8ed5-4fb8-aa36-95845707d8ee',
      postal_code: '90210',
      state: 'CA',
      street_address_1: '123 Any Street',
      street_address_2: 'P.O. Box 12345',
      street_address_3: 'c/o Some Person',
    },
    shipmentType: SHIPMENT_OPTIONS.HHG,
    status: 'SUBMITTED',
    updatedAt: '2020-06-10T15:58:02.404031Z',
  },
  {
    approvedDate: '0001-01-01',
    createdAt: '2020-06-10T15:58:02.404029Z',
    customerRemarks: 'please treat gently',
    destinationAddress: {
      city: 'Fairfield',
      country: 'US',
      eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi4zODk0MTJa',
      id: '672ff379-f6e3-48b4-a87d-796713f8f997',
      postal_code: '94535',
      state: 'CA',
      street_address_1: '987 Any Avenue',
      street_address_2: 'P.O. Box 9876',
      street_address_3: 'c/o Some Person',
    },
    eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi40MDQwMzFa',
    id: 'cd01a5b8-9b44-4511-8a8d-edb60f2a4aea',
    moveTaskOrderID: '9c7b255c-2981-4bf8-839f-61c7458e2b4d',
    pickupAddress: {
      city: 'Beverly Hills',
      country: 'US',
      eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi4zODQ3Njla',
      id: '1686751b-ab36-43cf-b3c9-c0f467d13c19',
      postal_code: '90210',
      state: 'CA',
      street_address_1: '123 Any Street',
      street_address_2: 'P.O. Box 12345',
      street_address_3: 'c/o Some Person',
    },
    rejectionReason: 'shipment not good enough',
    requestedPickupDate: '2018-03-15',
    scheduledPickupDate: '2018-03-16',
    secondaryDeliveryAddress: {
      city: 'Beverly Hills',
      country: 'US',
      eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi4zOTkzMlo=',
      id: '15e8f6cc-e1d7-44b2-b1e0-fcb3d6442831',
      postal_code: '90210',
      state: 'CA',
      street_address_1: '123 Any Street',
      street_address_2: 'P.O. Box 12345',
      street_address_3: 'c/o Some Person',
    },
    secondaryPickupAddress: {
      city: 'Beverly Hills',
      country: 'US',
      eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi4zOTM4OTZa',
      id: '9b79e0c3-8ed5-4fb8-aa36-95845707d8ee',
      postal_code: '90210',
      state: 'CA',
      street_address_1: '123 Any Street',
      street_address_2: 'P.O. Box 12345',
      street_address_3: 'c/o Some Person',
    },
    shipmentType: SHIPMENT_OPTIONS.HHG,
    status: 'SUBMITTED',
    updatedAt: '2020-06-10T15:58:02.404031Z',
  },
  {
    approvedDate: '0001-01-01',
    createdAt: '2020-06-10T15:58:02.431993Z',
    customerRemarks: 'please treat gently',
    destinationAddress: {
      city: 'Fairfield',
      country: 'US',
      eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi40MTcyMjZa',
      id: '00a5dfeb-c6a0-4ed8-965c-89943163fee4',
      postal_code: '94535',
      state: 'CA',
      street_address_1: '987 Any Avenue',
      street_address_2: 'P.O. Box 9876',
      street_address_3: 'c/o Some Person',
    },
    eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi40MzE5OTVa',
    id: 'c2f68d97-b960-4c86-a418-c70a0aeba04e',
    moveTaskOrderID: '9c7b255c-2981-4bf8-839f-61c7458e2b4d',
    pickupAddress: {
      city: 'Beverly Hills',
      country: 'US',
      eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi40MTMyNDha',
      id: '14b1d10d-b34b-4ec5-80e6-69d885206a2a',
      postal_code: '90210',
      state: 'CA',
      street_address_1: '123 Any Street',
      street_address_2: 'P.O. Box 12345',
      street_address_3: 'c/o Some Person',
    },
    rejectionReason: 'shipment not good enough',
    requestedPickupDate: '2018-03-15',
    scheduledPickupDate: '2018-03-16',
    secondaryDeliveryAddress: {
      city: 'Beverly Hills',
      country: 'US',
      eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi40MjYxODVa',
      id: '1a4f6fec-42b9-4dd2-b205-c6770ac7ea27',
      postal_code: '90210',
      state: 'CA',
      street_address_1: '123 Any Street',
      street_address_2: 'P.O. Box 12345',
      street_address_3: 'c/o Some Person',
    },
    secondaryPickupAddress: {
      city: 'Beverly Hills',
      country: 'US',
      eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi40MjIwNzVa',
      id: 'e188f33f-f84d-4f86-954a-938b52e38741',
      postal_code: '90210',
      state: 'CA',
      street_address_1: '123 Any Street',
      street_address_2: 'P.O. Box 12345',
      street_address_3: 'c/o Some Person',
    },
    shipmentType: SHIPMENT_OPTIONS.NTS,
    status: 'SUBMITTED',
    updatedAt: '2020-06-10T15:58:02.431995Z',
  },
];

const ordersInfo = {
  newDutyStation: {
    address: {
      city: 'Augusta',
      country: 'United States',
      eTag: 'MjAyMC0wOC0wNlQxNDo1Mjo0MS45NDQ0ODla',
      id: '5ac95be8-0230-47ea-90b4-b0f6f60de364',
      postal_code: '30813',
      state: 'GA',
      street_address_1: 'Fort Gordon',
    },
    address_id: '5ac95be8-0230-47ea-90b4-b0f6f60de364',
    eTag: 'MjAyMC0wOC0wNlQxNDo1Mjo0MS45NDQ0ODla',
    id: '2d5ada83-e09a-47f8-8de6-83ec51694a86',
    name: 'Fort Gordon',
  },
  currentDutyStation: {
    address: {
      city: 'Des Moines',
      country: 'US',
      eTag: 'MjAyMC0wOC0wNlQxNDo1MzozMC42NjEwODFa',
      id: '37880d6d-2c78-47f1-a71b-53c0ea1a0107',
      postal_code: '50309',
      state: 'IA',
      street_address_1: '987 Other Avenue',
      street_address_2: 'P.O. Box 1234',
      street_address_3: 'c/o Another Person',
    },
    address_id: '37880d6d-2c78-47f1-a71b-53c0ea1a0107',
    eTag: 'MjAyMC0wOC0wNlQxNDo1MzozMC42Njg5MDFa',
    id: '07282a8f-a496-4648-ae24-119775eef57d',
    name: 'vC6w22RPYC',
  },
  issuedDate: '2018-03-15',
  reportByDate: '2018-08-01',
  departmentIndicator: 'COAST_GUARD',
  ordersNumber: 'ORDER3',
  ordersType: 'PERMANENT_CHANGE_OF_STATION',
  ordersTypeDetail: 'TBD',
  tacMDC: '',
  sacSDN: '',
};

const allowancesInfo = {
  branch: 'Navy',
  rank: 'E-6',
  weightAllowance: '11,000 lbs',
  authorizedWeight: '11,000 lbs',
  progear: '2,000 lbs',
  spouseProgear: '500 lbs',
  storageInTransit: '90 days',
  dependents: 'Authorized',
};

const customerInfo = {
  name: 'Smith, Kerry',
  dodId: '9999999999',
  phone: '+1 999-999-9999',
  email: 'ksmith@email.com',
  currentAddress: {
    street_address_1: '812 S 129th St',
    city: 'San Antonio',
    state: 'TX',
    postal_code: '78234',
  },
  backupContactName: 'Quinn Ocampo',
  backupContactPhone: '+1 999-999-9999',
  backupContactEmail: 'quinnocampo@myemail.com',
};

const agents = [
  {
    type: 'RELEASING_AGENT',
    firstName: 'Dorothy',
    lastName: 'Lagomarsino',
    email: 'dorothyl@email.com',
    phone: '+1 999-999-9999',
    shipmentId: 'ce01a5b8-9b44-4511-8a8d-edb60f2a4aea',
  },
  {
    type: 'RECEIVING_AGENT',
    firstName: 'Dorothy Lagomarsino',
    lastName: 'Lagomarsino',
    email: 'dorothyl@email.com',
    phone: '+1 999-999-9999',
    shipmentId: 'ce01a5b8-9b44-4511-8a8d-edb60f2a4aea',
  },
];

const serviceItems = [
  {
    reServiceName: 'Move management',
    approvedAt: '2020-01-01',
    id: '76055c99-0990-410c-a7c9-69373b0b53eb',
    status: 'APPROVED',
    reServiceCode: 'MS',
  },
  {
    reServiceName: 'Counseling fee',
    id: '76055c99-0990-410c-a7c9-69373b0b5322',
    status: 'APPROVED',
    reServiceCode: 'CS',
    approvedAt: '2020-01-01',
  },
];

export default {
  title: 'Office Components/RequestedShipments',
};

export const withOneShipment = () => (
  <div style={{ padding: '20px' }}>
    <RequestedShipments
      mtoShipments={[shipments[0]]}
      ordersInfo={ordersInfo}
      allowancesInfo={allowancesInfo}
      customerInfo={customerInfo}
      mtoAgents={agents}
      shipmentsStatus="SUBMITTED"
    />
  </div>
);

export const withMultipleShipments = () => (
  <div style={{ position: 'relative', padding: '20px' }}>
    <RequestedShipments
      mtoShipments={shipments}
      ordersInfo={ordersInfo}
      allowancesInfo={allowancesInfo}
      customerInfo={customerInfo}
      mtoAgents={agents}
      shipmentsStatus="SUBMITTED"
    />
  </div>
);

export const withOneApprovedShipment = () => (
  <div style={{ position: 'relative', padding: '20px' }}>
    <RequestedShipments
      mtoShipments={[shipments[0]]}
      ordersInfo={ordersInfo}
      allowancesInfo={allowancesInfo}
      customerInfo={customerInfo}
      mtoAgents={agents}
      shipmentsStatus="APPROVED"
      mtoServiceItems={serviceItems}
    />
  </div>
);

export const withMultipleApprovedShipments = () => (
  <div style={{ position: 'relative', padding: '20px' }}>
    <RequestedShipments
      mtoShipments={shipments}
      ordersInfo={ordersInfo}
      allowancesInfo={allowancesInfo}
      customerInfo={customerInfo}
      mtoAgents={agents}
      shipmentsStatus="APPROVED"
      mtoServiceItems={serviceItems}
    />
  </div>
);
