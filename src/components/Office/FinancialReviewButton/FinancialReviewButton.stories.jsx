import React from 'react';
import { action } from '@storybook/addon-actions';

import FinancialReviewButton from './FinancialReviewButton';

export default {
  title: 'Office Components/FinancialReviewButton',
  component: FinancialReviewButton,
};

export const Basic = () => <FinancialReviewButton onClick={action('Click')} />;
