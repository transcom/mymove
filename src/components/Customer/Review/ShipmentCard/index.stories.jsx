/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';

import ShipmentCard from '.';

const defaultProps = {};

export default {
  title: 'Customer Components | ShipmentCard',
};

export const Basic = () => (
  <div style={{ padding: 40 }}>
    <ShipmentCard {...defaultProps} />
  </div>
);
