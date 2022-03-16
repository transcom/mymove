/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';

import ProfileTable from '.';

const defaultProps = {
  firstName: 'Jason',
  lastName: 'Ash',
  affiliation: 'Air Force',
  rank: 'E-5',
  edipi: '9999999999',
  currentDutyLocationName: 'Buckley AFB',
  telephone: '(999) 999-9999',
  email: 'test@example.com',
  streetAddress1: '17 8th St',
  state: 'New York',
  city: 'NY',
  postalCode: '11111',
};

export default {
  title: 'Customer Components / ProfileTable',
};

export const Basic = () => (
  <div style={{ padding: 40 }}>
    <ProfileTable {...defaultProps} />
  </div>
);
