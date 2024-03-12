import React from 'react';
import { Grid, GridContainer } from '@trussworks/react-uswds';

import ReviewProGear from './ReviewProGear';

import ppmDocumentStatus from 'constants/ppms';
import { MockProviders } from 'testUtils';

export default {
  title: 'Office Components / PPM / Review ProGear',
  component: ReviewProGear,
  decorators: [
    (Story) => (
      <MockProviders>
        <GridContainer>
          <Grid row>
            <Grid col desktop={{ col: 2, offset: 8 }}>
              <Story />
            </Grid>
          </Grid>
        </GridContainer>
      </MockProviders>
    ),
  ],
  argTypes: { onClose: { action: 'back button clicked' } },
};

const Template = (args) => <ReviewProGear {...args} />;

export const Blank = Template.bind({});
Blank.args = {
  mtoShipment: {
    ppmShipment: {
      actualMoveDate: '2022-04-30',
      actualPickupPostalCode: '90210',
      actualDestinationPostalCode: '94611',
      hasReceivedAdvance: true,
      advanceAmountReceived: 60000,
    },
  },
  tripNumber: 1,
  ppmNumber: 1,
};

export const FilledIn = Template.bind({});
FilledIn.args = {
  mtoShipment: {
    ppmShipment: {
      actualMoveDate: '2023-05-13',
      actualPickupPostalCode: '90210',
      actualDestinationPostalCode: '94611',
      hasReceivedAdvance: true,
      advanceAmountReceived: 60000,
    },
  },
  tripNumber: 1,
  ppmNumber: 1,
  proGear: {
    belongsToSelf: true,
    proGearDocument: [],
    weight: 1000,
    description: 'Description',
    missingWeightTicket: true,
    status: ppmDocumentStatus.REJECTED,
    reason: 'Rejection reason',
  },
};
