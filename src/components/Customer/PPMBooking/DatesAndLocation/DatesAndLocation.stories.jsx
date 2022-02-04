import React from 'react';
import { action } from '@storybook/addon-actions';

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
