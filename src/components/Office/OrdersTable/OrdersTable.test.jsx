import React from 'react';
import { shallow, mount } from 'enzyme';

import OrdersTable from './OrdersTable';

import { MockProviders } from 'testUtils';

const ordersInfo = {
  currentDutyStation: { name: 'JBSA Lackland' },
  newDutyStation: { name: 'JB Lewis-McChord' },
  issuedDate: '2020-03-08',
  reportByDate: '2020-04-01',
  departmentIndicator: 'NAVY_AND_MARINES',
  ordersNumber: '999999999',
  ordersType: 'PERMANENT_CHANGE_OF_STATION',
  ordersTypeDetail: 'HHG_PERMITTED',
  tacMDC: '9999',
  sacSDN: '999 999999 999',
};

const ordersInfoOptional = {
  currentDutyStation: { name: 'JBSA Lackland' },
  newDutyStation: { name: 'JB Lewis-McChord' },
  issuedDate: '2020-03-08',
  reportByDate: '2020-04-01',
  ordersType: 'SEPARATION',
};

const moveCode = 'TEST';

describe('Orders Table', () => {
  it('should render the data passed to its props', () => {
    const wrapper = shallow(<OrdersTable ordersInfo={ordersInfo} moveCode={moveCode} />);
    expect(wrapper.find({ 'data-testid': 'currentDutyStation' }).text()).toMatch('JBSA Lackland');
    expect(wrapper.find({ 'data-testid': 'newDutyStation' }).text()).toMatch('JB Lewis-McChord');
    expect(wrapper.find({ 'data-testid': 'issuedDate' }).text()).toMatch('08 Mar 2020');
    expect(wrapper.find({ 'data-testid': 'reportByDate' }).text()).toMatch('01 Apr 2020');
    expect(wrapper.find({ 'data-testid': 'departmentIndicator' }).text()).toMatch('17 Navy and Marine Corps');
    expect(wrapper.find({ 'data-testid': 'ordersNumber' }).text()).toMatch('999999999');
    expect(wrapper.find({ 'data-testid': 'ordersType' }).text()).toMatch('Permanent Change Of Station');
    expect(wrapper.find({ 'data-testid': 'ordersTypeDetail' }).text()).toMatch('Shipment of HHG Permitted');
    expect(wrapper.find({ 'data-testid': 'tacMDC' }).text()).toMatch('9999');
    expect(wrapper.find({ 'data-testid': 'sacSDN' }).text()).toMatch('999 999999 999');
  });

  it('should render the table with only required fields', () => {
    const wrapper = shallow(<OrdersTable ordersInfo={ordersInfoOptional} moveCode={moveCode} />);
    expect(wrapper.find({ 'data-testid': 'currentDutyStation' }).text()).toMatch('JBSA Lackland');
    expect(wrapper.find({ 'data-testid': 'newDutyStation' }).text()).toMatch('JB Lewis-McChord');
    expect(wrapper.find({ 'data-testid': 'issuedDate' }).text()).toMatch('08 Mar 2020');
    expect(wrapper.find({ 'data-testid': 'reportByDate' }).text()).toMatch('01 Apr 2020');
    expect(wrapper.find({ 'data-testid': 'departmentIndicator' }).text()).toMatch('');
    expect(wrapper.find({ 'data-testid': 'ordersNumber' }).text()).toMatch('');
    expect(wrapper.find({ 'data-testid': 'ordersType' }).text()).toMatch('Separation');
    expect(wrapper.find({ 'data-testid': 'ordersTypeDetail' }).text()).toMatch('');
    expect(wrapper.find({ 'data-testid': 'tacMDC' }).text()).toMatch('');
    expect(wrapper.find({ 'data-testid': 'sacSDN' }).text()).toMatch('');
  });

  it('should link to the edit orders page with order id param', () => {
    const wrapper = mount(
      <MockProviders initialEntries={[`/moves/1000/details`]}>
        <OrdersTable ordersInfo={ordersInfoOptional} moveCode={moveCode} />
      </MockProviders>,
    );
    expect(wrapper.find('Link').text()).toMatch('View & edit orders');
    expect(wrapper.find('a').prop('href')).toMatch('/moves/TEST/orders');
  });
});
