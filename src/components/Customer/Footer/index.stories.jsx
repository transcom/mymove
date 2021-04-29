import React from 'react';

import Footer from './index';

import { MockProviders } from 'testUtils';

export default {
  title: 'Customer Components / Page Footer',
  component: Footer,
  decorators: [
    (Story) => (
      <MockProviders>
        <Story />
      </MockProviders>
    ),
  ],
};

export const Basic = () => <Footer />;
