import React from 'react';
import { render, screen } from '@testing-library/react';

import ShipmentQAEReportHeader from './ShipmentQAEReportHeader';

import { SHIPMENT_OPTIONS } from 'shared/constants';
import { MockProviders } from 'testUtils';
import { permissionTypes } from 'constants/permissions';

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
  shipmentLocator: 'EVLRPT-01',
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
  shipmentLocator: 'EVLRPT-02',
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
  shipmentLocator: 'EVLRPT-03',
};
const ppmShipment = {
  id: 'c3c64a08-778d-4f9f-8b67-b2502e0fb5e9',
  shipmentType: SHIPMENT_OPTIONS.PPM,
  ppmShipment: {
    pickupAddress,
    destinationAddress,
  },
  status: 'SUBMITTED',
  createdAt: '2022-07-12T19:38:35.886Z',
  shipmentLocator: 'EVLRPT-04',
};
describe('ShipmentQAEReportHeader', () => {
  it('renders HHG shipment', () => {
    render(
      <MockProviders permissions={[permissionTypes.createEvaluationReport]}>
        <ShipmentQAEReportHeader destinationDutyLocationPostalCode="" shipment={hhgShipment} shipmentNumber={1} />
      </MockProviders>,
    );
    expect(screen.getByTestId('shipmentHeader')).toHaveTextContent('EVLRPT-01');

    expect(screen.getByTestId('shipmentHeader')).toHaveTextContent(pickupAddress.streetAddress1);
    expect(screen.getByTestId('shipmentHeader')).toHaveTextContent(destinationAddress.streetAddress1);
    expect(screen.getByRole('button', { name: 'Create report' })).toBeVisible();
  });
  it('renders NTS shipment', () => {
    render(
      <MockProviders permissions={[permissionTypes.createEvaluationReport]}>
        <ShipmentQAEReportHeader destinationDutyLocationPostalCode="" shipment={ntsShipment} shipmentNumber={1} />
      </MockProviders>,
    );
    expect(screen.getByTestId('shipmentHeader')).toHaveTextContent('EVLRPT-02');

    expect(screen.getByTestId('shipmentHeader')).toHaveTextContent(ntsShipment.pickupAddress.streetAddress1);
    expect(screen.getByTestId('shipmentHeader')).toHaveTextContent(ntsShipment.storageFacility.facilityName);
    expect(screen.getByRole('button', { name: 'Create report' })).toBeInTheDocument();
  });
  it('renders NTS-R shipment', () => {
    render(
      <MockProviders permissions={[permissionTypes.createEvaluationReport]}>
        <ShipmentQAEReportHeader
          destinationDutyLocationPostalCode=""
          shipment={ntsReleaseShipment}
          shipmentNumber={1}
        />
      </MockProviders>,
    );
    expect(screen.getByTestId('shipmentHeader')).toHaveTextContent('EVLRPT-03');

    expect(screen.getByTestId('shipmentHeader')).toHaveTextContent(ntsReleaseShipment.storageFacility.facilityName);
    expect(screen.getByTestId('shipmentHeader')).toHaveTextContent(
      ntsReleaseShipment.destinationAddress.streetAddress1,
    );
    expect(screen.getByRole('button', { name: 'Create report' })).toBeVisible();
  });
  it('renders PPM shipment', () => {
    render(
      <MockProviders permissions={[permissionTypes.createEvaluationReport]}>
        <ShipmentQAEReportHeader destinationDutyLocationPostalCode="" shipment={ppmShipment} shipmentNumber={1} />
      </MockProviders>,
    );
    expect(screen.getByTestId('shipmentHeader')).toHaveTextContent('EVLRPT-04');

    expect(screen.getByTestId('shipmentHeader')).toHaveTextContent(ppmShipment.ppmShipment.pickupAddress.postalCode);
    expect(screen.getByTestId('shipmentHeader')).toHaveTextContent(
      ppmShipment.ppmShipment.destinationAddress.postalCode,
    );
    expect(screen.getByRole('button', { name: 'Create report' })).toBeVisible();
  });
  it('renders a shipment but disables button when move is locked', () => {
    render(
      <MockProviders permissions={[permissionTypes.createEvaluationReport]}>
        <ShipmentQAEReportHeader
          destinationDutyLocationPostalCode=""
          shipment={hhgShipment}
          shipmentNumber={1}
          isMoveLocked
        />
      </MockProviders>,
    );

    expect(screen.getByRole('button', { name: 'Create report' })).toBeVisible();
    expect(screen.getByRole('button', { name: 'Create report' })).toBeDisabled();
  });
});
