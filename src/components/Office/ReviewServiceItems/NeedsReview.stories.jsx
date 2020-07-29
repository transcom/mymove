import React from 'react';
import { action } from '@storybook/addon-actions';

import NeedsReviewComponent from './NeedsReview';

export default {
  title: 'TOO/TIO Components|ReviewServiceItems',
  component: NeedsReviewComponent,
};

export const NeedsReview = () => (
  <NeedsReviewComponent numberOfItems={1} handleFinishReviewBtn={action('Finish button clicked!')} />
);
