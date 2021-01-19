/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { mount } from 'enzyme';

import CustomerHeader from './index';

const props = {
  customer: { last_name: 'Kerry', first_name: 'Smith', dodID: '999999999' },
  moveOrder: {
    departmentIndicator: 'Navy',
    grade: 'E_6',
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
  it('renders expected values', () => {
    expect(wrapper.find('[data-test="nameBlock"]').text()).toContain('Kerry, Smith');
    expect(wrapper.find('[data-test="nameBlock"]').text()).toContain('FKLCTR');
    expect(wrapper.find('[data-test="deptRank"]').text()).toContain('Navy E-6');
    expect(wrapper.find('[data-test="dodId"]').text()).toContain('DoD ID 999999999');
    expect(wrapper.find('[data-test="infoBlock"]').text()).toContain('JBSA Lackland');
    expect(wrapper.find('[data-test="infoBlock"]').text()).toContain('JB Lewis-McChord');
    expect(wrapper.find('[data-test="infoBlock"]').text()).toContain('01 Aug 2018');
  });
});
