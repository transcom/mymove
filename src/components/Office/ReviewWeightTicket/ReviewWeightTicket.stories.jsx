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
  mtoShipment: {
    id: 'f3c29ac7-823a-496a-90dd-b7ab0d4b0ece',
    moveTaskOrderId: 'e9864ee5-56e7-401d-9a7b-a5ea9a83bdea',
    shipmentType: SHIPMENT_OPTIONS.PPM,
    ppmShipment: {
      actualMoveDate: Date.now(),
      actualPickupPostalCode: '90210',
      actualDestinationPostalCode: '94611',
      hasRecievedAdvance: true,
      advanceAmountRequested: 60000,
    },
  },
  tripNumber: '1',
  ppmNumber: '1',
};
