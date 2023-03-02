import React from 'react';
import { Grid, GridContainer } from '@trussworks/react-uswds';

import ReviewShipmentWeightsTable from './ReviewShipmentWeightsTable';

import { PPMReviewWeightsTableColumns } from 'pages/Office/ServicesCounselingReviewShipmentWeights/ServicesCounselingReviewShipmentWeights';

export default {
  title: 'Office Components / PPM / Review Shipment Weights Table',
  component: ReviewShipmentWeightsTable,
  decorators: [
    (Story) => (
      <GridContainer>
        <Grid row>
          <Grid>
            <Story />
          </Grid>
        </Grid>
      </GridContainer>
    ),
  ],
  argTypes: { onClose: { action: 'back button clicked' } },
};

const Template = (args) => <ReviewShipmentWeightsTable {...args} />;

export const PPMShipments = Template.bind({});
PPMShipments.args = {
  tableData: [
    {
      shipmentType: 'PPM',
      ppmShipment: {
        actualMoveDate: '02-Dec-22',
        actualPickupPostalCode: '90210',
        actualDestinationPostalCode: '94611',
        hasReceivedAdvance: true,
        advanceAmountReceived: 60000,
        proGearWeight: 1000,
        spouseProGearWeight: 500,
        estimatedWeight: 4000,
      },
    },
    {
      shipmentType: 'PPM',
      ppmShipment: {
        actualMoveDate: '02-Dec-22',
        actualPickupPostalCode: '90210',
        actualDestinationPostalCode: '94611',
        hasReceivedAdvance: true,
        advanceAmountReceived: 60000,
        proGearWeight: 500,
        spouseProGearWeight: null,
        estimatedWeight: 2000,
      },
    },
  ],
  tableColumns: PPMReviewWeightsTableColumns,
};
