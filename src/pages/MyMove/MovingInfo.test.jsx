import React from 'react';
import { mount } from 'enzyme';

import { MovingInfo } from './MovingInfo';

const wrapper = mount(<MovingInfo />);

describe('MovingInfo component', () => {
  it('renders', () => {
    // expect(wrapper.find('MovingInfo').length).toBe(1);
    expect(wrapper.find('MovingInfo').exists()).toBe(true);
    // expect(wrapper.find(Radio).length).toBe(2);

    // check for weight estimate
    // expect(wrapper.find(Radio).at(0).text()).toContain('I’ll move things myself');
  });
});
