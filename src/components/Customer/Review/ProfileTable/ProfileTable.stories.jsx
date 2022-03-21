import React from 'react';

import ProfileTable from './ProfileTable';

export default {
  title: 'Customer Components / ProfileTable',
  component: ProfileTable,
  decorators: [
    (Story) => (
      <div style={{ padding: 40 }}>
        <Story />
      </div>
    ),
  ],
  argTypes: {
    onEditClick: { action: 'profile edit button clicked' },
  },
};

const defaultProps = {
  firstName: 'Jason',
  lastName: 'Ash',
  affiliation: 'Air Force',
  rank: 'E-5',
  edipi: '9999999999',
  currentDutyStationName: 'Buckley AFB',
  telephone: '(999) 999-9999',
  email: 'test@example.com',
  streetAddress1: '17 8th St',
  state: 'New York',
  city: 'NY',
  postalCode: '11111',
};

const ProfileTableTemplate = (args) => <ProfileTable {...args} />;
export const Basic = ProfileTableTemplate.bind({});
Basic.args = defaultProps;
