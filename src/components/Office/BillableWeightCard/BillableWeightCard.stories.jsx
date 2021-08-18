import React from 'react';

import BillableWeightCard from './BillableWeightCard';

export default {
  title: 'Office Components/BillableWeightCard',
  component: BillableWeightCard,
};

export const Card = () => (
  <div style={{ background: '#f9f9f9', width: '100%', height: '100%', paddingTop: 20 }}>
    <BillableWeightCard
      maxBillableWeight="13,750"
      totalBillableWeight="12,460"
      weightRequested="12,460"
      weightAllowance="8,000"
      shipments={[
        { id: '0001', shipmentType: 'HHG', billableWeightCap: '5,600' },
        { id: '0002', shipmentType: 'HHG', billableWeightCap: '3,200', reweigh: { id: '1234' } },
        { id: '0003', shipmentType: 'HHG', billableWeightCap: '3,400' },
      ]}
      entitlements={[
        { id: '1234', shipmentId: '0001', authorizedWeight: '4,600' },
        { id: '12346', shipmentId: '0002', authorizedWeight: '4,600' },
        { id: '12347', shipmentId: '0003', authorizedWeight: '4,600' },
      ]}
    />
  </div>
);
