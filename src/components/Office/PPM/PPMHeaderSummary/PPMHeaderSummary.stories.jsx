import React from 'react';
import { Grid, GridContainer } from '@trussworks/react-uswds';

import PPMHeaderSummary from './PPMHeaderSummary';

export default {
  title: 'Office Components / PPM / PPM Header Summary',
  component: PPMHeaderSummary,
  decorators: [
    (Story) => (
      <GridContainer>
        <Grid row>
          <Grid col desktop={{ col: 4, offset: 4 }}>
            <Story />
          </Grid>
        </Grid>
      </GridContainer>
    ),
  ],
};

const Template = (args) => <PPMHeaderSummary {...args} />;

export const WithAdvanceSingleDocument = Template.bind({});
WithAdvanceSingleDocument.args = {
  ppmShipment: {
    expectedDepartureDate: '2022-04-05',
    actualMoveDate: '2022-04-30',
    actualPickupPostalCode: '90210',
    actualDestinationPostalCode: '94611',
    hasReceivedAdvance: true,
    advanceAmountReceived: 60000,
    miles: 1358,
    estimatedWeight: 5500,
    actualWeight: 6000,
  },
  ppmNumber: 1,
  showAllFields: false,
};

export const WithNoAdvanceSingleDocument = Template.bind({});
WithNoAdvanceSingleDocument.args = {
  ppmShipment: {
    expectedDepartureDate: '2022-04-05',
    actualMoveDate: '2022-04-30',
    actualPickupPostalCode: '90210',
    actualDestinationPostalCode: '94611',
    hasReceivedAdvance: false,
    miles: 1358,
    estimatedWeight: 5500,
    actualWeight: 6000,
  },
  ppmNumber: 1,
  showAllFields: false,
};

export const WithAdvanceReviewAllDocuments = Template.bind({});
WithAdvanceReviewAllDocuments.args = {
  ppmShipment: {
    expectedDepartureDate: '2022-04-05',
    actualMoveDate: '2022-04-30',
    actualPickupPostalCode: '90210',
    actualDestinationPostalCode: '94611',
    hasReceivedAdvance: true,
    advanceAmountReceived: 60000,
    miles: 1358,
    estimatedWeight: 5500,
    actualWeight: 6000,
    incentives: {
      estimatedIncentive: 79796,
      gcc: 79796,
      remainingReimbursement: 76173,
    },
    gcc: {
      linehaulPrice: 60803,
      linehaulFuelSurcharge: 0,
      shorthaulPrice: 0,
      shorthaulFuelSurcharge: 0,
      fullPackUnpackCharge: 16427,
    },
  },
  ppmNumber: 1,
  showAllFields: true,
};

export const WithNoAdvanceReviewAllDocuments = Template.bind({});
WithNoAdvanceReviewAllDocuments.args = {
  ppmShipment: {
    expectedDepartureDate: '2022-04-05',
    actualMoveDate: '2022-04-30',
    actualPickupPostalCode: '90210',
    actualDestinationPostalCode: '94611',
    hasRequestedAdvance: false,
    hasReceivedAdvance: false,
    advanceAmountReceived: 0,
    miles: 1358,
    estimatedWeight: 5500,
    actualWeight: 6000,
    incentives: {
      estimatedIncentive: 79796,
      gcc: 79796,
      remainingReimbursement: 76173,
    },
    gcc: {
      linehaulPrice: 60803,
      linehaulFuelSurcharge: 0,
      shorthaulPrice: 0,
      shorthaulFuelSurcharge: 0,
      fullPackUnpackCharge: 16427,
    },
  },
  ppmNumber: 1,
  showAllFields: true,
};
