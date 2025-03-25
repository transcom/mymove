import React from 'react';
import { Grid, GridContainer } from '@trussworks/react-uswds';

import ProGearForm from 'components/Shared/PPM/Closeout/ProGearForm/ProGearForm';
import { MockProviders } from 'testUtils';

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
