import React from 'react';
import { Grid, GridContainer } from '@trussworks/react-uswds';

import PPMShipmentInfo from '../ppmTestData';

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
  ppmShipmentInfo: PPMShipmentInfo,
  tripNumber: 1,
  ppmNumber: 1,
};

export const NonStorage = Template.bind({});
NonStorage.args = {
  ppmShipmentInfo: PPMShipmentInfo,
  tripNumber: 1,
  ppmNumber: 1,
  categoryIndex: 1,
  expense: {
    movingExpenseType: expenseTypes.PACKING_MATERIALS,
    description: 'boxes, tape, bubble wrap',
    amount: 12345,
  },
};

export const Storage = Template.bind({});
Storage.args = {
  ppmShipmentInfo: PPMShipmentInfo,
  tripNumber: 1,
  ppmNumber: 1,
  categoryIndex: 1,
  expense: {
    movingExpenseType: expenseTypes.STORAGE,
    description: 'Pack n store',
    amount: 12345,
    sitStartDate: '2022-12-15',
    sitEndDate: '2022-12-25',
  },
};
