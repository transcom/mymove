import React from 'react';

import ShipmentCard from './ShipmentCard';

export default {
  title: 'Office Components/BillableWeightShipmentCard',
  component: ShipmentCard,
};

const props = {
  billableWeight: 4014,
  dateReweighRequested: new Date().toISOString(),
  departedDate: new Date().toISOString(),
  pickupAddress: {
    city: 'Rancho Santa Margarita',
    state: 'CA',
    postal_code: '92688',
  },
  destinationAddress: {
    city: 'West Springfield Town',
    state: 'MA',
    postal_code: '01089',
  },
  estimatedWeight: 5000,
  originalWeight: 4014,
  reweighRemarks: 'Unable to perform reweigh because shipment was already unloaded',
  reweightWeight: '',
};

export const Card = () => (
  <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '100%' }}>
    <ShipmentCard {...props} />
  </div>
);
