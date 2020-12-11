import React from 'react';
import { action } from '@storybook/addon-actions';

import AuthorizePaymentComponent from './AuthorizePayment';

export default {
  title: 'Office Components/ReviewServiceItems/ReviewDetails',
  component: AuthorizePaymentComponent,
};

export const AuthorizePayment = () => <AuthorizePaymentComponent onClick={action('Finish button clicked!')} />;
