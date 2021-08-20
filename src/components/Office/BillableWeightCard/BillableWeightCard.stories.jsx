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
        { id: '0001', shipmentType: 'HHG', billableWeightCap: '6,161', primeEstimatedWeight: '5,600' },
        {
          id: '0002',
          shipmentType: 'HHG',
          billableWeightCap: '3,200',
          primeEstimatedWeight: '5,000',
          reweigh: { id: '1234' },
        },
        { id: '0003', shipmentType: 'HHG', billableWeightCap: '3,400', primeEstimatedWeight: '5,000' },
      ]}
    />
  </div>
);
