/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { mount } from 'enzyme';

import CustomerHeader from './index';

const props = {
  customer: { last_name: 'Kerry', first_name: 'Smith', dodID: '999999999' },
  order: {
    agency: 'NAVY',
    grade: 'E_6',
    originDutyLocation: {
      name: 'JBSA Lackland',
    },
    originDutyLocationGBLOC: 'AGFM',
    destinationDutyLocation: {
      name: 'JB Lewis-McChord',
    },
    report_by_date: '2018-08-01',
  },
  moveCode: 'FKLCTR',
  move: {
    shipmentGBLOC: 'AGFM',
  },
};

const propsRetiree = {
  customer: { last_name: 'Kerry', first_name: 'Smith', dodID: '999999999' },
  order: {
    agency: 'NAVY',
    grade: 'E_6',
    order_type: 'RETIREMENT',
    originDutyLocation: {
      name: 'JBSA Lackland',
    },
    destinationDutyLocation: {
      name: 'JB Lewis-McChord',
    },
    report_by_date: '2018-08-01',
  },
  moveCode: 'FKLCTR',
  move: {
    shipmentGBLOC: 'AGFM',
  },
};

const propsUSMC = {
  customer: { last_name: 'Kerry', first_name: 'Smith', dodID: '999999999' },
  order: {
    agency: 'MARINES',
    grade: 'E_6',
    originDutyLocation: {
      name: 'JBSA Lackland',
    },
    originDutyLocationGBLOC: 'AGFM',
    destinationDutyLocation: {
      name: 'JB Lewis-McChord',
    },
    report_by_date: '2018-08-01',
  },
  moveCode: 'FKLCTR',
  move: {
    shipmentGBLOC: 'AGFM',
  },
};

const mountCustomerHeader = () => mount(<CustomerHeader {...props} />);
const mountCustomerHeaderRetiree = () => mount(<CustomerHeader {...propsRetiree} />);
const mountCustomerHeaderUSMC = () => mount(<CustomerHeader {...propsUSMC} />);

describe('CustomerHeader component', () => {
  const wrapper = mountCustomerHeader();
  it('renders without crashing', () => {
    expect(wrapper.find('CustomerHeader').length).toBe(1);
  });
  it('renders expected values', () => {
    expect(wrapper.find('[data-testid="nameBlock"]').text()).toContain('Kerry, Smith');
    expect(wrapper.find('[data-testid="nameBlock"]').text()).toContain('FKLCTR');
    expect(wrapper.find('[data-testid="deptPayGrade"]').text()).toContain('Navy E-6');
    expect(wrapper.find('[data-testid="dodId"]').text()).toContain('DoD ID 999999999');
    expect(wrapper.find('[data-testid="infoBlock"]').text()).toContain('JBSA Lackland');
    expect(wrapper.find('[data-testid="infoBlock"]').text()).toContain('JB Lewis-McChord');
    expect(wrapper.find('[data-testid="infoBlock"]').text()).toContain('01 Aug 2018');
    expect(wrapper.find('[data-testid="infoBlock"]').text()).toContain('AGFM');
  });

  const wrapperRetiree = mountCustomerHeaderRetiree();
  it('renders expected values for a retiree', () => {
    expect(wrapperRetiree.find('[data-testid="destinationLabel"]').text()).toContain('HOR, HOS or PLEAD');
    expect(wrapperRetiree.find('[data-testid="reportDateLabel"]').text()).toContain('Date of retirement');
  });

  const wrapperUSMC = mountCustomerHeaderUSMC();
  it('renders expected values for a USMC', () => {
    expect(wrapperUSMC.find('[data-testid="infoBlock"]').text()).toContain('AGFM / USMC');
  });
});
