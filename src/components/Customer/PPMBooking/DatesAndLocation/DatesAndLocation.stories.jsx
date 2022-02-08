import React from 'react';
import { action } from '@storybook/addon-actions';
import { within, userEvent } from '@storybook/testing-library';

import DatesAndLocation from './DatesAndLocation';

export default {
  title: 'Customer Components / PPM Booking / Dates and Location',
  component: DatesAndLocation,
};

const Template = (args) => <DatesAndLocation {...args} />;

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
  destinationDutyStation: {
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
    residentialAddress: {
      postalCode: '90210',
    },
  },
  destinationDutyStation: {
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
    residentialAddress: {
      postalCode: '99021',
    },
  },
  destinationDutyStation: {
    address: {
      postalCode: '94611',
    },
  },
  postalCodeValidator: () => {},
};
ErrorDatesAndLocation.play = async ({ canvasElement }) => {
  // Starts querying the component from its root element
  const canvas = within(canvasElement);

  // See https://storybook.js.org/docs/react/essentials/actions#automatically-matching-args to learn how to setup logging in the Actions panel
  await userEvent.click(canvas.getByText('Save & Continue'));
};
