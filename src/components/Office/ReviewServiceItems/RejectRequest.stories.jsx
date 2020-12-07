import React from 'react';
import { action } from '@storybook/addon-actions';

import RejectRequestComponent from './RejectRequest';

export default {
  title: 'TOO/TIO Components/ReviewServiceItems/ReviewDetails',
  component: RejectRequestComponent,
};

export const RejectRequest = () => (
  <RejectRequestComponent numberOfItems={1} onClick={action('Finish button clicked!')} />
);
