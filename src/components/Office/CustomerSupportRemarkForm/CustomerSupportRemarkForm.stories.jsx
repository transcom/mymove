import React from 'react';

import CustomerSupportRemarkForm from './CustomerSupportRemarkForm';

import { MockProviders } from 'testUtils';

export default {
  title: 'Office Components/CustomerSupportRemarkForm',
  decorators: [
    (Story) => (
      <MockProviders>
        <div style={{ padding: '40px', width: '550px' }}>
          <Story />
        </div>
      </MockProviders>
    ),
  ],
};

export const Default = () => <CustomerSupportRemarkForm />;
