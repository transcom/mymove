import React from 'react';
import { Grid, GridContainer } from '@trussworks/react-uswds';

import PPMShipmentInfo from '../ppmTestData';

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
  ppmShipmentInfo: PPMShipmentInfo,
  tripNumber: 1,
  ppmNumber: '1',
  weightTicket: {
    vehicleDescription: 'Kia Forte',
    emptyWeight: 600,
    fullWeight: 1200,
    ownsTrailer: true,
    trailerMeetsCriteria: false,
  },
  order: {
    entitlement: {
      totalWeight: 2000,
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
  ppmShipmentInfo: PPMShipmentInfo,
  tripNumber: 1,
  ppmNumber: '1',
  weightTicket: {
    vehicleDescription: 'Kia Forte',
    emptyWeight: 600,
    fullWeight: 1200,
    ownsTrailer: true,
    trailerMeetsCriteria: false,
  },
  order: {
    entitlement: {
      totalWeight: 2000,
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
  ppmShipmentInfo: PPMShipmentInfo,
  tripNumber: 1,
  ppmNumber: '1',
  weightTicket: {
    vehicleDescription: 'Kia Forte',
    emptyWeight: 6000,
    fullWeight: 8000,
    ownsTrailer: true,
    trailerMeetsCriteria: false,
    missingEmptyWeightTicket: true,
    missingFullWeightTicket: true,
  },
  order: {
    entitlement: {
      totalWeight: 2000,
    },
  },
};
