/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';

import { DeliveryFields } from './DeliveryFields';

const defaultProps = {
  newDutyStationAddress: {
    city: 'Fort Benning',
    state: 'GA',
    postal_code: '31905',
  },
};

export default {
  title: 'Customer Components | DeliveryFields',
};

export const Basic = () => <DeliveryFields {...defaultProps} />;
