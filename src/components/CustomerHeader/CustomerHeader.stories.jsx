import React from 'react';

import CustomerHeader from './index';

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
    report_by_date: '2018-08-01',
  },
  moveCode: 'FKLCTR',
};

// eslint-disable-next-line react/jsx-props-no-spreading
export const Customer = () => <CustomerHeader {...props} />;
