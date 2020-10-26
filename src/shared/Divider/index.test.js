/*  react/jsx-props-no-spreading */
import React from 'react';
import { mount } from 'enzyme';

import Divider from '.';

function mountDivider(props) {
  return mount(<Divider {...props} />);
}
describe('Divider component', () => {
  it('renders expected component', () => {
    const wrapper = mountDivider();
    expect(wrapper.find('hr').length).toBe(1);
  });
});
