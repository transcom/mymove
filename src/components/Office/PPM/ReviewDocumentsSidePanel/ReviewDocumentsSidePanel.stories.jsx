import React from 'react';
import { Grid, GridContainer } from '@trussworks/react-uswds';

import ReviewDocumentsSidePanel from './ReviewDocumentsSidePanel';

export default {
  title: 'Office Components / PPM / Review Documents Side Panel',
  component: ReviewDocumentsSidePanel,
  decorators: [
    (Story) => (
      <GridContainer>
        <Grid row>
          <Grid col desktop={{ col: 2, offset: 8 }}>
            <Story />
          </Grid>
        </Grid>
      </GridContainer>
    ),
  ],
  argTypes: { onClose: { action: 'back button clicked' } },
};

const Template = (args) => <ReviewDocumentsSidePanel {...args} />;

export const FilledIn = Template.bind({});
FilledIn.args = {
  ppmShipment: {
    actualMoveDate: '02-Dec-22',
    actualPickupPostalCode: '90210',
    actualDestinationPostalCode: '94611',
    hasReceivedAdvance: true,
    advanceAmountReceived: 60000,
  },
};
