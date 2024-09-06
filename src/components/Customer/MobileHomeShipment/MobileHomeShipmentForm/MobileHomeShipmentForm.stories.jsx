import React from 'react';
import { expect } from '@storybook/jest';
import { action } from '@storybook/addon-actions';
import { Grid, GridContainer } from '@trussworks/react-uswds';
import { within, userEvent } from '@storybook/testing-library';

import MobileHomeShipmentForm from './MobileHomeShipmentForm';

export default {
  title: 'Customer Components / Mobile Home Shipment / Mobile Home Shipment Form',
  component: MobileHomeShipmentForm,
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

const Template = (args) => <MobileHomeShipmentForm {...args} />;

export const BlankMobileHomeShipmentForm = Template.bind({});
BlankMobileHomeShipmentForm.args = {
  onSubmit: action('submit button clicked'),
  onBack: action('back button clicked'),
  mtoShipment: {
    mobileHomeShipment: {
      year: '',
      make: '',
      model: '',
      lengthInInches: null,
      widthInInches: null,
      heightInInches: null,
    },
  },
};

export const FilledMobileHomeShipmentForm = Template.bind({});
FilledMobileHomeShipmentForm.args = {
  onSubmit: action('submit button clicked'),
  onBack: action('back button clicked'),
  mtoShipment: {
    mobileHomeShipment: {
      year: 2022,
      make: 'Yamaha',
      model: '242X',
      lengthInInches: 288, // 24 feet
      widthInInches: 102, // 8 feet 6 inches
      heightInInches: 84, // 7 feet
    },
  },
};

export const ErrorMobileHomeShipmentForm = Template.bind({});
ErrorMobileHomeShipmentForm.args = {
  onSubmit: action('submit button clicked'),
  onBack: action('back button clicked'),
  mtoShipment: {
    mobileHomeShipment: {
      year: '',
      make: '',
      model: '',
      lengthInInches: null,
      widthInInches: null,
      heightInInches: null,
    },
  },
};
ErrorMobileHomeShipmentForm.play = async ({ canvasElement }) => {
  const canvas = within(canvasElement);

  await expect(canvas.getByRole('button', { name: 'Continue' })).toBeEnabled();

  await userEvent.click(canvas.getByRole('button', { name: 'Continue' }));
};
ErrorMobileHomeShipmentForm.parameters = {
  happo: false,
};
