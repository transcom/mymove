import React from 'react';
import { Grid, GridContainer } from '@trussworks/react-uswds';

import { SHIPMENT_OPTIONS } from '../../../../shared/constants';

import ReviewShipmentWeightsTable from './ReviewShipmentWeightsTable';
import { PPMReviewWeightsTableConfig, nonPPMReviewWeightsTableConfig } from './helpers';

export default {
  title: 'Office Components / PPM / Review Shipment Weights Table',
  component: ReviewShipmentWeightsTable,
  decorators: [
    (Story) => (
      <GridContainer>
        <Grid row>
          <Grid col desktop={{ col: 'fill' }}>
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
        hasReceivedAdvance: true,
        advanceAmountReceived: 60000,
        proGearWeight: 1000,
        spouseProGearWeight: 500,
        estimatedWeight: 4000,
        expectedDepartureDate: '01-Apr-23',
        weightTickets: [
          {
            emptyWeight: 1000,
            fullWeight: 2000,
          },
        ],
      },
    },
    {
      shipmentType: 'PPM',
      ppmShipment: {
        actualMoveDate: '02-Dec-22',
        hasReceivedAdvance: true,
        advanceAmountReceived: 60000,
        proGearWeight: 500,
        spouseProGearWeight: null,
        estimatedWeight: 2000,
        expectedDepartureDate: '04-Apr-23',
        weightTickets: [
          {
            emptyWeight: 1000,
            fullWeight: 2500,
          },
        ],
      },
    },
    {
      shipmentType: 'PPM',
      ppmShipment: {
        actualMoveDate: '02-Dec-22',
        hasReceivedAdvance: true,
        advanceAmountReceived: 2000,
        proGearWeight: 600,
        spouseProGearWeight: null,
        estimatedWeight: 1000,
        expectedDepartureDate: '04-Apr-23',
        weightTickets: [
          {
            emptyWeight: 1000,
            fullWeight: 2000,
          },
        ],
      },
    },
  ],
  tableConfig: PPMReviewWeightsTableConfig,
};

export const PPMShipmentsNoRows = Template.bind({});
PPMShipmentsNoRows.args = {
  tableData: [],
  tableConfig: PPMReviewWeightsTableConfig,
};

export const NonPPMShipments = Template.bind({});
NonPPMShipments.args = {
  tableData: [
    {
      shipmentType: SHIPMENT_OPTIONS.HHG,
      primeEstimatedWeight: 2500,
      calculatedBillableWeight: 3000,
      primeActualWeight: 3500,
      proGearWeight: 2000,
      spouseProGearWeight: 500,
      reweigh: {
        id: 'rw01',
        weight: 3200,
      },
      actualDeliveryDate: '04-Apr-23',
    },
    {
      shipmentType: SHIPMENT_OPTIONS.NTSR,
      primeActualWeight: 3500,
    },
    {
      shipmentType: SHIPMENT_OPTIONS.NTSR,
      primeActualWeight: 1200,
    },
    {
      shipmentType: SHIPMENT_OPTIONS.NTS,
      ntsRecordedWeight: 1500,
      calculatedBillableWeight: 2000,
      primeActualWeight: 2100,
      reweigh: {
        id: 'rw03',
        weight: 2700,
      },
      actualDeliveryDate: '04-Apr-23',
    },
  ],
  tableConfig: nonPPMReviewWeightsTableConfig,
};

export const NonPPMShipmentsNoRows = Template.bind({});
NonPPMShipmentsNoRows.args = {
  tableData: [],
  tableConfig: nonPPMReviewWeightsTableConfig,
};
