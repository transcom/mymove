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
  payGrade: 'E-5',
  edipi: '9999999999',
  originDutyLocationName: 'Buckley AFB',
  originTransportationOfficeName: 'Buckley AFB',
  originTransportationOfficePhone: '555-555-5555',
  editURL: '/',
};

const uscgProps = {
  firstName: 'Jason',
  lastName: 'Ash',
  affiliation: 'Coast Guard',
  payGrade: 'E-5',
  edipi: '9999999999',
  emplid: '1234567',
  originDutyLocationName: 'Buckley AFB',
  originTransportationOfficeName: 'Buckley AFB',
  originTransportationOfficePhone: '555-555-5555',
  editURL: '/',
};

export const Editable = () => (
  <div style={{ padding: 40 }}>
    <ServiceInfoDisplay {...defaultProps} isEmplidEnabled />
  </div>
);

export const NonEditableWithMessage = () => (
  <div style={{ padding: 40 }}>
    <ServiceInfoDisplay {...defaultProps} isEditable={false} showMessage isEmplidEnabled />
  </div>
);

export const NonEditableWithoutMessage = () => (
  <div style={{ padding: 40 }}>
    <ServiceInfoDisplay {...defaultProps} isEditable={false} isEmplidEnabled />
  </div>
);

export const CoastGuardCustomer = () => (
  <div style={{ padding: 40 }}>
    <ServiceInfoDisplay {...uscgProps} isEditable={false} isEmplidEnabled />
  </div>
);
