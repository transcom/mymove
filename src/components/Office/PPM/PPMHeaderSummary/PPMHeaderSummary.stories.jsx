import React from 'react';
import { Grid, GridContainer } from '@trussworks/react-uswds';

import PPMHeaderSummary from './PPMHeaderSummary';

export default {
  title: 'Office Components / PPM / PPM Header Summary',
  component: PPMHeaderSummary,
  decorators: [
    (Story) => (
      <GridContainer>
        <Grid row>
          <Grid col desktop={{ col: 2, offset: 8 }}>
            <Story />
          </Grid>
        </Grid>
      </GridContainer>
    ),
  ],
};

const Template = (args) => <PPMHeaderSummary {...args} />;

export const WithAdvance = Template.bind({});
WithAdvance.args = {
  ppmShipment: {
    actualMoveDate: Date.now(),
    actualPickupPostalCode: '90210',
    actualDestinationPostalCode: '94611',
    hasReceivedAdvance: true,
    advanceAmountReceived: 60000,
  },
  ppmNumber: '1',
};

export const WithNoAdvance = Template.bind({});
WithNoAdvance.args = {
  ppmShipment: {
    actualMoveDate: Date.now(),
    actualPickupPostalCode: '90210',
    actualDestinationPostalCode: '94611',
    hasReceivedAdvance: false,
  },
  ppmNumber: '1',
};
