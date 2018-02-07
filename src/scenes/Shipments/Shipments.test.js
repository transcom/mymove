import React from 'react';
import { mount } from 'enzyme';
import { Shipments } from './Shipments';

const loadShipments = () => {};

describe('No available shipments or errors', () => {
  let wrapper;

  beforeEach(() => {
    const shipments = null;
    const hasError = true;
    const shipmentStatus = 'available';
    wrapper = mount(
      <Shipments
        hasError={hasError}
        shipments={shipments}
        loadShipments={loadShipments}
        params={shipmentStatus}
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

describe('Has available shipments', () => {
  let wrapper;
  const pickup_date = new Date(2018, 11, 17);
  const delivery_date = new Date(2018, 11, 19);

  beforeEach(() => {
    const shipments = [
      {
        id: '10',
        name: 'Sally Shipment',
        pickup_date: pickup_date.toString(),
        delivery_date: delivery_date.toString(),
      },
      {
        id: '20',
        name: 'Nino Shipment',
        pickup_date: pickup_date.toString(),
        delivery_date: delivery_date.toString(),
      },
    ];
    const hasError = false;
    wrapper = mount(
      <Shipments
        hasError={hasError}
        shipments={shipments}
        loadShipments={loadShipments}
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
