/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { mount } from 'enzyme';
import { render, screen } from '@testing-library/react';

import ShipmentList from '.';

import { SHIPMENT_OPTIONS } from 'shared/constants';
import { formatWeight } from 'utils/formatters';

describe('ShipmentList component', () => {
  const shipments = [
    { id: 'ID-1', shipmentType: SHIPMENT_OPTIONS.PPM },
    { id: 'ID-2', shipmentType: SHIPMENT_OPTIONS.HHG },
    { id: 'ID-3', shipmentType: SHIPMENT_OPTIONS.NTS },
    { id: 'ID-4', shipmentType: SHIPMENT_OPTIONS.NTSR },
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
    expect(wrapper.find('.shipment-list-item-NTS-release').length).toBe(1);
    expect(wrapper.find('.shipment-list-item-NTS-release strong').text()).toBe('NTS-release');
    expect(wrapper.find('.shipment-list-item-NTS-release span').text()).toBe('#ID-4');
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

describe('Shipment List being used for billable weight', () => {
  it('renders maximum billable weight, total billable weight, weight requested and weight allowance with no flags', () => {
    const shipments = [
      {
        id: '0001',
        shipmentType: SHIPMENT_OPTIONS.HHG,
        calculatedBillableWeight: 1161,
        primeEstimatedWeight: 200,
        reweigh: { id: '1234', weight: 50 },
      },
      {
        id: '0002',
        shipmentType: SHIPMENT_OPTIONS.HHG,
        calculatedBillableWeight: 3200,
        primeEstimatedWeight: 3000,
        reweigh: { id: '1234' },
      },
      {
        id: '0003',
        shipmentType: SHIPMENT_OPTIONS.HHG,
        calculatedBillableWeight: 3000,
        primeEstimatedWeight: 3000,
        reweigh: { id: '1234', weight: 40 },
      },
    ];

    const defaultProps = {
      shipments,
      moveSubmitted: false,
      showShipmentWeight: true,
    };

    render(<ShipmentList {...defaultProps} />);

    // flags
    expect(screen.queryByText('Over weight')).toBeInTheDocument();
    expect(screen.queryByText('Missing weight')).toBeInTheDocument();

    // weights
    expect(screen.getByText(formatWeight(shipments[0].calculatedBillableWeight))).toBeInTheDocument();
    expect(screen.getByText(formatWeight(shipments[1].calculatedBillableWeight))).toBeInTheDocument();
    expect(screen.getByText(formatWeight(shipments[2].calculatedBillableWeight))).toBeInTheDocument();
  });

  it('does not display weight flags when not appropriate', () => {
    const shipments = [
      { id: '0001', shipmentType: SHIPMENT_OPTIONS.HHG, calculatedBillableWeight: 5666, primeEstimatedWeight: 5600 },
      {
        id: '0002',
        shipmentType: SHIPMENT_OPTIONS.HHG,
        calculatedBillableWeight: 3200,
        primeEstimatedWeight: 3000,
        reweigh: { id: '1234', weight: 3400 },
      },
      { id: '0003', shipmentType: SHIPMENT_OPTIONS.HHG, calculatedBillableWeight: 5400, primeEstimatedWeight: 5000 },
      // we don't show flags for ntsr shipments - if this was an hhg, it would show a missing weight warning
      { id: '0004', shipmentType: SHIPMENT_OPTIONS.NTSR },
    ];

    const defaultProps = {
      shipments,
      moveSubmitted: false,
      showShipmentWeight: true,
    };

    render(<ShipmentList {...defaultProps} />);

    // flags
    expect(screen.queryByText('Over weight')).not.toBeInTheDocument();
    expect(screen.queryByText('Missing weight')).not.toBeInTheDocument();
  });
});
