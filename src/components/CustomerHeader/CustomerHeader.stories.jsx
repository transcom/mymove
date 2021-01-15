import React from 'react';
import { Provider } from 'react-redux';

import CustomerHeader from './index';

import { store } from 'shared/store';

export default {
  title: 'Components/Headers/Customer Header',
  parameters: {
    abstract: {
      url: 'https://share.goabstract.com/d9ad20e6-944c-48a2-bbd2-1c7ed8bc1315?mode=design',
    },
  },
};

const props = {
  customer: { last_name: 'Kerry', first_name: 'Smith', dodID: '999999999' },
  moveOrder: {
    departmentIndicator: 'Navy',
    grade: 'E-6',
    originDutyStation: {
      name: 'JBSA Lackland',
    },
    destinationDutyStation: {
      name: 'JB Lewis-McChord',
    },
  },
  moveCode: 'FKLCTR',
};

/* eslint-disable react/jsx-props-no-spreading */
export const Customer = () => {
  return (
    <Provider store={store}>
      <CustomerHeader {...props} />
    </Provider>
  );
};
