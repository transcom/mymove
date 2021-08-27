import React from 'react';

import WeightSummary from 'components/Office/WeightSummary/WeightSummary';

export default {
  title: 'Office Components/WeightSummary',
  component: WeightSummary,
};

const props = {
  maxBillableWeight: 13750,
  totalBillableWeight: 12460,
  weightRequested: 12460,
  weightAllowance: 8000,
  totalBillableWeightFlag: true,
  shipments: [
    { id: '0001', shipmentType: 'HHG', billableWeight: 6161, estimatedWeight: 5600 },
    {
      id: '0002',
      shipmentType: 'HHG',
      billableWeight: 3200,
      estimatedWeight: 5000,
      reweigh: { id: '1234' },
    },
    { id: '0003', shipmentType: 'HHG', billableWeight: 3400, estimatedWeight: 5000 },
  ],
};

export const WeightSummaryCard = () => (
  <div style={{ maxWidth: '336px' }}>
    <WeightSummary {...props} />
  </div>
);
