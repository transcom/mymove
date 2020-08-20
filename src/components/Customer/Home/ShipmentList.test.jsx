/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { mount } from 'enzyme';

import ShipmentList from './ShipmentList';

const defaultProps = {
  shipments: [],
  onShipmentClick: () => {},
};

function mountShipmentList(props = defaultProps) {
  return mount(<ShipmentList {...props} />);
}
describe('ShipmentList component', () => {
  it('renders ShipmentList with shipments', () => {
    const shipments = [
      { id: '#ID-1', type: 'PPM' },
      { id: '#ID-2', type: 'HHG' },
      { id: '#ID-3', type: 'NTS' },
    ];
    const onShipmentClick = () => {};
    const props = {
      shipments,
      onShipmentClick,
    };
    const wrapper = mountShipmentList(props);
    expect(wrapper.find('ShipmentListItem').length).toBe(3);
    expect(wrapper.find('.shipment-list-item-PPM').length).toBe(1);
    expect(wrapper.find('.shipment-list-item-PPM strong').text()).toBe('PPM');
    expect(wrapper.find('.shipment-list-item-PPM span').text()).toBe('#ID-1');
    expect(wrapper.find('.shipment-list-item-HHG').length).toBe(1);
    expect(wrapper.find('.shipment-list-item-HHG strong').text()).toBe('HHG');
    expect(wrapper.find('.shipment-list-item-HHG span').text()).toBe('#ID-2');
    expect(wrapper.find('.shipment-list-item-NTS').length).toBe(1);
    expect(wrapper.find('.shipment-list-item-NTS strong').text()).toBe('NTS');
    expect(wrapper.find('.shipment-list-item-NTS span').text()).toBe('#ID-3');
  });

  it('ShipmentList calls onShipmentClick when clicked', () => {
    const shipments = [
      { id: '#ID-1', type: 'PPM' },
      { id: '#ID-2', type: 'HHG' },
      { id: '#ID-3', type: 'NTS' },
    ];
    const onShipmentClick = jest.fn();
    const props = {
      shipments,
      onShipmentClick,
    };
    const wrapper = mountShipmentList(props);
    expect(onShipmentClick.mock.calls.length).toBe(0);
    wrapper.find('ShipmentListItem').at(0).simulate('click');
    expect(onShipmentClick.mock.calls.length).toBe(1);
    wrapper.find('ShipmentListItem').at(1).simulate('click');
    expect(onShipmentClick.mock.calls.length).toBe(2);
    wrapper.find('ShipmentListItem').at(2).simulate('click');
    expect(onShipmentClick.mock.calls.length).toBe(3);
    const [shipmentOne, shipmentTwo, shipmentThree] = onShipmentClick.mock.calls;
    expect(shipmentOne[0]).toEqual(shipments[0]);
    expect(shipmentTwo[0]).toEqual(shipments[1]);
    expect(shipmentThree[0]).toEqual(shipments[2]);
  });
});
