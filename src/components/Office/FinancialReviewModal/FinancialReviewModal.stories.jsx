import React from 'react';
import { action } from '@storybook/addon-actions';

import FinancialReviewModal from './FinancialReviewModal';

export default {
  title: 'Office Components/FinancialReviewModal',
  component: FinancialReviewModal,
};

export const Basic = () => <FinancialReviewModal Submit={action('Submit')} onClose={action('Cancel')} />;
