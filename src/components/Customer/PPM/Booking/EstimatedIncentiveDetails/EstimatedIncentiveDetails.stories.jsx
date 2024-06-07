import React from 'react';
import { Grid, GridContainer } from '@trussworks/react-uswds';

import EstimatedIncentiveDetails from 'components/Customer/PPM/Booking/EstimatedIncentiveDetails/EstimatedIncentiveDetails';

export default {
  title: 'Customer Components / PPM Booking / Estimated Incentive Details',
  component: EstimatedIncentiveDetails,
  decorators: [
    (Story) => (
      <GridContainer>
        <Grid row>
          <Grid col desktop={{ col: 8, offset: 2 }}>
            <Story />
          </Grid>
        </Grid>
      </GridContainer>
    ),
  ],
};

const Template = (args) => <EstimatedIncentiveDetails {...args} />;

export const WithoutSecondaryPostalCodes = Template.bind({});
WithoutSecondaryPostalCodes.args = {
  shipment: {
    ppmShipment: {
      expectedDepartureDate: '2022-07-04',
      estimatedWeight: 4999,
      estimatedIncentive: 123499,
    },
  },
};

export const WithSecondaryPostalCodes = Template.bind({});
WithSecondaryPostalCodes.args = {
  shipment: {
    ppmShipment: {
      expectedDepartureDate: '2022-07-04',
      estimatedWeight: 4999,
      estimatedIncentive: 123499,
    },
  },
};
