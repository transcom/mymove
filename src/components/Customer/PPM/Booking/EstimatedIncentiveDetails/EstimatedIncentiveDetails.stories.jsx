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
      pickupAddress: {
        streetAddress1: '812 S 129th St',
        streetAddress2: '#123',
        city: 'San Antonio',
        state: 'TX',
        postalCode: '10001',
      },
      destinationAddress: {
        streetAddress1: '813 S 129th St',
        streetAddress2: '#124',
        city: 'San Antonio',
        state: 'TX',
        postalCode: '10002',
      },
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
      pickupAddress: {
        streetAddress1: '812 S 129th St',
        streetAddress2: '#123',
        city: 'San Antonio',
        state: 'TX',
        postalCode: '10001',
      },
      destinationAddress: {
        streetAddress1: '813 S 129th St',
        streetAddress2: '#124',
        city: 'San Antonio',
        state: 'TX',
        postalCode: '10002',
      },
      secondaryPickupAddress: {
        streetAddress1: '814 S 129th St',
        streetAddress2: '#125',
        city: 'San Antonio',
        state: 'TX',
        postalCode: '10001',
      },
      secondaryDestinationAddress: {
        streetAddress1: '815 S 129th St',
        streetAddress2: '#126',
        city: 'San Antonio',
        state: 'TX',
        postalCode: '10002',
      },
      hasSecondaryDestinationAddress: true,
      hasSecondaryPickupAddress: true,
      expectedDepartureDate: '2022-07-04',
      estimatedWeight: 4999,
      estimatedIncentive: 123499,
    },
  },
};
