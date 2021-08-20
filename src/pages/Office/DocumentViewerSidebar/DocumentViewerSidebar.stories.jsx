import React from 'react';
import { Button } from '@trussworks/react-uswds';

import DocumentViewerSidebar from './DocumentViewerSidebar';

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
