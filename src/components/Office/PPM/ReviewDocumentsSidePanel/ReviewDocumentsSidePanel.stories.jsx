import React from 'react';
import { Grid, GridContainer } from '@trussworks/react-uswds';

import PPMShipmentInfo from '../ppmTestData';

import ReviewDocumentsSidePanel from './ReviewDocumentsSidePanel';

import PPMDocumentsStatus from 'constants/ppms';
import { expenseTypes } from 'constants/ppmExpenseTypes';
import { createCompleteMovingExpense } from 'utils/test/factories/movingExpense';
import { createBaseProGearWeightTicket } from 'utils/test/factories/proGearWeightTicket';
import { createBaseGunSafeWeightTicket } from 'utils/test/factories/gunSafeWeightTicket';
import { createCompleteWeightTicket } from 'utils/test/factories/weightTicket';
import { MockProviders } from 'testUtils';

export default {
  title: 'Office Components / PPM / Review Documents Side Panel',
  component: ReviewDocumentsSidePanel,
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

const Template = (args) => <ReviewDocumentsSidePanel {...args} />;

export const FilledIn = Template.bind({});
FilledIn.args = {
  ppmShipmentInfo: PPMShipmentInfo,
  expenseTickets: [
    createCompleteMovingExpense(
      {},
      {
        movingExpenseType: expenseTypes.STORAGE,
        status: PPMDocumentsStatus.REJECTED,
        sitStartDate: '2023-02-01',
        sitEndDate: '2023-03-01',
        amount: 30000,
        reason: 'Too large',
      },
    ),
    createCompleteMovingExpense(
      {},
      {
        movingExpenseType: expenseTypes.PACKING_MATERIALS,
        status: PPMDocumentsStatus.APPROVED,
        amount: 20000,
        reason: null,
      },
    ),
  ],
  proGearTickets: [
    createBaseProGearWeightTicket({}, { status: PPMDocumentsStatus.EXCLUDED, reason: 'Objects not applicable' }),
    createBaseProGearWeightTicket({}, { status: PPMDocumentsStatus.APPROVED, reason: null }),
  ],
  gunSafeTickets: [
    createBaseGunSafeWeightTicket({}, { status: PPMDocumentsStatus.REJECTED, reason: 'Objects not applicable' }),
    createBaseGunSafeWeightTicket({}, { status: PPMDocumentsStatus.APPROVED, reason: null }),
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
