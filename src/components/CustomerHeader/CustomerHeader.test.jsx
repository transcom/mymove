/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { mount } from 'enzyme';

import CustomerHeader from './index';

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

const mountCustomerHeader = () => mount(<CustomerHeader {...props} />);

describe('CustomerHeader component', () => {
  const wrapper = mountCustomerHeader();
  it('renders without crashing', () => {
    expect(wrapper.find('CustomerHeader').length).toBe(1);
  });
});
