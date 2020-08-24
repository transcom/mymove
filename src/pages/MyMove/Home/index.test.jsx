/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { mount } from 'enzyme';

import Home from '.';

const defaultProps = {};

function mountHome(props = defaultProps) {
  return mount(<Home {...props} />);
}
describe('Home component', () => {
  it('renders Home with the right amount of components', () => {
    const wrapper = mountHome();
    expect(wrapper.find('Step').length).toBe(4);
    expect(wrapper.find('.usa-alert--success').length).toBe(1);
    expect(wrapper.find('Helper').length).toBe(1);
    expect(wrapper.find('Contact').length).toBe(1);
  });
});
