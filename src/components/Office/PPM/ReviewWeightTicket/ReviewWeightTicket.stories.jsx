import React from 'react';
import { Grid, GridContainer } from '@trussworks/react-uswds';

import ReviewWeightTicket from './ReviewWeightTicket';

import { SHIPMENT_OPTIONS } from 'shared/constants';

export default {
  title: 'Office Components / Review Weight Ticket',
  component: ReviewWeightTicket,
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
  argTypes: { onClose: { action: 'back button clicked' } },
};

const Template = (args) => <ReviewWeightTicket {...args} />;

export const Blank = Template.bind({});
Blank.args = {
  ppmShipment: {
    actualMoveDate: Date.now(),
    actualPickupPostalCode: '90210',
    actualDestinationPostalCode: '94611',
    hasReceivedAdvance: true,
    advanceAmountReceived: 60000,
  },
  tripNumber: '1',
  ppmNumber: '1',
};
