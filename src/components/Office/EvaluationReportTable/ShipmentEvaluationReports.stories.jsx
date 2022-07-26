import React from 'react';

import ShipmentEvaluationReports from './ShipmentEvaluationReports';

import { SHIPMENT_OPTIONS } from 'shared/constants';

export default {
  title: 'Office Components/ShipmentEvaluationReports',
  component: ShipmentEvaluationReports,
};

const hhgShipment = {
  id: '11111111-1111-1111-1111-111111111111',
  shipmentType: SHIPMENT_OPTIONS.HHG,
  pickupAddress: {
    streetAddress1: '123 Any St',
    city: 'Anytown',
    state: 'AK',
    postalCode: '90210',
  },
  destinationAddress: {
    streetAddress1: '123 Any St',
    city: 'Anytown',
    state: 'AK',
    postalCode: '90210',
  },
};
const ppmShipment = {
  id: '22222222-2222-2222-2222-222222222222',
  shipmentType: SHIPMENT_OPTIONS.PPM,
  ppmShipment: {
    pickupPostalCode: '89503',
    destinationPostalCode: '90210',
  },
};

const ntsShipment = {
  id: '33333333-3333-3333-3333-333333333333',
  shipmentType: SHIPMENT_OPTIONS.NTS,
  pickupAddress: {
    streetAddress1: '123 Any St',
    city: 'Anytown',
    state: 'AK',
    postalCode: '90210',
  },
  storageFacility: {
    facilityName: 'Awesome Storage LLC',
  },
};
const ntsrShipment = {
  id: '44444444-4444-4444-4444-444444444444',
  shipmentType: SHIPMENT_OPTIONS.NTSR,
  destinationAddress: {
    streetAddress1: '123 Any St',
    city: 'Anytown',
    state: 'AK',
    postalCode: '90210',
  },
  storageFacility: {
    facilityName: 'Awesome Storage LLC',
  },
};

const shipments = [hhgShipment, ppmShipment, ntsShipment, ntsrShipment];

const reports = [
  {
    id: 'a7fdb0b3-f876-4686-b94f-ad20a2c9a63d',
    createdAt: '2022-07-14T19:21:27.573Z',
    location: 'DESTINATION',
    shipmentID: '11111111-1111-1111-1111-111111111111',
    submittedAt: '2022-07-14T19:21:27.565Z',
    type: 'SHIPMENT',
    violationsObserved: true,
  },
  {
    id: '1f9d18a8-7688-428d-be8e-3f3c59a0ae5e',
    createdAt: '2022-07-14T19:21:27.579Z',
    location: null,
    shipmentID: '22222222-2222-2222-2222-222222222222',
    submittedAt: null,
    type: 'SHIPMENT',
    violationsObserved: true,
  },
];
export const empty = () => (
  <div className="officeApp">
    <ShipmentEvaluationReports reports={[]} />
  </div>
);

export const single = () => (
  <div className="officeApp">
    <ShipmentEvaluationReports shipments={shipments} reports={reports} />
  </div>
);
