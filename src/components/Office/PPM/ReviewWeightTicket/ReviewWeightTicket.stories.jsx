import React from 'react';
import { Grid, GridContainer } from '@trussworks/react-uswds';

import ReviewWeightTicket from './ReviewWeightTicket';

export default {
  title: 'Office Components / PPM / Review Weight Ticket',
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

export const FilledIn = Template.bind({});
FilledIn.args = {
  ppmShipment: {
    actualMoveDate: '2022-04-30',
    actualPickupPostalCode: '90210',
    actualDestinationPostalCode: '94611',
    hasReceivedAdvance: true,
    advanceAmountReceived: 60000,
  },
  tripNumber: 1,
  ppmNumber: 1,
  weightTicket: {
    vehicleDescription: 'Kia Forte',
    emptyWeight: 600,
    fullWeight: 1200,
    ownsTrailer: true,
    trailerMeetsCriteria: false,
  },
};
