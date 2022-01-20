import React from 'react';

import {
  ordersLOA,
  hhgInfo,
  ntsInfo,
  ntsMissingInfo,
  ntsReleaseInfo,
  ntsReleaseMissingInfo,
  postalOnlyInfo,
  diversionInfo,
  cancelledInfo,
  usesExternalVendor,
} from './ShipmentDisplayTestData';

import ShipmentDisplay from 'components/Office/ShipmentDisplay/ShipmentDisplay';
import { SHIPMENT_OPTIONS } from 'shared/constants';
import { MockProviders } from 'testUtils';

export default {
  title: 'Office Components/Shipment Display',
  component: ShipmentDisplay,
  decorators: [
    (Story) => (
      <MockProviders>
        <Story />
      </MockProviders>
    ),
  ],
};

const warnIfMissing = ['ntsRecordedWeight', 'serviceOrderNumber', 'counselorRemarks', 'tacType', 'sacType'];
const errorIfMissing = ['storageFacility'];
const errorIfMissingTACType = ['tacType'];
const errorIfMissingStorageFacility = ['storageFacility'];

const warnIfMissing = ['primeActualWeight', 'serviceOrderNumber', 'counselorRemarks', 'tacType', 'sacType'];
const showWhenCollapsed = ['counselorRemarks'];

const ordersLOA = {
  tac: '1111',
  sac: '2222222222',
  ntsTac: '3333',
  ntsSac: '4444444444',
};

const hhgInfo = {
  heading: 'HHG',
  shipmentId: 'testShipmentId394',
  requestedPickupDate: '26 Mar 2020',
  pickupAddress: {
    streetAddress1: '812 S 129th St',
    city: 'San Antonio',
    state: 'TX',
    postalCode: '78234',
  },
  destinationAddress: {
    streetAddress1: '441 SW Rio de la Plata Drive',
    city: 'Tacoma',
    state: 'WA',
    postalCode: '98421',
  },
};

const ntsInfo = {
  heading: 'NTS',
  requestedPickupDate: '26 Mar 2020',
  shipmentId: 'testShipmentId394',
  pickupAddress: {
    streetAddress1: '812 S 129th St',
    city: 'San Antonio',
    state: 'TX',
    postalCode: '78234',
  },
  destinationAddress: {
    streetAddress1: '441 SW Rio de la Plata Drive',
    city: 'Tacoma',
    state: 'WA',
    postalCode: '98421',
  },
};

const ntsReleaseInfo = {
  heading: 'NTS-release',
  shipmentId: 'testShipmentId111',
  ntsRecordedWeight: 2000,
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
  destinationAddress: {
    streetAddress1: '441 SW Rio de la Plata Drive',
    city: 'Tacoma',
    state: 'WA',
    postalCode: '98421',
  },
  secondaryDeliveryAddress: {
    streetAddress1: '812 S 129th St',
    city: 'San Antonio',
    state: 'TX',
    postalCode: '78234',
  },
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

const ntsReleaseMissingInfo = {
  heading: 'NTS-release',
  shipmentId: 'testShipmentId222',
  ntsRecordedWeight: 2000,
  requestedDeliveryDate: '26 Mar 2020',
  destinationAddress: {
    streetAddress1: '441 SW Rio de la Plata Drive',
    city: 'Tacoma',
    state: 'WA',
    postalCode: '98421',
  },
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

const postalOnlyInfo = {
  heading: 'HHG',
  requestedPickupDate: '26 Mar 2020',
  shipmentId: 'testShipmentId394',
  pickupAddress: {
    streetAddress1: '812 S 129th St',
    city: 'San Antonio',
    state: 'TX',
    postalCode: '78234',
  },
  destinationAddress: {
    postalCode: '98421',
  },
};

const diversionInfo = {
  heading: 'HHG',
  shipmentId: 'testShipmentId394',
  isDiversion: true,
  requestedPickupDate: '26 Mar 2020',
  pickupAddress: {
    streetAddress1: '812 S 129th St',
    city: 'San Antonio',
    state: 'TX',
    postalCode: '78234',
  },
  destinationAddress: {
    streetAddress1: '441 SW Rio de la Plata Drive',
    city: 'Tacoma',
    state: 'WA',
    postalCode: '98421',
  },
};

const cancelledInfo = {
  heading: 'HHG',
  shipmentId: 'testShipmentId394',
  isDiversion: false,
  shipmentStatus: 'CANCELED',
  requestedPickupDate: '26 Mar 2020',
  pickupAddress: {
    streetAddress1: '812 S 129th St',
    city: 'San Antonio',
    state: 'TX',
    postalCode: '78234',
  },
  destinationAddress: {
    streetAddress1: '441 SW Rio de la Plata Drive',
    city: 'Tacoma',
    state: 'WA',
    postalCode: '98421',
  },
};

export const HHGShipment = () => (
  <div style={{ padding: '20px' }}>
    <ShipmentDisplay displayInfo={hhgInfo} ordersLOA={ordersLOA} shipmentType={SHIPMENT_OPTIONS.HHG} isSubmitted />
  </div>
);

export const HHGShipmentServiceCounselor = () => (
  <div style={{ padding: '20px' }}>
    <ShipmentDisplay
      displayInfo={hhgInfo}
      ordersLOA={ordersLOA}
      shipmentType={SHIPMENT_OPTIONS.HHG}
      isSubmitted
      allowApproval={false}
    />
  </div>
);

export const HHGShipmentWithCounselorRemarks = () => (
  <div style={{ padding: '20px' }}>
    <ShipmentDisplay
      displayInfo={{ ...hhgInfo, counselorRemarks: 'counselor approved' }}
      shipmentType={SHIPMENT_OPTIONS.HHG}
      ordersLOA={ordersLOA}
      isSubmitted
    />
  </div>
);

export const HHGShipmentEditable = () => (
  <div style={{ padding: '20px' }}>
    <ShipmentDisplay
      displayInfo={{ ...hhgInfo, counselorRemarks: 'counselor approved' }}
      shipmentType={SHIPMENT_OPTIONS.HHG}
      ordersLOA={ordersLOA}
      isSubmitted
      editURL="/"
    />
  </div>
);

export const NTSShipment = () => (
  <div style={{ padding: '20px' }}>
    <ShipmentDisplay
      displayInfo={ntsInfo}
      ordersLOA={ordersLOA}
      shipmentType={SHIPMENT_OPTIONS.NTS}
      isSubmitted
      editURL="/"
    />
  </div>
);

export const NTSShipmentMissingInfo = () => (
  <div style={{ padding: '20px' }}>
    <ShipmentDisplay
      displayInfo={ntsMissingInfo}
      shipmentType={SHIPMENT_OPTIONS.NTS}
      shipmentId={ntsMissingInfo.shipmentId}
      ordersLOA={ordersLOA}
      isSubmitted
      warnIfMissing={warnIfMissing}
      errorIfMissing={errorIfMissingTACType}
      showWhenCollapsed={showWhenCollapsed}
      editURL="/"
    />
  </div>
);

export const NTSShipmentExternalVendor = () => (
  <div style={{ padding: '20px' }}>
    <ShipmentDisplay
      displayInfo={ntsInfo}
      shipmentType={SHIPMENT_OPTIONS.NTS}
      ordersLOA={ordersLOA}
      usesExternalVendor={usesExternalVendor}
      isSubmitted
      editURL="/"
    />
  </div>
);

export const NTSReleaseShipment = () => (
  <div style={{ padding: '20px' }}>
    <ShipmentDisplay
      displayInfo={ntsReleaseInfo}
      shipmentType={SHIPMENT_OPTIONS.NTSR}
      shipmentId={ntsReleaseInfo.shipmentId}
      ordersLOA={ordersLOA}
      showWhenCollapsed={showWhenCollapsed}
      isSubmitted
      editURL="/"
    />
  </div>
);

export const NTSReleaseShipmentExternalVendor = () => (
  <div style={{ padding: '20px' }}>
    <ShipmentDisplay
      displayInfo={ntsReleaseInfo}
      shipmentType={SHIPMENT_OPTIONS.NTSR}
      shipmentId={ntsReleaseInfo.shipmentId}
      ordersLOA={ordersLOA}
      showWhenCollapsed={showWhenCollapsed}
      usesExternalVendor={usesExternalVendor}
      isSubmitted
      editURL="/"
    />
  </div>
);

export const NTSReleaseShipmentMissingInfo = () => (
  <div style={{ padding: '20px' }}>
    <ShipmentDisplay
      displayInfo={ntsReleaseMissingInfo}
      shipmentType={SHIPMENT_OPTIONS.NTSR}
      shipmentId={ntsReleaseMissingInfo.shipmentId}
      ordersLOA={ordersLOA}
      isSubmitted
      warnIfMissing={warnIfMissing}
      errorIfMissing={errorIfMissingStorageFacility}
      showWhenCollapsed={showWhenCollapsed}
      editURL="/"
      allowApproval={false}
    />
  </div>
);

export const ApprovedShipment = () => (
  <div style={{ padding: '20px' }}>
    <ShipmentDisplay
      displayInfo={hhgInfo}
      ordersLOA={ordersLOA}
      shipmentType={SHIPMENT_OPTIONS.HHG}
      isSubmitted={false}
    />
  </div>
);

export const PostalOnlyDestination = () => (
  <div style={{ padding: '20px' }}>
    <ShipmentDisplay
      displayInfo={postalOnlyInfo}
      ordersLOA={ordersLOA}
      shipmentType={SHIPMENT_OPTIONS.HHG}
      isSubmitted
      editURL="/"
    />
  </div>
);

export const DivertedShipment = () => (
  <div style={{ padding: '20px' }}>
    <ShipmentDisplay
      shipmentId="1"
      displayInfo={diversionInfo}
      shipmentType={SHIPMENT_OPTIONS.HHG}
      ordersLOA={ordersLOA}
      isSubmitted
      editURL="/"
    />
  </div>
);

export const CancelledShipment = () => (
  <div style={{ padding: '20px' }}>
    <ShipmentDisplay
      shipmentId="1"
      displayInfo={cancelledInfo}
      ordersLOA={ordersLOA}
      shipmentType={SHIPMENT_OPTIONS.HHG}
      isSubmitted
    />
  </div>
);
