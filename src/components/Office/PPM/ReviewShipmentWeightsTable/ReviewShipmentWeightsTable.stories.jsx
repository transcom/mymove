import React from 'react';
import { Grid, GridContainer } from '@trussworks/react-uswds';

import { SHIPMENT_OPTIONS } from '../../../../shared/constants';

import ReviewShipmentWeightsTable from './ReviewShipmentWeightsTable';
import { NonPPMTableColumns, NoRowsMessages, PPMReviewWeightsTableColumns, ProGearTableColumns } from './helpers';

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
      showNumber: true,
      shipmentNumber: 1,
      ppmShipment: {
        actualMoveDate: '02-Dec-22',
        actualPickupPostalCode: '90210',
        actualDestinationPostalCode: '94611',
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
      showNumber: true,
      shipmentNumber: 2,
      ppmShipment: {
        actualMoveDate: '02-Dec-22',
        actualPickupPostalCode: '90210',
        actualDestinationPostalCode: '94611',
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
  ],
  tableColumns: PPMReviewWeightsTableColumns,
  noRowsMsg: NoRowsMessages.PPM,
};

export const PPMShipmentsNoRows = Template.bind({});
PPMShipmentsNoRows.args = {
  tableData: [],
  tableColumns: PPMReviewWeightsTableColumns,
  noRowsMsg: NoRowsMessages.PPM,
};

export const ProGearWeights = Template.bind({});
ProGearWeights.args = {
  tableData: [
    {
      entitlement: {
        proGearWeight: 2000,
        spouseProGearWeight: 500,
      },
    },
  ],
  tableColumns: ProGearTableColumns,
};

export const NonPPMShipments = Template.bind({});
NonPPMShipments.args = {
  tableData: [
    {
      shipmentType: SHIPMENT_OPTIONS.HHG,
      primeEstimatedWeight: 2500,
      calculatedBillableWeight: 3000,
      primeActualWeight: 3500,
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
  tableColumns: NonPPMTableColumns,
  noRowsMsg: NoRowsMessages.NonPPM,
};

export const NonPPMShipmentsNoRows = Template.bind({});
NonPPMShipmentsNoRows.args = {
  tableData: [],
  tableColumns: NonPPMTableColumns,
  noRowsMsg: NoRowsMessages.NonPPM,
};
