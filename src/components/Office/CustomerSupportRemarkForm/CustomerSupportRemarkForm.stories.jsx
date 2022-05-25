import React from 'react';
import { MemoryRouter, Route } from 'react-router';

import CustomerSupportRemarkForm from './CustomerSupportRemarkForm';

import { MockProviders } from 'testUtils';

export default {
  title: 'Office Components/CustomerSupportRemarkForm',
  decorators: [
    (Story) => (
      <MemoryRouter initialEntries={['/moves/AMDORD/customer-support-remarks']}>
        <Route path="/moves/:moveCode/customer-support-remarks">
          <MockProviders>
            <div style={{ padding: '40px', width: '550px' }}>
              <Story />
            </div>
          </MockProviders>
        </Route>
      </MemoryRouter>
    ),
  ],
};

export const Default = () => <CustomerSupportRemarkForm />;
