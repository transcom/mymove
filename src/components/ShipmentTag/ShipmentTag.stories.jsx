import React from 'react';

import ShipmentTag from './ShipmentTag';

import { SHIPMENT_OPTIONS } from 'shared/constants';

export default {
  title: 'Components / Shipment Tag',
  component: ShipmentTag,
};

export const HHG = () => <ShipmentTag shipmentType={SHIPMENT_OPTIONS.HHG} />;
export const HHGWithNumber = () => <ShipmentTag shipmentType={SHIPMENT_OPTIONS.HHG} shipmentNumber="1A2B3C" />;
export const PPM = () => <ShipmentTag shipmentType={SHIPMENT_OPTIONS.PPM} />;
export const NTS = () => <ShipmentTag shipmentType={SHIPMENT_OPTIONS.NTS} />;
export const NTSR = () => <ShipmentTag shipmentType={SHIPMENT_OPTIONS.NTSR} />;
