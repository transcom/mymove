import React from 'react';
import { mount } from 'enzyme';
import { AvailableShipments } from './AvailableShipments';

const loadAvailableShipments = () => {};

describe('No Available Shipments and Errors', () => {
  let wrapper;

  beforeEach(() => {
    const shipments = null;
    const hasError = true;
    wrapper = mount(
      <AvailableShipments
        hasError={hasError}
        shipments={shipments}
        loadAvailableShipments={loadAvailableShipments}
      />,
    );
  });

  it('renders an alert', () => {
    expect(wrapper.find('Alert').length).toBe(1);
  });

  it('does not render issue cards', () => {
    expect(wrapper.find('ShipmentCards').length).toBe(0);
  });
});

describe('Has shipments', () => {
  let wrapper;

  beforeEach(() => {
    const shipments = [
      { id: '10', name: 'Sally Shipment' },
      { id: '20', name: 'Nino Shipment' },
    ];
    const hasError = false;
    wrapper = mount(
      <AvailableShipments
        hasError={hasError}
        shipments={shipments}
        loadAvailableShipments={loadAvailableShipments}
      />,
    );
  });

  it('renders without an alert', () => {
    expect(wrapper.find('Alert').length).toBe(0);
  });

  it('renders issue cards without crashing', () => {
    expect(wrapper.find('ShipmentCards').length).toBe(1);
  });
});
