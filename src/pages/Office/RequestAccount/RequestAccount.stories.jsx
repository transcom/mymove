import React from 'react';

import RequestAccount from './RequestAccount';

import { MockProviders } from 'testUtils';

export default {
  title: 'Office Components/RequestAccount',
  parameters: { layout: 'fullscreen' },
};

export const Form = () => {
  return (
    <MockProviders>
      <RequestAccount />
    </MockProviders>
  );
};
