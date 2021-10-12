import React from 'react';

import { SITStatusOrigin } from '../ShipmentSITDisplay/ShipmentSITDisplayTestParams';

import ShipmentDetails from './ShipmentDetails';

export default {
  title: 'Office Components/Shipment Details',
};

const shipment = {
  requestedPickupDate: '2021-06-01',
  scheduledPickupDate: '2021-06-03',
  customerRemarks: 'Please treat gently.',
  counselorRemarks: 'This shipment is to be treated with care.',
  pickupAddress: {
    street_address_1: '812 S 129th St',
    city: 'San Antonio',
    state: 'TX',
    postal_code: '78234',
  },
  secondaryPickupAddress: {
    street_address_1: '444 S 131st St',
    city: 'San Antonio',
    state: 'TX',
    postal_code: '78234',
  },
  destinationAddress: {
    street_address_1: '7 Q St',
    city: 'Austin',
    state: 'TX',
    postal_code: '78722',
  },
  secondaryDeliveryAddress: {
    street_address_1: '17 8th St',
    city: 'Austin',
    state: 'TX',
    postal_code: '78751',
  },
  primeEstimatedWeight: 4000,
  primeActualWeight: 3800,
  agents: [
    {
      agentType: 'RELEASING_AGENT',
      firstName: 'Quinn',
      lastName: 'Ocampo',
      phone: '999-999-9999',
      email: 'quinnocampo@myemail.com',
    },
    {
      agentType: 'RECEIVING_AGENT',
      firstName: 'Kate',
      lastName: 'Smith',
      phone: '419-555-9999',
      email: 'ksmith@email.com',
    },
  ],
  reweigh: {
    id: '00000000-0000-0000-0000-000000000000',
  },
  sitExtensions: [
    {
      contractorRemarks: 'The customer requested an extension.',
      createdAt: '2021-09-13T15:41:59.373Z',
      decisionDate: '0001-01-01T00:00:00.000Z',
      eTag: 'MjAyMS0wOS0xM1QxNTo0MTo1OS4zNzM2NTRa',
      id: '7af5d51a-789c-4f5e-83dd-d905daed0785',
      mtoShipmentID: '8afd043a-8304-4e36-a695-7728e415990d',
      officeRemarks: 'The service member is unable to move into their new home at the expected time.',
      requestReason: 'SERIOUS_ILLNESS_MEMBER',
      requestedDays: 30,
      approvedDays: 30,
      status: 'APPROVED',
      updatedAt: '2021-09-13T15:41:59.373Z',
    },
  ],
  sitStatus: SITStatusOrigin,
  sitDaysAllowance: 270,
};

const order = {
  originDutyStation: {
    address: {
      street_address_1: '444 S 131st St',
      city: 'San Antonio',
      state: 'TX',
      postal_code: '78234',
    },
  },
  destinationDutyStation: {
    address: {
      street_address_1: '17 8th St',
      city: 'Austin',
      state: 'TX',
      postal_code: '78751',
    },
  },
};

export const Default = () => <ShipmentDetails shipment={shipment} order={order} />;
