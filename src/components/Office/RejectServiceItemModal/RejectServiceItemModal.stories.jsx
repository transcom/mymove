import React from 'react';
import { action } from '@storybook/addon-actions';

import RejectServiceItemModal from './RejectServiceItemModal';

export default {
  title: 'TOO/TIO Components|RejectServiceItemModal',
  component: RejectServiceItemModal,
};

const serviceItem = {
  id: 'abc123',
  serviceItem: 'Domestic Crating',
  code: 'DCRT',
  status: 'SUBMITTED',
  submittedAt: '2020-10-31',
  details: {
    description: '',
    itemDimensions: {},
    crateDimensions: {},
    imgURL: '',
  },
};

export const Basic = () => (
  <RejectServiceItemModal serviceItem={serviceItem} onSubmit={action('Submit')} onClose="Close" />
);
