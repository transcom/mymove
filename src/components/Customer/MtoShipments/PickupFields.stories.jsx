/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';

import { PickupFields } from './PickupFields';

const defaultProps = {
  currentResidence: {
    city: 'Fort Benning',
    state: 'GA',
    postal_code: '31905',
    street_address_1: '123 Main',
  },
};

export default {
  title: 'Customer Components | PickupFields',
};

export const Basic = () => <PickupFields {...defaultProps} />;
