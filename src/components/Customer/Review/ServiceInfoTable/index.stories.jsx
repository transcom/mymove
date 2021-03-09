/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';

import ServiceInfoTable from '.';

const defaultProps = {
  firstName: 'Jason',
  lastName: 'Ash',
  affiliation: 'Air Force',
  rank: 'E-5',
  edipi: '9999999999',
  currentDutyStationName: 'Buckley AFB',
  currentDutyStationPhone: '555-555-5555',
};

export default {
  title: 'Customer Components / ServiceInfoTable',
};

export const Basic = () => (
  <div style={{ padding: 40 }}>
    <ServiceInfoTable {...defaultProps} />
  </div>
);
