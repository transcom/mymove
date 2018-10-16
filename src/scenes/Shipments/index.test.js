import React from 'react';
import { mount } from 'enzyme';
import { Shipments } from '.';

const loadShipments = () => {};

describe('No available shipments or errors', () => {
  let wrapper;
  const shipments = [];
  const hasError = true;
  // Match is a param on props that allows access to url parameter
  const match = { params: { shipmentsStatus: 'available' } };

  beforeEach(() => {
    wrapper = mount(
      <Shipments hasError={hasError} shipments={shipments} loadShipments={loadShipments} match={match} />,
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
  const match = { params: { shipmentsStatus: 'available' } };

  beforeEach(() => {
    const shipments = [
      {
        id: '10',
        name: 'Sally Shipment',
        pickup_date: pickup_date.toString(),
        delivery_date: delivery_date.toString(),
        traffic_distribution_list_id: '123',
        transportation_service_provider_id: null,
      },
      {
        id: '20',
        name: 'Nino Shipment',
        pickup_date: pickup_date.toString(),
        delivery_date: delivery_date.toString(),
        traffic_distribution_list_id: '123',
        transportation_service_provider_id: null,
      },
    ];
    const hasError = false;
    wrapper = mount(
      <Shipments hasError={hasError} shipments={shipments} loadShipments={loadShipments} match={match} />,
    );
  });

  it('renders without an alert', () => {
    expect(wrapper.find('Alert').length).toBe(0);
  });

  it('renders issue cards without crashing', () => {
    expect(wrapper.find('ShipmentCards').length).toBe(1);
  });
});

describe('No awarded shipments or errors', () => {
  let wrapper;
  const shipments = [];
  const hasError = true;
  // Match is a param on props that allows access to url parameter
  const match = { params: { shipmentsStatus: 'awarded' } };

  beforeEach(() => {
    wrapper = mount(
      <Shipments hasError={hasError} shipments={shipments} loadShipments={loadShipments} match={match} />,
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
  const match = { params: { shipmentsStatus: 'awarded' } };

  beforeEach(() => {
    const shipments = [
      {
        id: '10',
        name: 'Sally Shipment',
        pickup_date: pickup_date.toString(),
        delivery_date: delivery_date.toString(),
        traffic_distribution_list_id: '123',
        shipment_id: '13',
        transportation_service_provider_id: '20',
        administrative_shipment: false,
      },
      {
        id: '20',
        name: 'Nino Shipment',
        pickup_date: pickup_date.toString(),
        delivery_date: delivery_date.toString(),
        traffic_distribution_list_id: '123',
        shipment_id: '13',
        transportation_service_provider_id: '20',
        administrative_shipment: false,
      },
    ];
    const hasError = false;
    wrapper = mount(
      <Shipments hasError={hasError} shipments={shipments} loadShipments={loadShipments} match={match} />,
    );
  });

  it('renders without an alert', () => {
    expect(wrapper.find('Alert').length).toBe(0);
  });

  it('renders issue cards without crashing', () => {
    expect(wrapper.find('ShipmentCards').length).toBe(1);
  });
});
