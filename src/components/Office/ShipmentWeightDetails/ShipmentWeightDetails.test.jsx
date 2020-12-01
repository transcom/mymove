import React from 'react';
import { mount } from 'enzyme';

import ShipmentWeightDetails from './ShipmentWeightDetails';

describe('ShipmentWeightDetails', () => {
  it('renders without crashing', () => {
    const wrapper = mount(<ShipmentWeightDetails />);
    expect(wrapper.find('DataPointGroup')).toHaveLength(1);
  });

  it('renders with estimated weight', () => {
    const wrapper = mount(<ShipmentWeightDetails estimatedWeight={1111} />);
    const text = wrapper.find('DataPointGroup').text();
    expect(text).toContain('Estimated weight');
    expect(text).toContain('1,111 lbs');
  });

  it('renders with actual weight', () => {
    const wrapper = mount(<ShipmentWeightDetails actualWeight={1111} />);
    const text = wrapper.find('DataPointGroup').text();
    expect(text).toContain('Actual weight');
    expect(text).toContain('1,111 lbs');
  });
});
