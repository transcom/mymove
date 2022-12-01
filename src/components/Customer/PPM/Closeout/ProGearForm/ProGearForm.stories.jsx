import React from 'react';
import { Grid, GridContainer } from '@trussworks/react-uswds';

import { MockProviders } from 'testUtils';
import ProGearForm from 'components/Customer/PPM/Closeout/ProGearForm/ProGearForm';

export default {
  title: 'Customer Components / PPM Closeout / Pro Gear',
  component: ProGearForm,
  decorators: [
    (Story) => (
      <MockProviders>
        <GridContainer>
          <Grid row>
            <Grid col desktop={{ col: 8, offset: 2 }}>
              <Story />
            </Grid>
          </Grid>
        </GridContainer>
      </MockProviders>
    ),
  ],
  argTypes: {
    onBack: { action: 'back button clicked' },
    onSubmit: { action: 'submit button clicked' },
  },
};

export const Default = (args) => <ProGearForm {...args} />;
