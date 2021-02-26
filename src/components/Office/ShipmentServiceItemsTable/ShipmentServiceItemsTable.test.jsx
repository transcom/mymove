import React from 'react';
import { shallow } from 'enzyme';

import ShipmentServiceItemsTable from './ShipmentServiceItemsTable';

import { SHIPMENT_OPTIONS } from 'shared/constants';

describe('Shipment Service Items Table', () => {
  it('renders the hhg longhaul shipment type with service items', () => {
    const wrapper = shallow(<ShipmentServiceItemsTable shipmentType={SHIPMENT_OPTIONS.HHG_LONGHAUL_DOMESTIC} />);
    expect(wrapper.find('table.serviceItemsTable').exists()).toBe(true);
    expect(wrapper.find('h4').text()).toEqual('Service items for this shipment 6 items');
    expect(wrapper.find('th').text()).toEqual('Service item');
    expect(wrapper.find('td').length).toBe(6);
    expect(wrapper.find('td').at(0).text()).toEqual('Domestic linehaul');
    expect(wrapper.find('td').at(1).text()).toEqual('Fuel surcharge');
    expect(wrapper.find('td').at(2).text()).toEqual('Domestic origin price');
    expect(wrapper.find('td').at(3).text()).toEqual('Domestic destination price');
    expect(wrapper.find('td').at(4).text()).toEqual('Domestic packing');
    expect(wrapper.find('td').at(5).text()).toEqual('Domestic unpacking');
  });

  it('renders the hhg shorthaul shipment type with service items', () => {
    const wrapper = shallow(<ShipmentServiceItemsTable shipmentType={SHIPMENT_OPTIONS.HHG_SHORTHAUL_DOMESTIC} />);
    expect(wrapper.find('table.serviceItemsTable').exists()).toBe(true);
    expect(wrapper.find('h4').text()).toEqual('Service items for this shipment 6 items');
    expect(wrapper.find('th').text()).toEqual('Service item');
    expect(wrapper.find('td').length).toBe(6);
    expect(wrapper.find('td').at(0).text()).toEqual('Domestic shorthaul');
    expect(wrapper.find('td').at(1).text()).toEqual('Fuel surcharge');
    expect(wrapper.find('td').at(2).text()).toEqual('Domestic origin price');
    expect(wrapper.find('td').at(3).text()).toEqual('Domestic destination price');
    expect(wrapper.find('td').at(4).text()).toEqual('Domestic packing');
    expect(wrapper.find('td').at(5).text()).toEqual('Domestic unpacking');
  });

  it('renders the nts shipment type with service items', () => {
    const wrapper = shallow(<ShipmentServiceItemsTable shipmentType={SHIPMENT_OPTIONS.NTS} />);
    expect(wrapper.find('table.serviceItemsTable').exists()).toBe(true);
    expect(wrapper.find('h4').text()).toEqual('Service items for this shipment 6 items');
    expect(wrapper.find('th').text()).toEqual('Service item');
    expect(wrapper.find('td').length).toBe(6);
    expect(wrapper.find('td').at(0).text()).toEqual('Domestic linehaul');
    expect(wrapper.find('td').at(1).text()).toEqual('Fuel surcharge');
    expect(wrapper.find('td').at(2).text()).toEqual('Domestic origin price');
    expect(wrapper.find('td').at(3).text()).toEqual('Domestic destination price');
    expect(wrapper.find('td').at(4).text()).toEqual('Domestic packing');
    expect(wrapper.find('td').at(5).text()).toEqual('Domestic NTS packing factor');
  });

  it('renders the nts release shipment type with service items', () => {
    const wrapper = shallow(<ShipmentServiceItemsTable shipmentType={SHIPMENT_OPTIONS.NTSR} />);
    expect(wrapper.find('table.serviceItemsTable').exists()).toBe(true);
    expect(wrapper.find('h4').text()).toEqual('Service items for this shipment 5 items');
    expect(wrapper.find('th').text()).toEqual('Service item');
    expect(wrapper.find('td').length).toBe(5);
    expect(wrapper.find('td').at(0).text()).toEqual('Domestic linehaul');
    expect(wrapper.find('td').at(1).text()).toEqual('Fuel surcharge');
    expect(wrapper.find('td').at(2).text()).toEqual('Domestic origin price');
    expect(wrapper.find('td').at(3).text()).toEqual('Domestic destination price');
    expect(wrapper.find('td').at(4).text()).toEqual('Domestic unpacking');
  });
});
