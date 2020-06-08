import React from 'react';
import { shallow } from 'enzyme';
import OrdersTable from './OrdersTable';

const ordersInfo = {
  currentDutyStation: 'JBSA Lackland',
  newDutyStation: 'JB Lewis-McChord',
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
    expect(wrapper.find({ 'data-cy': 'currentDutyStation' }).text()).toMatch('JBSA Lackland');
    expect(wrapper.find({ 'data-cy': 'newDutyStation' }).text()).toMatch('JB Lewis-McChord');
    expect(wrapper.find({ 'data-cy': 'issuedDate' }).text()).toMatch('8 Mar 2020');
    expect(wrapper.find({ 'data-cy': 'reportByDate' }).text()).toMatch('1 Apr 2020');
    expect(wrapper.find({ 'data-cy': 'departmentIndicator' }).text()).toMatch('17 (Navy and Marine Corps)');
    expect(wrapper.find({ 'data-cy': 'ordersNumber' }).text()).toMatch('999999999');
    expect(wrapper.find({ 'data-cy': 'ordersType' }).text()).toMatch('Permanent Change of Duty Station');
    expect(wrapper.find({ 'data-cy': 'ordersTypeDetail' }).text()).toMatch('Shipment of HHG permitted');
    expect(wrapper.find({ 'data-cy': 'tacMDC' }).text()).toMatch('9999');
    expect(wrapper.find({ 'data-cy': 'sacSDN' }).text()).toMatch('999 999999 999');
  });
});
