/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';

import ShipmentList, { ShipmentListItem } from '.';

import { SHIPMENT_OPTIONS } from 'shared/constants';

export const Basic = () => (
  <div className="grid-container">
    <h3>Single Shipment</h3>
    <ShipmentList shipments={[{ id: '0001', shipmentType: SHIPMENT_OPTIONS.PPM }]} onShipmentClick={() => {}} />
    <ShipmentListItem
      shipment={{ id: '0001', shipmentType: SHIPMENT_OPTIONS.HHG, primeActualWeight: 4200 }}
      showNumber={false}
      showShipmentWeight
      isOverweight
      onShipmentClick={() => {}}
    />
    <ShipmentListItem
      shipment={{ id: '0001', shipmentType: SHIPMENT_OPTIONS.HHG, primeActualWeight: 6800 }}
      showShipmentWeight
      isMissingWeight
      showNumber={false}
      onShipmentClick={() => {}}
    />
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

export default {
  title: 'Customer Components / ShipmentList',
};
