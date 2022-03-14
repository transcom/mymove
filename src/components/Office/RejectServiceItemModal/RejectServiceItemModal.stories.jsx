import React from 'react';
import { action } from '@storybook/addon-actions';

import RejectServiceItemModal from './RejectServiceItemModal';

export default {
  title: 'Office Components/RejectServiceItemModal',
  component: RejectServiceItemModal,
};

const serviceItem = {
  id: 'abc123',
  serviceItem: 'Domestic Crating',
  code: 'DCRT',
  status: 'APPROVED',
  createdAt: '2020-10-31T00:00:00.12345',
  approvedAt: '2020-11-01T00:00:00.12345',
  details: {
    description: 'Trombone',
    itemDimensions: { length: 1000, width: 2500, height: 3000 },
    crateDimensions: { length: 1000, width: 2500, height: 3000 },
  },
};

export const Basic = () => (
  <div className="officeApp">
    <RejectServiceItemModal serviceItem={serviceItem} onSubmit={action('Submit')} onClose={action('Close')} />
  </div>
);
