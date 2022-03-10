import React from 'react';
import { action } from '@storybook/addon-actions';
// eslint-disable-next-line import/order
import { Grid, GridContainer } from '@trussworks/react-uswds';
// import { within, userEvent } from '@storybook/testing-library';

import AdvanceForm from 'components/Customer/PPMBooking/Advance/AdvanceForm';

export default {
  title: 'Customer Components / PPM Booking / Advance Form',
  component: AdvanceForm,
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

const Template = (args) => <AdvanceForm {...args} />;

export const BlankAdvanceForm = Template.bind({});
BlankAdvanceForm.args = {
  onSubmit: action('submit button clicked'),
  onBack: action('back button clicked'),
  mtoShipment: {
    id: '123',
    ppmShipment: {
      id: '123',
      pickupPostalCode: '12345',
      secondaryPickupPostalCode: '34512',
      destinationPostalCode: '94611',
      secondaryDestinationPostalCode: '90210',
      sitExpected: 'true',
      expectedDepartureDate: '2022-09-23',
      estimatedIncentive: 10000,
    },
  },
};

// export const MTOShipmentDatesAndLocation = Template.bind({});
// MTOShipmentDatesAndLocation.args = {
//   onSubmit: action('submit button clicked'),
//   onBack: action('back button clicked'),
//   serviceMember: {
//     id: '123',
//     residential_address: {
//       postalCode: '90210',
//     },
//   },
//   destinationDutyStation: {
//     address: {
//       postalCode: '94611',
//     },
//   },
//   postalCodeValidator: () => {},
//   mtoShipment: {
//     id: '123',
//     ppmShipment: {
//       id: '123',
//       pickupPostalCode: '12345',
//       secondaryPickupPostalCode: '34512',
//       destinationPostalCode: '94611',
//       secondaryDestinationPostalCode: '90210',
//       sitExpected: 'true',
//       expectedDepartureDate: '2022-09-23',
//     },
//   },
// };

// export const ErrorDatesAndLocation = Template.bind({});
// ErrorDatesAndLocation.args = {
//   onSubmit: action('submit button clicked'),
//   onBack: action('back button clicked'),
//   serviceMember: {
//     id: '123',
//     residential_address: {
//       postalCode: '99021',
//     },
//   },
//   destinationDutyStation: {
//     address: {
//       postalCode: '94611',
//     },
//   },
//   postalCodeValidator: () =>
//     'Sorry, we don’t support that zip code yet. Please contact your local PPPO for assistance.',
// };
// ErrorDatesAndLocation.play = async ({ canvasElement }) => {
//   // Starts querying the component from its root element
//   const canvas = within(canvasElement);

//   // See https://storybook.js.org/docs/react/essentials/actions#automatically-matching-args to learn how to setup logging in the Actions panel
//   await userEvent.click(canvas.getByText('Save & Continue'));
// };
