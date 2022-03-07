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
  order: {
    agency: 'MARINES',
    grade: 'E_6',
    originDutyLocation: {
      name: 'JBSA Lackland',
    },
    destinationDutyLocation: {
      name: 'JB Lewis-McChord',
    },
    report_by_date: '2018-08-01',
  },
  moveCode: 'FKLCTR',
};

const propsRetirement = {
  ...props,
  order: {
    ...props.order,
    order_type: 'RETIREMENT',
  },
};

const propsSeparation = {
  ...props,
  order: {
    ...props.order,
    order_type: 'SEPARATION',
  },
};

// eslint-disable-next-line react/jsx-props-no-spreading
export const Customer = () => (
  <div style={{ minWidth: '1000px' }}>
    <CustomerHeader {...props} />
  </div>
);

export const CustomerRetirement = () => (
  <div style={{ minWidth: '1000px' }}>
    <CustomerHeader {...propsRetirement} />
  </div>
);

export const CustomerSeparation = () => (
  <div style={{ minWidth: '1000px' }}>
    <CustomerHeader {...propsSeparation} />
  </div>
);
