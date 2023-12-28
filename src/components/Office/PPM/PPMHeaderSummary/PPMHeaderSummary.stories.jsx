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
      baseLinehaul: 60803,
      originLinehaulFactor: 109,
      destinationLinehaulFactor: 68,
      linehaulAdjustment: 433,
      shorthaulCharge: 0,
      transportationCost: 61276,
      linehaulFuelSurcharge: 0,
      fuelSurchargePercent: 0,
      originServiceAreaFee: 512,
      originFactor: 1253,
      destinationServiceAreaFee: 343,
      destinationFactor: 839,
      fullPackUnpackCharge: 16427,
      ppmFactor: 0,
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
      baseLinehaul: 60803,
      originLinehaulFactor: 109,
      destinationLinehaulFactor: 68,
      linehaulAdjustment: 433,
      shorthaulCharge: 0,
      transportationCost: 61276,
      linehaulFuelSurcharge: 0,
      fuelSurchargePercent: 0,
      originServiceAreaFee: 512,
      originFactor: 1253,
      destinationServiceAreaFee: 343,
      destinationFactor: 839,
      fullPackUnpackCharge: 16427,
      ppmFactor: 0,
    },
  },
  ppmNumber: 1,
  showAllFields: true,
};
