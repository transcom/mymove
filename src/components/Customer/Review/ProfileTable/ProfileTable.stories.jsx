import React from 'react';
import { Grid, GridContainer } from '@trussworks/react-uswds';

import ProfileTable from './ProfileTable';

export default {
  title: 'Customer Components / ProfileTable',
  component: ProfileTable,
  decorators: [
    (Story) => (
      <GridContainer>
        <Grid row>
          <Grid col desktop={{ col: 8, offset: 2 }}>
            <Story />
          </Grid>
        </Grid>
      </GridContainer>
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
  edipi: '9999999999',
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
