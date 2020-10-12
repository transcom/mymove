import React from 'react';
import { mount } from 'enzyme';

import MoveQueue from './MoveQueue';

describe('MoveQueue', () => {
  const wrapper = mount(<MoveQueue />);

  it('should render the h1', () => {
    expect(wrapper.find('h1').text()).toBe('All moves');
  });

  it('should render the table', () => {
    expect(wrapper.find('Table').exists()).toBe(true);
  });
});
