import React from 'react';
import { MemoryRouter, Route } from 'react-router';

import { SITStatusOrigin } from '../ShipmentSITDisplay/ShipmentSITDisplayTestParams';

import ShipmentDetails from './ShipmentDetails';

import { LOA_TYPE } from 'shared/constants';

export default {
  title: 'Office Components/Shipment Details',
  decorators: [
    (Story) => (
      <MemoryRouter initialEntries={['/moves/HGNTSR/mto']}>
        <Route path="/moves/:moveCode/mto">
          <Story />
        </Route>
      </MemoryRouter>
    ),
  ],
};

const shipment = {
  requestedPickupDate: '2021-06-01',
  scheduledPickupDate: '2021-06-03',
  customerRemarks: 'Please treat gently.',
  counselorRemarks: 'This shipment is to be treated with care.',
  pickupAddress: {
    streetAddress1: '812 S 129th St',
    city: 'San Antonio',
    state: 'TX',
    postalCode: '78234',
  },
  secondaryPickupAddress: {
    streetAddress1: '444 S 131st St',
    city: 'San Antonio',
    state: 'TX',
    postalCode: '78234',
  },
  destinationAddress: {
    streetAddress1: '7 Q St',
    city: 'Austin',
    state: 'TX',
    postalCode: '78722',
  },
  secondaryDeliveryAddress: {
    streetAddress1: '17 8th St',
    city: 'Austin',
    state: 'TX',
    postalCode: '78751',
  },
  primeEstimatedWeight: 4000,
  primeActualWeight: 3800,
  mtoAgents: [
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
  storageFacility: {
    facilityName: 'Most Excellent Storage',
    address: {
      streetAddress1: '3373 NW Martin Luther King Jr Blvd',
      city: 'San Antonio',
      state: 'TX',
      postalCode: '78212',
    },
    phone: '555-555-5555',
    lotNumber: '64321',
  },
  serviceOrderNumber: '1234',
  tacType: LOA_TYPE.HHG,
  sacType: LOA_TYPE.NTS,
};

const order = {
  originDutyLocation: {
    address: {
      streetAddress1: '444 S 131st St',
      city: 'San Antonio',
      state: 'TX',
      postalCode: '78234',
    },
  },
  destinationDutyLocation: {
    address: {
      streetAddress1: '17 8th St',
      city: 'Austin',
      state: 'TX',
      postalCode: '78751',
    },
  },
  tac: '1234',
  sac: '567',
  ntsTac: '8912',
  ntsSac: '345',
};

export const Default = () => {
  const [modifiedShipment, setModifiedShipment] = React.useState(shipment);
  const handleEditSon = (values) => {
    setModifiedShipment({
      ...shipment,
      serviceOrderNumber: values.serviceOrderNumber,
    });
  };

  return (
    <div className="officeApp">
      <ShipmentDetails shipment={modifiedShipment} order={order} handleEditServiceOrderNumber={handleEditSon} />
    </div>
  );
};
