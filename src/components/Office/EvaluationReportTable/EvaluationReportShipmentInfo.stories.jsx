import React from 'react';

import EvaluationReportShipmentInfo from './EvaluationReportShipmentInfo';

import { SHIPMENT_OPTIONS } from 'shared/constants';
import { MockProviders } from 'testUtils';

export default {
  title: 'Office Components/EvaluationReportShipmentInfo',
  component: EvaluationReportShipmentInfo,
  decorators: [(Story) => <MockProviders>{Story()}</MockProviders>],
};

const hhgShipment = {
  id: '111111',
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
  id: '22222',
  shipmentType: SHIPMENT_OPTIONS.PPM,
  ppmShipment: {
    pickupPostalCode: '89503',
    destinationPostalCode: '90210',
  },
};

const ntsShipment = {
  id: '3333333',
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
  id: '444444',
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

export const hhg = () => (
  <div className="officeApp">
    <EvaluationReportShipmentInfo shipment={hhgShipment} shipmentNumber={1} />
  </div>
);

export const nts = () => (
  <div className="officeApp">
    <EvaluationReportShipmentInfo shipment={ntsShipment} shipmentNumber={2} />
  </div>
);
export const ntsr = () => (
  <div className="officeApp">
    <EvaluationReportShipmentInfo shipment={ntsrShipment} shipmentNumber={3} />
  </div>
);
export const ppm = () => (
  <div className="officeApp">
    <EvaluationReportShipmentInfo shipment={ppmShipment} shipmentNumber={4} />
  </div>
);
