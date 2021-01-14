import React from 'react';
import { mount } from 'enzyme';

import App from './index';

it('renders without crashing', () => {
  const wrapper = mount(<App />);
  expect(wrapper.exists()).toBe(true);
});

// todo: add tests for routing
