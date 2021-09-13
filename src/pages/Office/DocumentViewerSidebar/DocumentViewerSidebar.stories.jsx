import React from 'react';
import { Button } from '@trussworks/react-uswds';

import DocumentViewerSidebar from './DocumentViewerSidebar';

import ShipmentCard from 'components/Office/BillableWeight/ShipmentCard/ShipmentCard';
import WeightSummary from 'components/Office/WeightSummary/WeightSummary';

export default {
  title: 'Office Components/DocumentViewerSidebar',
  component: DocumentViewerSidebar,
};

export const Sidebar = () => (
  <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center' }}>
    <DocumentViewerSidebar
      title="Review weights"
      subtitle="Shipment weights"
      description="Shipment 1 of 2"
      onClose={() => {}}
    >
      <DocumentViewerSidebar.Content />
      <DocumentViewerSidebar.Footer>
        <Button>Review billable weight</Button>
      </DocumentViewerSidebar.Footer>
    </DocumentViewerSidebar>
  </div>
);

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

const shipmentCardProps = {
  billableWeight: 4014,
  dateReweighRequested: new Date('1/1/2020').toISOString(),
  departedDate: new Date('12/25/2019').toISOString(),
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
};

export const WeightsSidebar = () => (
  <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center' }}>
    <DocumentViewerSidebar
      title="Review weights"
      subtitle="Shipment weights"
      description="Shipment 1 of 2"
      onClose={() => {}}
    >
      <DocumentViewerSidebar.Content>
        <div style={{ maxWidth: '336px', backgroundColor: 'white', marginBottom: '16px' }}>
          <WeightSummary {...props} />
        </div>
        <div style={{ height: '100%', width: 336 }}>
          <ShipmentCard {...shipmentCardProps} />
        </div>
      </DocumentViewerSidebar.Content>
      <DocumentViewerSidebar.Footer>
        <Button>Review billable weight</Button>
      </DocumentViewerSidebar.Footer>
    </DocumentViewerSidebar>
  </div>
);
