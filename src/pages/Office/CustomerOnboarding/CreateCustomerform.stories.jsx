import React from 'react';

import CreateCustomerForm from './CreateCustomerForm';

import { MockProviders } from 'testUtils';

export default {
  title: 'Office Components/CreateCustomerForm',
  parameters: { layout: 'fullscreen' },
};

export const Form = () => {
  return (
    <MockProviders>
      <CreateCustomerForm />
    </MockProviders>
  );
};
