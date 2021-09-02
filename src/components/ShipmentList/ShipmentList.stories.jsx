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
      shipments={[{ id: '0001', shipmentType: SHIPMENT_OPTIONS.HHG, billableWeight: 4600 }]}
      showShipmentWeight
    />
    <br />
    <h3>Multiple shipments</h3>
    <ShipmentList
      shipments={[
        { id: '0001', shipmentType: SHIPMENT_OPTIONS.HHG, billableWeight: 6161, estimatedWeight: 5600 },
        { id: '0002', shipmentType: SHIPMENT_OPTIONS.HHG, billableWeight: 3200, reweigh: { id: '1234' } },
        { id: '0003', shipmentType: SHIPMENT_OPTIONS.HHG, billableWeight: 3400, estimatedWeight: 5000 },
      ]}
      showShipmentWeight
    />
  </div>
);

export default {
  title: 'Components / ShipmentList',
};
