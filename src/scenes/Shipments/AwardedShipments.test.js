import React from 'react';
import { mount } from 'enzyme';
import { AwardedShipments } from './AwardedShipments';

const loadAwardedShipments = () => {};

describe('No awarded shipments or errors', () => {
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

describe('Has awarded shipments', () => {
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
        traffic_distribution_list_id: 'ab1eace7-ec68-4794-883d-bc69b16f0fe',
      },
      {
        id: '20',
        name: 'Nino Shipment',
        pickup_date: pickup_date.toString(),
        delivery_date: delivery_date.toString(),
        traffic_distribution_list_id: 'ab1eace7-ec68-4794-883d-bc6db16f20fe',
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
