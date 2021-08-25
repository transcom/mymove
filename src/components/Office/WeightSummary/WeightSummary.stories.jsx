import React from 'react';

import WeightSummary from 'components/Office/WeightSummary/WeightSummary';

export default {
  title: 'Office Components/WeightSummary',
  component: WeightSummary,
};

export const WeightSummaryCard = () => (
  <WeightSummary
    maxBillableWeight="13,750"
    totalBillableWeight="12,460"
    weightRequested="12,460"
    weightAllowance="8,000"
    totalBillableWeightFlag
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
);
