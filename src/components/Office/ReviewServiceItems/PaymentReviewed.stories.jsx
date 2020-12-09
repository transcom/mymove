import React from 'react';

import PaymentReviewed from './PaymentReviewed';

export default {
  title: 'Office Components/ReviewServiceItems/ReviewDetails',
  component: PaymentReviewed,
};

export const PaymentAuthorized = () => (
  <PaymentReviewed authorizedAmount={499.99} dateAuthorized="2020-01-02T23:59:59.456Z" />
);

export const PaymentAuthorizedRejected = () => (
  <PaymentReviewed authorizedAmount={0} dateAuthorized="2020-01-02T23:59:59.456Z" />
);
