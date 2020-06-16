import React from 'react';
import { shallow } from 'enzyme';
import { Tag } from '@trussworks/react-uswds';
import TabNav from '.';

describe('TabNav', () => {
  it('should render the tab navigation', () => {
    const options = [
      <a href="/test">
        <span>Link 1</span>
      </a>,
      <a href="/test2">
        <span>Link 2</span>
      </a>,
    ];
    const wrapper = shallow(<TabNav items={options} />);
    expect(wrapper.find('.tab-title').first().text()).toBe('Option 1');
    expect(wrapper.find(Tag).length).toBe(1);
    expect(wrapper.find(Tag).children().text()).toBe('2');
    expect(wrapper.find('.tab-title').last().text()).toBe('Option 3');
  });
});
