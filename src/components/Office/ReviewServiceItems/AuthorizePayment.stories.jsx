import React from 'react';
import { action } from '@storybook/addon-actions';

import AuthorizePaymentComponent from './AuthorizePayment';

export default {
  title: 'TOO/TIO Components|ReviewServiceItems/ReviewDetails',
  component: AuthorizePaymentComponent,
};

export const AuthorizePayment = () => (
  <AuthorizePaymentComponent handleFinishReviewBtn={action('Finish button clicked!')} />
);
