import React from 'react';
import { expect } from '@storybook/jest';
import { action } from '@storybook/addon-actions';
import { Grid, GridContainer } from '@trussworks/react-uswds';
import { within, userEvent } from '@storybook/testing-library';

import BoatShipmentForm from './BoatShipmentForm';

export default {
  title: 'Customer Components / Boat Shipment / Boat Shipment Form',
  component: BoatShipmentForm,
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

const Template = (args) => <BoatShipmentForm {...args} />;

export const BlankBoatShipmentForm = Template.bind({});
BlankBoatShipmentForm.args = {
  onSubmit: action('submit button clicked'),
  onBack: action('back button clicked'),
  mtoShipment: {
    boatShipment: {
      year: '',
      make: '',
      model: '',
      lengthInInches: null,
      widthInInches: null,
      heightInInches: null,
      hasTrailer: false,
      isRoadworthy: null,
    },
  },
};

export const FilledBoatShipmentForm = Template.bind({});
FilledBoatShipmentForm.args = {
  onSubmit: action('submit button clicked'),
  onBack: action('back button clicked'),
  mtoShipment: {
    boatShipment: {
      year: 2022,
      make: 'Yamaha',
      model: '242X',
      lengthInInches: 288, // 24 feet
      widthInInches: 102, // 8 feet 6 inches
      heightInInches: 84, // 7 feet
      hasTrailer: true,
      isRoadworthy: true,
    },
  },
};

export const ErrorBoatShipmentForm = Template.bind({});
ErrorBoatShipmentForm.args = {
  onSubmit: action('submit button clicked'),
  onBack: action('back button clicked'),
  mtoShipment: {
    boatShipment: {
      year: '',
      make: '',
      model: '',
      lengthInInches: null,
      widthInInches: null,
      heightInInches: null,
      hasTrailer: false,
      isRoadworthy: null,
    },
  },
};
ErrorBoatShipmentForm.play = async ({ canvasElement }) => {
  const canvas = within(canvasElement);

  await expect(canvas.getByRole('button', { name: 'Continue' })).toBeEnabled();

  await userEvent.click(canvas.getByRole('button', { name: 'Continue' }));
};
ErrorBoatShipmentForm.parameters = {
  happo: false,
};
