/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';

import ServiceInfoDisplay from './ServiceInfoDisplay';

const defaultProps = {
  firstName: 'Jason',
  lastName: 'Ash',
  affiliation: 'Air Force',
  rank: 'E-5',
  edipi: '9999999999',
  originDutyStationName: 'Buckley AFB',
  originTransportationOfficeName: 'Buckley AFB',
  originTransportationOfficePhone: '555-555-5555',
  editURL: '/',
};

export default {
  title: 'Customer Components / ServiceInfoDisplay',
};

export const Editable = () => (
  <div style={{ padding: 40 }}>
    <ServiceInfoDisplay {...defaultProps} />
  </div>
);

export const NonEditable = () => (
  <div style={{ padding: 40 }}>
    <ServiceInfoDisplay {...defaultProps} isEditable={false} />
  </div>
);
