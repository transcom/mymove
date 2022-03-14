import React from 'react';
import { action } from '@storybook/addon-actions';
// eslint-disable-next-line import/order
import { Grid, GridContainer } from '@trussworks/react-uswds';
// import { within, userEvent } from '@storybook/testing-library';

import AdvanceForm from 'components/Customer/PPMBooking/Advance/AdvanceForm';

export default {
  title: 'Customer Components / PPM Booking / Advance Form',
  component: AdvanceForm,
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
};

const Template = (args) => <AdvanceForm {...args} />;

export const BlankAdvanceForm = Template.bind({});
BlankAdvanceForm.args = {
  onSubmit: action('submit button clicked'),
  onBack: action('back button clicked'),
  mtoShipment: {
    id: '123',
    ppmShipment: {
      id: '123',
      pickupPostalCode: '12345',
      secondaryPickupPostalCode: '34512',
      destinationPostalCode: '94611',
      secondaryDestinationPostalCode: '90210',
      sitExpected: 'true',
      expectedDepartureDate: '2022-09-23',
      estimatedIncentive: 1000000,
    },
  },
};

export const PreFilledAdvanceForm = Template.bind({});
PreFilledAdvanceForm.args = {
  onSubmit: action('submit button clicked'),
  onBack: action('back button clicked'),
  mtoShipment: {
    id: '123',
    ppmShipment: {
      id: '123',
      pickupPostalCode: '12345',
      secondaryPickupPostalCode: '34512',
      destinationPostalCode: '94611',
      secondaryDestinationPostalCode: '90210',
      expectedDepartureDate: '2022-09-23',
      estimatedIncentive: 1000000,
      advance: 30000,
    },
  },
};

export const MaxRequestedExceededAdvanceForm = Template.bind({});
MaxRequestedExceededAdvanceForm.args = {
  onSubmit: action('submit button clicked'),
  onBack: action('back button clicked'),
  mtoShipment: {
    id: '123',
    ppmShipment: {
      id: '123',
      pickupPostalCode: '12345',
      secondaryPickupPostalCode: '34512',
      destinationPostalCode: '94611',
      secondaryDestinationPostalCode: '90210',
      expectedDepartureDate: '2022-09-23',
      estimatedIncentive: 1000000,
      advance: 300000000,
    },
  },
};
