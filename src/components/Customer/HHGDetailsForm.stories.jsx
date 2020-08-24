/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';

import { HHGDetailsFormComponent as HHGDetailsForm } from './HHGDetailsForm';

const defaultProps = {
  pageList: ['page1', 'anotherPage/:foo/:bar'],
  pageKey: 'page1',
  match: { isExact: false, path: '', url: '', params: { moveId: '123' } },
  newDutyStationAddress: {
    city: 'Fort Benning',
    state: 'GA',
    postal_code: '31905',
  },
  currentResidence: {
    city: 'Fort Benning',
    state: 'GA',
    postal_code: '31905',
    street_address_1: '123 Main',
  },
};

export default {
  title: 'Customer Components | HHGDetailsForm',
};

export const Basic = () => <HHGDetailsForm {...defaultProps} />;
