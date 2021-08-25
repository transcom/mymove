import React from 'react';

import ShipmentCard from './ShipmentCard';

export default {
  title: 'Office Components/BillableWeightShipmentCard',
  component: ShipmentCard,
};

export const Card = () => (
  <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '100%' }}>
    <ShipmentCard />
  </div>
);
