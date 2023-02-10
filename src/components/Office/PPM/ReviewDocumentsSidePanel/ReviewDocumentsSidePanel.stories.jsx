import React from 'react';
import { Grid, GridContainer } from '@trussworks/react-uswds';

import ReviewDocumentsSidePanel from './ReviewDocumentsSidePanel';

import PPMDocumentsStatus from 'constants/ppms';
import { expenseTypes } from 'constants/ppmExpenseTypes';

export default {
  title: 'Office Components / PPM / Review Documents Side Panel',
  component: ReviewDocumentsSidePanel,
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

const Template = (args) => <ReviewDocumentsSidePanel {...args} />;

export const FilledIn = Template.bind({});
FilledIn.args = {
  ppmShipment: {
    actualMoveDate: '02-Dec-22',
    actualPickupPostalCode: '90210',
    actualDestinationPostalCode: '94611',
    hasReceivedAdvance: true,
    advanceAmountReceived: 60000,
  },
  expenseTickets: [
    { movingExpenseType: expenseTypes.STORAGE, status: PPMDocumentsStatus.REJECTED, reason: 'Too large' },
    { movingExpenseType: expenseTypes.PACKING_MATERIALS, status: PPMDocumentsStatus.ACCEPTED, reason: null },
  ],
  proGearTickets: [
    { status: PPMDocumentsStatus.EXCLUDED, reason: 'Objects not applicable' },
    { status: PPMDocumentsStatus.ACCEPTED, reason: null },
  ],
  weightTickets: [
    {
      status: PPMDocumentsStatus.APPROVED,
    },
  ],
};
