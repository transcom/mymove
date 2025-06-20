import React from 'react';
import { Grid, GridContainer } from '@trussworks/react-uswds';

import PPMShipmentInfo from '../ppmTestData';

import ReviewGunSafe from './ReviewGunSafe';

import ppmDocumentStatus from 'constants/ppms';
import { MockProviders } from 'testUtils';

export default {
  title: 'Office Components / PPM / Review Gun Safe',
  component: ReviewGunSafe,
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

const Template = (args) => <ReviewGunSafe {...args} />;

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
  gunSafe: {
    gunSafeDocument: [],
    weight: 1000,
    description: 'Description',
    hasWeightTickets: false,
    status: ppmDocumentStatus.REJECTED,
    reason: 'Rejection reason',
  },
};
