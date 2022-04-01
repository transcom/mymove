/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';

import ServiceInfoDisplay from './ServiceInfoDisplay';

import { MockProviders } from 'testUtils';

export default {
  title: 'Customer Components / ServiceInfoDisplay',
  component: ServiceInfoDisplay,
  decorators: [
    (Story) => (
      <MockProviders>
        <Story />
      </MockProviders>
    ),
  ],
};

const defaultProps = {
  firstName: 'Jason',
  lastName: 'Ash',
  affiliation: 'Air Force',
  rank: 'E-5',
  edipi: '9999999999',
  originDutyLocationName: 'Buckley AFB',
  originTransportationOfficeName: 'Buckley AFB',
  originTransportationOfficePhone: '555-555-5555',
  editURL: '/',
};

export const Editable = () => (
  <div style={{ padding: 40 }}>
    <ServiceInfoDisplay {...defaultProps} />
  </div>
);

export const NonEditableWithMessage = () => (
  <div style={{ padding: 40 }}>
    <ServiceInfoDisplay {...defaultProps} isEditable={false} showMessage />
  </div>
);

export const NonEditableWithoutMessage = () => (
  <div style={{ padding: 40 }}>
    <ServiceInfoDisplay {...defaultProps} isEditable={false} />
  </div>
);
