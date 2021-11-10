import React from 'react';
import { action } from '@storybook/addon-actions';

import FinancialReviewButton from './FinancialReviewButton';

export default {
  title: 'Office Components/FinancialReviewButton',
  component: FinancialReviewButton,
};

export const MoveNotFlaggedForFinancialReview = () => <FinancialReviewButton onClick={action('Click')} />;

export const MoveFlaggedForFinancialReview = () => <FinancialReviewButton onClick={action('Click')} reviewRequested />;
