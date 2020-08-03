import React from 'react';
import { action } from '@storybook/addon-actions';

import AuthorizePaymentComponent from './AuthorizePayment';

export default {
  title: 'TOO/TIO Components|ReviewServiceItems/ReviewDetails',
  component: AuthorizePaymentComponent,
  parameters: {
    loki: {
      skip: true,
    },
  },
};

export const AuthorizePayment = () => <AuthorizePaymentComponent onClick={action('Finish button clicked!')} />;
