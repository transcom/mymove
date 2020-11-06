/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { mount } from 'enzyme';

import ShipmentList from '.';

const shipments = [
  { id: 'ID-1', shipmentType: 'PPM' },
  { id: 'ID-2', shipmentType: 'HHG' },
  { id: 'ID-3', shipmentType: 'HHG_INTO_NTS_DOMESTIC' },
  { id: 'ID-4', shipmentType: 'HHG_OUTOF_NTS_DOMESTIC' },
];
const onShipmentClick = jest.fn();
const defaultProps = {
  shipments,
  onShipmentClick,
  moveSubmitted: false,
};

describe('ShipmentList component', () => {
  it('renders ShipmentList with shipments', () => {
    const wrapper = mount(<ShipmentList {...defaultProps} />);
    expect(wrapper.find('ShipmentListItem').length).toBe(4);
    expect(wrapper.find('.shipment-list-item-PPM').length).toBe(1);
    expect(wrapper.find('.shipment-list-item-PPM strong').text()).toBe('PPM');
    expect(wrapper.find('.shipment-list-item-PPM span').text()).toBe('#ID-1');
    expect(wrapper.find('.shipment-list-item-HHG').length).toBe(1);
    expect(wrapper.find('.shipment-list-item-HHG strong').text()).toBe('HHG');
    expect(wrapper.find('.shipment-list-item-HHG span').text()).toBe('#ID-2');
    expect(wrapper.find('.shipment-list-item-NTS').length).toBe(1);
    expect(wrapper.find('.shipment-list-item-NTS strong').text()).toBe('NTS');
    expect(wrapper.find('.shipment-list-item-NTS span').text()).toBe('#ID-3');
    expect(wrapper.find('.shipment-list-item-NTS-R').length).toBe(1);
    expect(wrapper.find('.shipment-list-item-NTS-R strong').text()).toBe('NTS-R');
    expect(wrapper.find('.shipment-list-item-NTS-R span').text()).toBe('#ID-4');
  });

  it('ShipmentList calls onShipmentClick when clicked', () => {
    const wrapper = mount(<ShipmentList {...defaultProps} />);
    expect(onShipmentClick.mock.calls.length).toBe(0);
    wrapper.find('ShipmentListItem').at(0).simulate('click');
    expect(onShipmentClick.mock.calls.length).toBe(1);
    wrapper.find('ShipmentListItem').at(1).simulate('click');
    expect(onShipmentClick.mock.calls.length).toBe(2);
    wrapper.find('ShipmentListItem').at(2).simulate('click');
    expect(onShipmentClick.mock.calls.length).toBe(3);
    const [shipmentOneId, shipmentTwoId, shipmentThreeId] = onShipmentClick.mock.calls;
    expect(shipmentOneId[0]).toEqual(shipments[0].id);
    expect(shipmentTwoId[0]).toEqual(shipments[1].id);
    expect(shipmentThreeId[0]).toEqual(shipments[2].id);
  });
});
