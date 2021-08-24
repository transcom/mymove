/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';

import ShipmentList from '.';

import { SHIPMENT_OPTIONS } from 'shared/constants';

export const Basic = () => (
  <div className="grid-container">
    <h3>Single Shipment</h3>
    <ShipmentList shipments={[{ id: '0001', shipmentType: SHIPMENT_OPTIONS.PPM }]} onShipmentClick={() => {}} />
    <br />
    <h3>Multiple shipments</h3>
    <ShipmentList
      shipments={[
        { id: '0001', shipmentType: SHIPMENT_OPTIONS.HHG },
        { id: '0002', shipmentType: SHIPMENT_OPTIONS.NTS },
        { id: '0003', shipmentType: SHIPMENT_OPTIONS.PPM },
      ]}
      onShipmentClick={() => {}}
    />
  </div>
);

export const ShipmentListWithWeights = () => (
  <div className="grid-container">
    <h3>Single Shipment</h3>
    <ShipmentList
      shipments={[{ id: '0001', shipmentType: SHIPMENT_OPTIONS.HHG, billableWeightCap: '4,600' }]}
      showShipmentWeight
    />
    <br />
    <h3>Multiple shipments</h3>
    <ShipmentList
      shipments={[
        { id: '0001', shipmentType: SHIPMENT_OPTIONS.HHG, billableWeightCap: '6,161', primeEstimatedWeight: '5,600' },
        { id: '0002', shipmentType: SHIPMENT_OPTIONS.HHG, billableWeightCap: '3,200', reweigh: { id: '1234' } },
        { id: '0003', shipmentType: SHIPMENT_OPTIONS.HHG, billableWeightCap: '3,400', primeEstimatedWeight: '5,000' },
      ]}
      showShipmentWeight
    />
  </div>
);

export default {
  title: 'Components / ShipmentList',
};
