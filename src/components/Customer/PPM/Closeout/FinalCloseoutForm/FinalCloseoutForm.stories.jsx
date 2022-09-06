import React from 'react';
import { Grid, GridContainer } from '@trussworks/react-uswds';

import FinalCloseoutForm from 'components/Customer/PPM/Closeout/FinalCloseoutForm/FinalCloseoutForm';
import { createPPMShipmentWithFinalIncentive } from 'utils/test/factories/ppmShipment';

export default {
  title: 'Customer Components / PPM Closeout / Final Closeout Form',
  component: FinalCloseoutForm,
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
  argTypes: { onBack: { action: 'back button clicked' }, onSubmit: { action: 'submit button clicked' } },
};

const Template = (args) => <FinalCloseoutForm {...args} />;

export const Blank = Template.bind({});
Blank.args = {
  mtoShipment: createPPMShipmentWithFinalIncentive(),
};
