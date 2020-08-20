/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { mount } from 'enzyme';

import Helper from './Helper';

const defaultProps = {
  containerStyles: {},
  title: '',
  helpList: [],
  description: '',
};
function mountHelper(props = defaultProps) {
  return mount(<Helper {...props} />);
}
describe('Helper component', () => {
  it('renders Helper with description', () => {
    const title = 'Title';
    const description = 'description';
    const props = {
      title,
      description,
    };
    const wrapper = mountHelper(props);

    expect(wrapper.find('h3').text()).toBe(title);
    expect(wrapper.find('p').text()).toBe(description);
  });
  it('renders Helper with helpList', () => {
    const title = 'Title';
    const helpList = ['bullet 1', 'bullet 2', 'bullet 3'];
    const props = {
      title,
      helpList,
    };
    const wrapper = mountHelper(props);
    expect(wrapper.find('h3').text()).toBe(title);
    // first item does not have padding
    expect(wrapper.find('li').at(0).props().className).toBe('');
    expect(wrapper.find('li').at(0).text()).toBe(helpList[0]);
    expect(wrapper.find('li').at(1).text()).toBe(helpList[1]);
    expect(wrapper.find('li').at(2).text()).toBe(helpList[2]);
  });
});
