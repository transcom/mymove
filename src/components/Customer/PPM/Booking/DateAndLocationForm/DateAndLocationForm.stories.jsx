import React from 'react';
import { action } from '@storybook/addon-actions';
import { Grid, GridContainer } from '@trussworks/react-uswds';
import { within, userEvent } from '@storybook/testing-library';

import DateAndLocationForm from 'components/Customer/PPM/Booking/DateAndLocationForm/DateAndLocationForm';
import { UnsupportedZipCodePPMErrorMsg } from 'utils/validation';

export default {
  title: 'Customer Components / PPM Booking / Date and Location Form',
  component: DateAndLocationForm,
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

const Template = (args) => <DateAndLocationForm {...args} />;

export const BlankDatesAndLocation = Template.bind({});
BlankDatesAndLocation.args = {
  onSubmit: action('submit button clicked'),
  onBack: action('back button clicked'),
  serviceMember: {
    id: '123',
    residentialAddress: {
      postalCode: '90210',
    },
  },
  destinationDutyLocation: {
    address: {
      postalCode: '94611',
    },
  },
  postalCodeValidator: () => {},
};

export const MTOShipmentDatesAndLocation = Template.bind({});
MTOShipmentDatesAndLocation.args = {
  onSubmit: action('submit button clicked'),
  onBack: action('back button clicked'),
  serviceMember: {
    id: '123',
    residential_address: {
      postalCode: '90210',
    },
  },
  destinationDutyLocation: {
    address: {
      postalCode: '94611',
    },
  },
  postalCodeValidator: () => {},
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
    },
  },
};

export const ErrorDatesAndLocation = Template.bind({});
ErrorDatesAndLocation.args = {
  onSubmit: action('submit button clicked'),
  onBack: action('back button clicked'),
  serviceMember: {
    id: '123',
    residential_address: {
      postalCode: '99021',
    },
  },
  destinationDutyLocation: {
    address: {
      postalCode: '94611',
    },
  },
  postalCodeValidator: () => UnsupportedZipCodePPMErrorMsg,
};
ErrorDatesAndLocation.play = async ({ canvasElement }) => {
  // Starts querying the component from its root element
  const canvas = within(canvasElement);

  // See https://storybook.js.org/docs/react/essentials/actions#automatically-matching-args to learn how to setup logging in the Actions panel
  await userEvent.click(await canvas.getByRole('button', { name: 'Save & Continue' }));
};
