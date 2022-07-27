import React from 'react';
import { render, screen } from '@testing-library/react';

import EvaluationReportShipmentInfo from './EvaluationReportShipmentInfo';

import { SHIPMENT_OPTIONS } from 'shared/constants';

const pickupAddress = {
  city: 'Beverly Hills',
  country: 'US',
  postalCode: '90210',
  state: 'CA',
  streetAddress1: '123 Any Street',
  streetAddress2: 'P.O. Box 12345',
  streetAddress3: 'c/o Some Person',
};
const destinationAddress = {
  city: 'Fairfield',
  country: 'US',
  postalCode: '94535',
  state: 'CA',
  streetAddress1: '987 Any Avenue',
  streetAddress2: 'P.O. Box 9876',
  streetAddress3: 'c/o Some Person',
};

const hhgShipment = {
  id: 'c3c64a08-778d-4f9f-8b67-b2502e0fb5e9',
  pickupAddress,
  destinationAddress,
  shipmentType: SHIPMENT_OPTIONS.HHG,
  status: 'SUBMITTED',
  createdAt: '2022-07-12T19:38:35.886Z',
};
const ntsShipment = {
  id: 'c3c64a08-778d-4f9f-8b67-b2502e0fb5e9',
  pickupAddress,
  shipmentType: SHIPMENT_OPTIONS.NTS,
  storageFacility: {
    facilityName: 'Storage Facility',
    address: {},
  },
  status: 'SUBMITTED',
  createdAt: '2022-07-12T19:38:35.886Z',
};
const ntsReleaseShipment = {
  id: 'c3c64a08-778d-4f9f-8b67-b2502e0fb5e9',
  destinationAddress,
  shipmentType: SHIPMENT_OPTIONS.NTSR,
  storageFacility: {
    facilityName: 'Storage Facility',
    address: {},
  },
  status: 'SUBMITTED',
  createdAt: '2022-07-12T19:38:35.886Z',
};
const ppmShipment = {
  id: 'c3c64a08-778d-4f9f-8b67-b2502e0fb5e9',
  shipmentType: SHIPMENT_OPTIONS.PPM,
  ppmShipment: {
    pickupPostalCode: '90210',
    destinationPostalCode: '94535',
  },
  status: 'SUBMITTED',
  createdAt: '2022-07-12T19:38:35.886Z',
};
describe('EvaluationReportShipmentInfo', () => {
  it('renders HHG shipment', () => {
    render(<EvaluationReportShipmentInfo shipment={hhgShipment} shipmentNumber={1} />);
    expect(screen.getByRole('heading', { level: 4, name: /HHG Shipment ID #C3C64/ })).toBeInTheDocument();
    expect(screen.getByText(/123 Any Street/)).toBeInTheDocument();
    expect(screen.getByText(/987 Any Avenue/)).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Create report' })).toBeInTheDocument();
  });
  it('renders NTS shipment', () => {
    render(<EvaluationReportShipmentInfo shipment={ntsShipment} shipmentNumber={1} />);
    expect(screen.getByRole('heading', { level: 4, name: /NTS Shipment ID #C3C64/ })).toBeInTheDocument();
    expect(screen.getByText(/123 Any Street/)).toBeInTheDocument();
    expect(screen.getByText(/Storage Facility/)).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Create report' })).toBeInTheDocument();
  });
  it('renders NTS-R shipment', () => {
    render(<EvaluationReportShipmentInfo shipment={ntsReleaseShipment} shipmentNumber={1} />);
    expect(screen.getByRole('heading', { level: 4, name: /NTS-Release Shipment ID #C3C64/ })).toBeInTheDocument();
    expect(screen.getByText(/Storage Facility/)).toBeInTheDocument();
    expect(screen.getByText(/987 Any Avenue/)).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Create report' })).toBeInTheDocument();
  });
  it('renders PPM shipment', () => {
    render(<EvaluationReportShipmentInfo shipment={ppmShipment} shipmentNumber={1} />);
    expect(screen.getByRole('heading', { level: 4, name: /PPM Shipment ID #C3C64/ })).toBeInTheDocument();
    expect(screen.getByText(/90210/)).toBeInTheDocument();
    expect(screen.getByText(/94535/)).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Create report' })).toBeInTheDocument();
  });
});
