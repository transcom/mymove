/*  react/jsx-props-no-spreading */
import React from 'react';
import { mount } from 'enzyme';

import Hint from '.';

function mountHint(props) {
  return mount(<Hint {...props} />);
}
describe('Hint component', () => {
  it('renders expected component with class', () => {
    const wrapper = mountHint();
    expect(wrapper.find('div.hint').length).toBe(1);
  });
});
