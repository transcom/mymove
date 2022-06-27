import React from 'react';

import CustomerSupportRemarkText from './CustomerSupportRemarkText';

import { MockProviders } from 'testUtils';

export default {
  title: 'Office Components/CustomerSupportRemarkText',
  component: CustomerSupportRemarkText,
  decorators: [
    (Story) => (
      <div style={{ padding: '40px', width: '550px', minWidth: '550px' }}>
        <Story />
      </div>
    ),
  ],
};

const customerSupportRemark = {
  id: '672ff379-f6e3-48b4-a87d-796713f8f997',
  moveID: '9c7b255c-2981-4bf8-839f-61c7458e2b4d',
  officeUserID: 'ce01a5b8-9b44-4511-8a8d-edb60f2a4aee',
  content: 'This is a comment.',
  officeUserFirstName: 'Grace',
  officeUserLastName: 'Griffin',
  createdAt: '2020-06-10T15:58:02.404031Z',
  updatedAt: '2020-06-10T15:58:02.404031Z',
};

export const Default = () => (
  <MockProviders>
    <CustomerSupportRemarkText customerSupportRemark={customerSupportRemark} />
  </MockProviders>
);

export const Edited = () => (
  <MockProviders>
    <CustomerSupportRemarkText
      customerSupportRemark={{ ...customerSupportRemark, updatedAt: '2022-06-15T13:34:34.434543Z' }}
    />
  </MockProviders>
);

export const WithEditPermission = () => (
  <MockProviders currentUserId={customerSupportRemark.officeUserID}>
    <CustomerSupportRemarkText customerSupportRemark={customerSupportRemark} />
  </MockProviders>
);
