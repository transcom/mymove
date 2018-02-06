import React from 'react';
import { mount } from 'enzyme';
import { AwardedShipments } from './AwardedShipments';

const loadAwardedShipments = () => {};

describe('No Awarded Shipments and Errors', () => {
  let wrapper;

  beforeEach(() => {
    const shipments = null;
    const hasError = true;
    wrapper = mount(
      <AwardedShipments
        hasError={hasError}
        shipments={shipments}
        loadAwardedShipments={loadAwardedShipments}
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
      {
        id: '10',
        name: 'Sally Shipment',
        traffic_distribution_list: 'Piggy Porters',
      },
      {
        id: '20',
        name: 'Nino Shipment',
        traffic_distribution_list: 'Shleppers',
      },
    ];
    const hasError = false;
    wrapper = mount(
      <AwardedShipments
        hasError={hasError}
        shipments={shipments}
        loadAwardedShipments={loadAwardedShipments}
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
