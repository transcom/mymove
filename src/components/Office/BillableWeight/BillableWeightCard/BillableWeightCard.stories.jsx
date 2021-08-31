import React from 'react';

import BillableWeightCard from './BillableWeightCard';

export default {
  title: 'Office Components/BillableWeightCard',
  component: BillableWeightCard,
  argTypes: {
    onReviewWeights: { action: 'review weights' },
  },
};

export const Card = (argTypes) => (
  <div style={{ background: '#f9f9f9', width: '100%', height: '100%', paddingTop: 20 }}>
    <BillableWeightCard
      maxBillableWeight={13750}
      totalBillableWeight={12460}
      weightRequested={12460}
      weightAllowance={8000}
      onReviewWeights={argTypes.onReviewWeights}
      shipments={[
        { id: '0001', shipmentType: 'HHG', billableWeight: 6161, estimatedWeight: 5600 },
        {
          id: '0002',
          shipmentType: 'HHG',
          billableWeight: 3200,
          estimatedWeight: 5000,
          reweigh: { id: '1234' },
        },
        { id: '0003', shipmentType: 'HHG', billableWeight: 3400, estimatedWeight: 5000 },
      ]}
    />
  </div>
);
