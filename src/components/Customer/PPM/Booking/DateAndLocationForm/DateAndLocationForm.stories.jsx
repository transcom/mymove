import React from 'react';
import { expect } from '@storybook/jest';
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
    residential_address: {
      streetAddress1: '123 Any St',
      streetAddress2: '',
      city: 'Beverly Hills',
      state: 'CA',
      postalCode: '90210',
    },
  },
  destinationDutyLocation: {
    address: {
      streetAddress1: '234 Any Dr',
      streetAddress2: '',
      city: 'Oakland',
      state: 'CA',
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
  const canvas = within(canvasElement);

  await expect(canvas.getByRole('button', { name: 'Save & Continue' })).toBeEnabled();

  await userEvent.click(canvas.getByRole('button', { name: 'Save & Continue' }));
};
ErrorDatesAndLocation.parameters = {
  happo: false,
};
