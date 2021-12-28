import React from 'react';
import { action } from '@storybook/addon-actions';

import FinancialReviewButton from './FinancialReviewButton';

export default {
  title: 'Office Components/FinancialReviewButton',
  component: FinancialReviewButton,
};

export const MoveNotFlaggedForFinancialReview = () => (
  <div className="officeApp">
    <FinancialReviewButton onClick={action('Click')} />
  </div>
);

export const MoveFlaggedForFinancialReview = () => (
  <div className="officeApp">
    <FinancialReviewButton onClick={action('Click')} reviewRequested />
  </div>
);
