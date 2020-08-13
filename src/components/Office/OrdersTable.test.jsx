import React from 'react';
import { shallow } from 'enzyme';

import OrdersTable from './OrdersTable';

const ordersInfo = {
  currentDutyStation: { name: 'JBSA Lackland' },
  newDutyStation: { name: 'JB Lewis-McChord' },
  issuedDate: '8 Mar 2020',
  reportByDate: '1 Apr 2020',
  departmentIndicator: '17 (Navy and Marine Corps)',
  ordersNumber: '999999999',
  ordersType: 'Permanent Change of Duty Station',
  ordersTypeDetail: 'Shipment of HHG permitted',
  tacMDC: '9999',
  sacSDN: '999 999999 999',
};

describe('Orders Table', () => {
  it('should render the data passed to its props', () => {
    const wrapper = shallow(<OrdersTable ordersInfo={ordersInfo} />);
    expect(wrapper.find({ 'data-testid': 'currentDutyStation' }).text()).toMatch('JBSA Lackland');
    expect(wrapper.find({ 'data-testid': 'newDutyStation' }).text()).toMatch('JB Lewis-McChord');
    expect(wrapper.find({ 'data-testid': 'issuedDate' }).text()).toMatch('8 Mar 2020');
    expect(wrapper.find({ 'data-testid': 'reportByDate' }).text()).toMatch('1 Apr 2020');
    expect(wrapper.find({ 'data-testid': 'departmentIndicator' }).text()).toMatch('17 (Navy and Marine Corps)');
    expect(wrapper.find({ 'data-testid': 'ordersNumber' }).text()).toMatch('999999999');
    expect(wrapper.find({ 'data-testid': 'ordersType' }).text()).toMatch('Permanent Change of Duty Station');
    expect(wrapper.find({ 'data-testid': 'ordersTypeDetail' }).text()).toMatch('Shipment of HHG permitted');
    expect(wrapper.find({ 'data-testid': 'tacMDC' }).text()).toMatch('9999');
    expect(wrapper.find({ 'data-testid': 'sacSDN' }).text()).toMatch('999 999999 999');
  });
});
