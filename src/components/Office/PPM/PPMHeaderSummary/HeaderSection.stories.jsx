import React from 'react';
import { Grid, GridContainer } from '@trussworks/react-uswds';

import HeaderSection, { sectionTypes } from './HeaderSection';

export default {
  title: 'Office Components / PPM / Header Section',
  component: HeaderSection,
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

const Template = (args) => <HeaderSection {...args} />;

export const ShipmentSummary = Template.bind({});
ShipmentSummary.args = {
  sectionInfo: {
    type: sectionTypes.shipmentInfo,
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
};

export const Incentives = Template.bind({});
Incentives.args = {
  sectionInfo: {
    type: sectionTypes.incentives,
    estimatedIncentive: 79796,
    gcc: 79796,
    remainingReimbursement: 76173,
  },
  ppmNumber: 1,
};

export const GCC = Template.bind({});
GCC.args = {
  sectionInfo: {
    type: sectionTypes.gcc,
    linehaulPrice: 60803,
    linehaulFuelSurcharge: 0,
    shorthaulPrice: 0,
    shorthaulFuelSurcharge: 0,
    fullPackUnpackCharge: 16427,
  },
  ppmNumber: 1,
};
