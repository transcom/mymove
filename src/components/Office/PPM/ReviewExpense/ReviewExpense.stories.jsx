import React from 'react';
import { Grid, GridContainer } from '@trussworks/react-uswds';

import ReviewExpense from './ReviewExpense';

import { expenseTypes } from 'constants/ppmExpenseTypes';
import { MockProviders } from 'testUtils';

export default {
  title: 'Office Components / PPM / Review Expense',
  component: ReviewExpense,
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

const Template = (args) => <ReviewExpense {...args} />;

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

export const NonStorage = Template.bind({});
NonStorage.args = {
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
  expense: {
    movingExpenseType: expenseTypes.PACKING_MATERIALS,
    description: 'boxes, tape, bubble wrap',
    amount: 12345,
  },
};

export const Storage = Template.bind({});
Storage.args = {
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
  expense: {
    movingExpenseType: expenseTypes.STORAGE,
    description: 'Pack n store',
    amount: 12345,
    sitStartDate: '2022-12-15',
    sitEndDate: '2022-12-25',
  },
};
