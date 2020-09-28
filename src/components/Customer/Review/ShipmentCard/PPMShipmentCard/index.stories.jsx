/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';

import PPMShipmentCard from '.';

const defaultProps = {
  destinationZIP: '11111',
  estimatedWeight: '5,000',
  expectedDepartureDate: new Date('01/01/2020').toISOString(),
  shipmentId: '#ABC123K-001',
  sitDays: '24',
  originZIP: '00000',
};

export default {
  title: 'Customer Components | PPMShipmentCard',
};

export const Basic = () => (
  <div style={{ padding: 10 }}>
    <PPMShipmentCard {...defaultProps} />
  </div>
);
