import React from 'react';
import { Grid, GridContainer } from '@trussworks/react-uswds';

import ReviewWeightTicket from './ReviewWeightTicket';

import { MockProviders } from 'testUtils';
import { createCompleteWeightTicket } from 'utils/test/factories/weightTicket';

export default {
  title: 'Office Components / PPM / Review Weight Ticket',
  component: ReviewWeightTicket,
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

const Template = (args) => <ReviewWeightTicket {...args} />;

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
  weightTicket: {
    vehicleDescription: 'Kia Forte',
    emptyWeight: 600,
    fullWeight: 1200,
    ownsTrailer: true,
    trailerMeetsCriteria: false,
  },
  order: {
    entitlement: {
      authorizedWeight: 2000,
    },
  },
  mtoShipments: [
    { primeActualWeight: 1000, reweigh: null, status: 'APPROVED' },
    { primeActualWeight: 2000, reweigh: { weight: 1000 }, status: 'APPROVED' },
    {
      ppmShipment: {
        weightTickets: [
          createCompleteWeightTicket({}, { fullWeight: 1200, emptyWeight: 200 }),
          createCompleteWeightTicket({}, { fullWeight: 1200, emptyWeight: 200 }),
        ],
      },
      status: 'APPROVED',
    },
  ],
};

export const FilledIn = Template.bind({});
FilledIn.args = {
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
  weightTicket: {
    vehicleDescription: 'Kia Forte',
    emptyWeight: 600,
    fullWeight: 1200,
    ownsTrailer: true,
    trailerMeetsCriteria: false,
  },
  order: {
    entitlement: {
      authorizedWeight: 2000,
    },
  },
  mtoShipments: [
    { primeActualWeight: 1000, reweigh: null, status: 'APPROVED' },
    { primeActualWeight: 2000, reweigh: { weight: 1000 }, status: 'APPROVED' },
    {
      ppmShipment: {
        weightTickets: [
          createCompleteWeightTicket({}, { fullWeight: 1200, emptyWeight: 200 }),
          createCompleteWeightTicket({}, { fullWeight: 1200, emptyWeight: 200 }),
        ],
      },
      status: 'APPROVED',
    },
  ],
};

export const MissingWeightTickets = Template.bind({});
MissingWeightTickets.args = {
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
  weightTicket: {
    vehicleDescription: 'Kia Forte',
    emptyWeight: 6000,
    fullWeight: 8000,
    ownsTrailer: true,
    trailerMeetsCriteria: false,
    missingEmptyWeightTicket: true,
    missingFullWeightTicket: true,
  },
};
