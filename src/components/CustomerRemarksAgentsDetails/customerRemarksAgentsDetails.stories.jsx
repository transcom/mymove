import React from 'react';
import { withKnobs, text, object } from '@storybook/addon-knobs';

import CustomerRemarksAgentsDetails from 'components/CustomerRemarksAgentsDetails';

export default {
  title: 'TOO/TIO Components/CustomerRemarksAgentsDetails',
  decorator: withKnobs,
};

export const WithNoDetails = () => (
  <div id="containers" style={{ padding: '20px' }}>
    <CustomerRemarksAgentsDetails />
  </div>
);

const customerRemarks = 'This is a customer remark.';
const releasingAgent = {
  firstName: 'FirstName',
  lastName: 'LastName',
  phone: '(111) 111-1111',
  email: 'test@test.com',
};
const receivingAgent = {
  firstName: 'FirstName',
  lastName: 'LastName',
  phone: '(111) 111-1111',
  email: 'test@test.com',
};
export const WithDetails = () => (
  <div id="containers" style={{ padding: '20px' }}>
    <CustomerRemarksAgentsDetails
      customerRemarks={text('customerRemarks', customerRemarks)}
      releasingAgent={object('releasingAgent', releasingAgent)}
      receivingAgent={object('receivingAgent', receivingAgent)}
    />
  </div>
);
