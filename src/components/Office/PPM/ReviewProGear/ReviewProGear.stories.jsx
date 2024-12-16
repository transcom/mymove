import React from 'react';
import { Grid, GridContainer } from '@trussworks/react-uswds';

import PPMShipmentInfo from '../ppmTestData';

import ReviewProGear from './ReviewProGear';

import ppmDocumentStatus from 'constants/ppms';
import { MockProviders } from 'testUtils';

export default {
  title: 'Office Components / PPM / Review ProGear',
  component: ReviewProGear,
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

const Template = (args) => <ReviewProGear {...args} />;

export const Blank = Template.bind({});
Blank.args = {
  ppmShipmentInfo: PPMShipmentInfo,
  tripNumber: 1,
  ppmNumber: '1',
};

export const FilledIn = Template.bind({});
FilledIn.args = {
  ppmShipmentInfo: PPMShipmentInfo,
  tripNumber: 1,
  ppmNumber: '1',
  proGear: {
    belongsToSelf: true,
    proGearDocument: [],
    weight: 1000,
    description: 'Description',
    missingWeightTicket: true,
    status: ppmDocumentStatus.REJECTED,
    reason: 'Rejection reason',
  },
};
