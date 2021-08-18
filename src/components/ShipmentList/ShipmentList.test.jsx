/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { mount } from 'enzyme';
import { render } from '@testing-library/react';

import ShipmentList from '.';

describe('ShipmentList component', () => {
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

describe('BillableWeightCard', () => {
  it('renders maximum billable weight, total billable weight, weight requested and weight allowance', () => {
    const shipments = [
      { id: '0001', shipmentType: 'HHG', billableWeightCap: '5,600' },
      { id: '0002', shipmentType: 'HHG', billableWeightCap: '3,200', reweigh: { id: '1234' } },
      { id: '0003', shipmentType: 'HHG', billableWeightCap: '3,400' },
    ];
    const entitlements = [
      { id: '1234', shipmentId: '0001', authorizedWeight: '4,600' },
      { id: '12346', shipmentId: '0002', authorizedWeight: '4,600' },
      { id: '12347', shipmentId: '0003', authorizedWeight: '4,600' },
    ];
    const defaultProps = {
      shipments,
      entitlements,
      moveSubmitted: false,
      showShipmentWeight: true,
    };

    const { getByText } = render(<ShipmentList {...defaultProps} />);

    // flags
    expect(getByText('Over weight')).toBeInTheDocument();
    expect(getByText('Missing weight')).toBeInTheDocument();

    // weights
    expect(getByText(`${shipments[0].billableWeightCap} lbs`)).toBeInTheDocument();
    expect(getByText(`${shipments[1].billableWeightCap} lbs`)).toBeInTheDocument();
    expect(getByText(`${shipments[2].billableWeightCap} lbs`)).toBeInTheDocument();
  });
});
