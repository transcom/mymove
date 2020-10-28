/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { shallow } from 'enzyme';

import Hint from '.';

describe('Hint component', () => {
  it('renders expected component with class', () => {
    const wrapper = shallow(<Hint>Test Hint</Hint>);
    expect(wrapper.find('div.Hint').length).toBe(1);
    expect(wrapper.text()).toBe('Test Hint');
  });
});
