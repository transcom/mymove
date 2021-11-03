import React from 'react';
import { action } from '@storybook/addon-actions';

import FinancialReviewModal from './FinancialReviewModal';

export default {
  title: 'Office Components/FinancialReviewModal',
  component: FinancialReviewModal,
};

export const Basic = () => <FinancialReviewModal onSubmit={action('Submit')} onClose={action('Cancel')} />;
