import React from 'react';
import { Grid, GridContainer } from '@trussworks/react-uswds';

import ReviewProGear from './ReviewProGear';

import ppmDocumentStatus from 'constants/ppms';

export default {
  title: 'Office Components / PPM / Review ProGear',
  component: ReviewProGear,
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

const Template = (args) => <ReviewProGear {...args} />;

export const Blank = Template.bind({});
Blank.args = {
  ppmShipment: {
    actualMoveDate: '02-Dec-22',
    actualPickupPostalCode: '90210',
    actualDestinationPostalCode: '94611',
    hasReceivedAdvance: true,
    advanceAmountReceived: 60000,
  },
  tripNumber: 1,
  ppmNumber: 1,
};

export const FilledIn = Template.bind({});
FilledIn.args = {
  ppmShipment: {
    actualMoveDate: '02-Dec-22',
    actualPickupPostalCode: '90210',
    actualDestinationPostalCode: '94611',
    hasReceivedAdvance: true,
    advanceAmountReceived: 60000,
  },
  tripNumber: 1,
  ppmNumber: 1,
  proGear: {
    selfProGear: true,
    proGearDocument: [],
    proGearWeight: 1,
    description: 'Description',
    missingWeightTicket: true,
    status: ppmDocumentStatus.REJECTED,
    reason: 'Rejection reason',
  },
};
