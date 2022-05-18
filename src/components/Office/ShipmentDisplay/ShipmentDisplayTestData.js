import { shipmentStatuses } from 'constants/shipments';

export const ordersLOA = {
  tac: '1111',
  sac: '2222222222',
  ntsTac: '3333',
  ntsSac: '4444444444',
};

const pickupAddress = {
  streetAddress1: '812 S 129th St',
  city: 'San Antonio',
  state: 'TX',
  postalCode: '78234',
};

const destinationAddress = {
  streetAddress1: '441 SW Rio de la Plata Drive',
  city: 'Tacoma',
  state: 'WA',
  postalCode: '98421',
};

export const usesExternalVendor = true;

export const hhgInfo = {
  heading: 'HHG',
  shipmentId: 'testShipmentId394',
  requestedPickupDate: '26 Mar 2020',
  pickupAddress,
  destinationAddress,
};

export const ntsInfo = {
  heading: 'NTS',
  requestedPickupDate: '26 Mar 2020',
  shipmentId: 'testShipmentId394',
  pickupAddress,
  destinationAddress,
};

export const ntsMissingInfo = {
  heading: 'NTS',
  requestedPickupDate: '26 Mar 2020',
  shipmentId: 'testShipmentId394',
  pickupAddress,
  destinationAddress,
};

export const ntsReleaseInfo = {
  heading: 'NTS-release',
  shipmentId: 'testShipmentId111',
  shipmentStatus: shipmentStatuses.SUBMITTED,
  ntsRecordedWeight: 2000,
  isDiversion: false,
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
  requestedDeliveryDate: '26 Mar 2020',
  destinationAddress,
  secondaryDeliveryAddress: pickupAddress,
  agents: [
    {
      agentType: 'RECEIVING_AGENT',
      firstName: 'Kate',
      lastName: 'Smith',
      phone: '419-555-9999',
      email: 'ksmith@email.com',
    },
  ],
  counselorRemarks:
    'Lorem ipsum dolor sit amet, consectetur adipiscing elit. Aliquam vulputate commodo erat. ' +
    'Morbi porta nibh nibh, ac malesuada tortor egestas.',
  customerRemarks: 'Ut enim ad minima veniam',
  tacType: 'HHG',
  sacType: 'NTS',
};

export const ntsReleaseMissingInfo = {
  heading: 'NTS-release',
  shipmentId: 'testShipmentId222',
  ntsRecordedWeight: 2000,
  requestedDeliveryDate: '26 Mar 2020',
  destinationAddress,
  agents: [
    {
      agentType: 'RECEIVING_AGENT',
      firstName: 'Kate',
      lastName: 'Smith',
      phone: '419-555-9999',
      email: 'ksmith@email.com',
    },
  ],
  sacType: 'NTS',
};

export const postalOnlyInfo = {
  heading: 'HHG',
  requestedPickupDate: '26 Mar 2020',
  shipmentId: 'testShipmentId394',
  pickupAddress,
  destinationAddress: {
    postalCode: '98421',
  },
};

export const diversionInfo = {
  heading: 'HHG',
  shipmentId: 'testShipmentId394',
  isDiversion: true,
  requestedPickupDate: '26 Mar 2020',
  pickupAddress,
  destinationAddress,
  counselorRemarks: 'counselor approved',
};

export const cancelledInfo = {
  heading: 'HHG',
  shipmentId: 'testShipmentId394',
  isDiversion: false,
  shipmentStatus: shipmentStatuses.CANCELED,
  requestedPickupDate: '26 Mar 2020',
  pickupAddress,
  destinationAddress,
  counselorRemarks: 'counselor approved',
};

export const ppmInfo = {
  heading: 'PPM',
  ppmShipment: {
    actualMoveDate: null,
    advance: 598700,
    advanceRequested: true,
    approvedAt: null,
    createdAt: '2022-04-29T21:48:21.581Z',
    deletedAt: null,
    destinationPostalCode: '30813',
    eTag: 'MjAyMi0wNC0yOVQyMTo0ODoyMS41ODE0MzFa',
    estimatedIncentive: 1000000,
    estimatedWeight: 4000,
    expectedDepartureDate: '2020-03-15',
    hasProGear: true,
    id: 'b6ec215c-2cef-45fe-8d4a-35f445cd4768',
    netWeight: null,
    pickupPostalCode: '90210',
    proGearWeight: 1987,
    reviewedAt: null,
    secondaryDestinationPostalCode: '30814',
    secondaryPickupPostalCode: '90211',
    shipmentId: 'b5c2d9a1-d1e6-485d-9678-8b62deb0d801',
    spouseProGearWeight: 498,
    status: 'SUBMITTED',
    submittedAt: '2022-04-29T21:48:21.573Z',
    updatedAt: '2022-04-29T21:48:21.581Z',
  },
};
