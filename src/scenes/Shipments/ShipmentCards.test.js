import React from 'react';
import { shallow } from 'enzyme';
import ShipmentCards from './ShipmentCards';

describe('Null Shipments on ShipmentCards', () => {
  let wrapper;
  const shipments = null;

  beforeEach(() => {
    wrapper = shallow(<ShipmentCards shipments={shipments} />);
  });

  it('renders without crashing', () => {
    expect(wrapper.find('ShipmentCards').toExist);
  });
});

describe('Empty Shipments on ShipmentCards', () => {
  let wrapper;
  const shipments = [];

  beforeEach(() => {
    wrapper = shallow(<ShipmentCards shipments={shipments} />);
  });

  it('renders without crashing', () => {
    expect(wrapper.find('ShipmentCards').toExist);
  });
});

describe('Shipments on ShipmentCards', () => {
  let wrapper;
  const shipments = [
    {
      id: '13',
      name: 'Sally Shipment',
      traffic_distribution_list: 'Piggy Packers',
    },
  ];

  beforeEach(() => {
    wrapper = shallow(<ShipmentCards shipments={shipments} />);
  });

  it('renders without crashing', () => {
    expect(wrapper.find('ShipmentCards').toExist);
  });
});
