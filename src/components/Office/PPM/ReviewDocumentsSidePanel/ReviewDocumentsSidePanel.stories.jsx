import React from 'react';
import { Grid, GridContainer } from '@trussworks/react-uswds';

import ReviewDocumentsSidePanel from './ReviewDocumentsSidePanel';

import PPMDocumentsStatus from 'constants/ppms';
import { expenseTypes } from 'constants/ppmExpenseTypes';
import { createCompleteMovingExpense } from 'utils/test/factories/movingExpense';
import { createBaseProGearWeightTicket } from 'utils/test/factories/proGearWeightTicket';
import { createCompleteWeightTicket } from 'utils/test/factories/weightTicket';

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
    createCompleteMovingExpense(
      {},
      { movingExpenseType: expenseTypes.STORAGE, status: PPMDocumentsStatus.REJECTED, reason: 'Too large' },
    ),
    createCompleteMovingExpense(
      {},
      { movingExpenseType: expenseTypes.PACKING_MATERIALS, status: PPMDocumentsStatus.APPROVED, reason: null },
    ),
  ],
  proGearTickets: [
    createBaseProGearWeightTicket({}, { status: PPMDocumentsStatus.EXCLUDED, reason: 'Objects not applicable' }),
    createBaseProGearWeightTicket({}, { status: PPMDocumentsStatus.APPROVED, reason: null }),
  ],
  weightTickets: [
    createCompleteWeightTicket(
      {},
      {
        status: PPMDocumentsStatus.APPROVED,
      },
    ),
  ],
};
