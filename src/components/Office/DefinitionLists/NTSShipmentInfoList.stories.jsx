import React from 'react';

import NTSShipmentInfoList from './NTSShipmentInfoList';

export default {
  title: 'Office Components/Shipment Info List',
  component: NTSShipmentInfoList,
};

const info = {
  storageFacility: {
    address: {
      city: 'Anytown',
      country: 'USA',
      postalCode: '90210',
      state: 'OK',
      streetAddress1: '555 Main Ave',
      streetAddress2: 'Apartment 900',
    },
    facilityName: 'my storage',
    lotNumber: '2222',
  },
  serviceOrderNumber: '12341234',
  requestedPickupDate: '26 Mar 2020',
  pickupAddress: {
    streetAddress1: '441 SW Rio de la Plata Drive',
    city: 'Tacoma',
    state: 'WA',
    postalCode: '98421',
  },
  secondaryPickupAddress: {
    streetAddress1: '812 S 129th St',
    city: 'San Antonio',
    state: 'TX',
    postalCode: '78234',
  },
  agents: [
    {
      agentType: 'RELEASING_AGENT',
      firstName: 'Jason',
      lastName: 'Ash',
      phone: '419-555-5555',
      email: 'jash@email.com',
    },
  ],
  counselorRemarks:
    'Lorem ipsum dolor sit amet, consectetur adipiscing elit. Aliquam vulputate commodo erat. ' +
    'Morbi porta nibh nibh, ac malesuada tortor egestas.',
  customerRemarks: 'Ut enim ad minima veniam',
  tacType: 'HHG',
  sacType: 'NTS',
  tac: '123',
  sac: '456',
};

export const NTSBasic = () => (
  <NTSShipmentInfoList
    shipment={{
      counselorRemarks: info.counselorRemarks,
      requestedPickupDate: info.requestedPickupDate,
      storageFacility: info.storageFacility,
      pickupAddress: info.pickupAddress,
      tacType: info.tacType,
      sacType: info.sacType,
      serviceOrderNumber: info.serviceOrderNumber,
    }}
  />
);

export const NTSBasicEvaluationReport = () => (
  <NTSShipmentInfoList
    shipment={{
      counselorRemarks: info.counselorRemarks,
      requestedPickupDate: info.requestedPickupDate,
      storageFacility: info.storageFacility,
      pickupAddress: info.pickupAddress,
      tacType: info.tacType,
      sacType: info.sacType,
      serviceOrderNumber: info.serviceOrderNumber,
    }}
    isForEvaluationReport
  />
);

export const NTSMissingInfo = () => (
  <NTSShipmentInfoList
    isExpanded
    shipment={{
      requestedPickupDate: info.requestedPickupDate,
      pickupAddress: info.pickupAddress,
      sacType: info.sacType,
    }}
    errorIfMissing={['storageFacility', 'serviceOrderNumber', 'tacType']}
  />
);

export const NTSWarning = () => (
  <NTSShipmentInfoList
    isExpanded
    shipment={{
      requestedPickupDate: info.requestedPickupDate,
      pickupAddress: info.pickupAddress,
      sacType: info.sacType,
    }}
    warnIfMissing={['storageFacility', 'serviceOrderNumber', 'tacType']}
  />
);

export const NTSWithAllInfo = () => (
  <NTSShipmentInfoList
    isExpanded
    shipment={{
      requestedPickupDate: info.requestedPickupDate,
      storageFacility: info.storageFacility,
      tacType: info.tacType,
      sacType: info.sacType,
      tac: info.tac,
      sac: info.sac,
      serviceOrderNumber: info.serviceOrderNumber,
      pickupAddress: info.pickupAddress,
      secondaryPickupAddress: info.secondaryPickupAddress,
      agents: info.agents,
      counselorRemarks: info.counselorRemarks,
      customerRemarks: info.customerRemarks,
    }}
  />
);
