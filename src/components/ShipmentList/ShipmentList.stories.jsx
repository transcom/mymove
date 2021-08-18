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
      entitlements={[{ id: '1234', shipmentId: '0001', authorizedWeight: '4,600' }]}
      showShipmentWeight
      onShipmentClick={() => {}}
    />
    <br />
    <h3>Multiple shipments</h3>
    <ShipmentList
      shipments={[
        { id: '0001', shipmentType: SHIPMENT_OPTIONS.HHG, billableWeightCap: '5,600' },
        { id: '0002', shipmentType: SHIPMENT_OPTIONS.HHG, billableWeightCap: '3,200', reweigh: { id: '1234' } },
        { id: '0003', shipmentType: SHIPMENT_OPTIONS.HHG, billableWeightCap: '3,200' },
      ]}
      entitlements={[
        { id: '1234', shipmentId: '0001', authorizedWeight: '4,600' },
        { id: '12346', shipmentId: '0002', authorizedWeight: '4,600' },
        { id: '12347', shipmentId: '0003', authorizedWeight: '4,600' },
      ]}
      showShipmentWeight
      onShipmentClick={() => {}}
    />
  </div>
);

export default {
  title: 'Customer Components / ShipmentList',
};
