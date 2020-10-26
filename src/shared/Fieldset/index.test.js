/*  react/jsx-props-no-spreading */
import React from 'react';
import { mount } from 'enzyme';

import Fieldset from '.';

const defaultProps = {
  children: null,
};

function mountFieldset(props = defaultProps) {
  return mount(<Fieldset {...props}>{props.children}</Fieldset>);
}
describe('Fieldset component', () => {
  it('renders expected component with class', () => {
    const wrapper = mountFieldset();
    expect(wrapper.find('fieldset.usa-fieldset').length).toBe(1);
  });
  it('renders a legend', () => {
    const legendText = 'Example legend';
    const wrapper = mountFieldset({ legend: legendText });

    expect(wrapper.find('legend').length).toBe(1);
    expect(wrapper.find('legend').text()).toBe(legendText);
  });
  it('renders hint text', () => {
    const legendText = 'Example legend';
    const hintText = 'Example hint';
    const wrapper = mountFieldset({ legend: legendText, hintText });
    expect(wrapper.find('Hint').length).toBe(1);
    expect(wrapper.find('div.hint').length).toBe(1);
    expect(wrapper.find('div.hint').text()).toBe(hintText);
  });
});
