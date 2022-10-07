import React from 'react';
import { Grid, GridContainer } from '@trussworks/react-uswds';
import { userEvent } from '@storybook/testing-library';

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

export const W2AddressFormFilledOut = Template.bind({});
W2AddressFormFilledOut.storyName = 'W2 Address Form Filled Out';
W2AddressFormFilledOut.args = {
  mtoShipment: {
    ppmShipment: {
      w2Address: {
        streetAddress1: '123 Anywhere St',
        streetAddress2: '',
        city: 'Santa Monica',
        state: 'CA',
        postalCode: '90402',
      },
    },
  },
};

export const W2AddressFormWithErrors = Template.bind({});
W2AddressFormWithErrors.storyName = 'W2 Address Form With Errors';
W2AddressFormWithErrors.args = {
  mtoShipment: {
    ppmShipment: {
      w2Address: {
        streetAddress1: '',
        streetAddress2: '',
        city: '',
        state: '',
        postalCode: '90',
      },
    },
  },
};
W2AddressFormWithErrors.play = async () => {
  await userEvent.tab();
  await userEvent.tab();
  await userEvent.tab();
  await userEvent.tab();
  await userEvent.tab();
  await userEvent.tab();
};

W2AddressFormWithErrors.parameters = {
  happo: false,
};
