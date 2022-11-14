import React from 'react';

import ShipmentQAEReportHeader from './ShipmentQAEReportHeader';

import { SHIPMENT_OPTIONS } from 'shared/constants';
import { MockProviders } from 'testUtils';
import { permissionTypes } from 'constants/permissions';

export default {
  title: 'Office Components/ShipmentQAEReportHeader',
  component: ShipmentQAEReportHeader,
  decorators: [
    (Story, context) => {
      // Dont wrap with permissions for the read only tests
      if (context.name.includes('Read Only')) {
        return <MockProviders>{Story()}</MockProviders>;
      }

      // By default, show component with permissions
      return (
        <MockProviders permissions={[permissionTypes.createEvaluationReport]}>
          <Story />
        </MockProviders>
      );
    },
  ],
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
    <ShipmentQAEReportHeader shipment={hhgShipment} />
  </div>
);

export const nts = () => (
  <div className="officeApp">
    <ShipmentQAEReportHeader shipment={ntsShipment} />
  </div>
);
export const ntsr = () => (
  <div className="officeApp">
    <ShipmentQAEReportHeader shipment={ntsrShipment} />
  </div>
);
export const ppm = () => (
  <div className="officeApp">
    <ShipmentQAEReportHeader shipment={ppmShipment} />
  </div>
);
export const hhgReadOnly = () => (
  <div className="officeApp">
    <ShipmentQAEReportHeader shipment={hhgShipment} />
  </div>
);

export const ntsReadOnly = () => (
  <div className="officeApp">
    <ShipmentQAEReportHeader shipment={ntsShipment} />
  </div>
);
export const ntsrReadOnly = () => (
  <div className="officeApp">
    <ShipmentQAEReportHeader shipment={ntsrShipment} />
  </div>
);
export const ppmReadOnly = () => (
  <div className="officeApp">
    <ShipmentQAEReportHeader shipment={ppmShipment} />
  </div>
);
