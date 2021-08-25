import React from 'react';

import BillableWeightCard from './BillableWeightCard';

export default {
  title: 'Office Components/BillableWeightCard',
  component: BillableWeightCard,
  argTypes: {
    reviewWeights: { action: 'review weights' },
  },
};

export const Card = (argTypes) => (
  <div style={{ background: '#f9f9f9', width: '100%', height: '100%', paddingTop: 20 }}>
    <BillableWeightCard
      maxBillableWeight={13750}
      totalBillableWeight={12460}
      weightRequested={12460}
      weightAllowance={8000}
      reviewWeights={argTypes.reviewWeights}
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
