import React from 'react';
import { expect } from '@storybook/jest';
import { action } from '@storybook/addon-actions';
import { Grid, GridContainer } from '@trussworks/react-uswds';
import { within, userEvent, waitFor } from '@storybook/testing-library';

import EstimatedWeightsProGearForm from 'components/Customer/PPM/Booking/EstimatedWeightsProGearForm/EstimatedWeightsProGearForm';

export default {
  title: 'Customer Components / PPM Booking / Estimated Weights and Pro-gear',
  component: EstimatedWeightsProGearForm,
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

const Template = (args) => <EstimatedWeightsProGearForm {...args} />;

export const BlankEstimatedWeightsProGear = Template.bind({});
BlankEstimatedWeightsProGear.args = {
  onSubmit: action('submit button clicked'),
  onBack: action('back button clicked'),
  orders: {
    has_dependents: true,
  },
  serviceMember: {
    weight_allotment: {
      total_weight_self_plus_dependents: 8000,
    },
  },
};

export const WarningForOverweightEstimatedWeightProGear = Template.bind({});
WarningForOverweightEstimatedWeightProGear.args = {
  onSubmit: action('submit button clicked'),
  onBack: action('back button clicked'),
  orders: {
    has_dependents: true,
  },
  serviceMember: {
    weight_allotment: {
      total_weight_self_plus_dependents: 5000,
    },
  },
  mtoShipment: {
    id: '123',
    ppmShipment: {
      id: '123',
      estimatedWeight: 7000,
    },
  },
};

export const MTOShipmentEstimatedWeightProGear = Template.bind({});
MTOShipmentEstimatedWeightProGear.args = {
  onSubmit: action('submit button clicked'),
  onBack: action('back button clicked'),
  orders: {
    has_dependents: true,
  },
  serviceMember: {
    weight_allotment: {
      total_weight_self_plus_dependents: 5000,
    },
  },
  mtoShipment: {
    id: '123',
    ppmShipment: {
      id: '123',
      hasProGear: true,
      proGearWeight: 1000,
      spouseProGearWeight: 100,
      estimatedWeight: 4000,
    },
  },
};

export const ErrorEstimatedWeightsProGear = Template.bind({});
ErrorEstimatedWeightsProGear.args = {
  onSubmit: action('submit button clicked'),
  onBack: action('back button clicked'),
  orders: {
    has_dependents: true,
  },
  serviceMember: {
    weight_allotment: {
      total_weight_self_plus_dependents: 5000,
    },
  },
};
ErrorEstimatedWeightsProGear.play = async ({ canvasElement }) => {
  const canvas = within(canvasElement);

  await waitFor(() => {
    expect(canvas.getByRole('button', { name: 'Save & Continue' })).toBeEnabled();
  });
  // See https://storybook.js.org/docs/react/essentials/actions#automatically-matching-args to learn how to setup logging in the Actions panel
  await userEvent.click(canvas.getByRole('button', { name: 'Save & Continue' }));
};
