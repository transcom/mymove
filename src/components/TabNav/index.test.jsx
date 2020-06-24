import React from 'react';
import { shallow } from 'enzyme';
import { Tag } from '@trussworks/react-uswds';

import TabNav from '.';

describe('TabNav', () => {
  it('should render the tab navigation', () => {
    const options = [
      <a href="/test">
        <span className="tab-title">Link 1</span>
        <Tag>5</Tag>
      </a>,
      <a href="/test2">
        <span className="tab-title">Link 2</span>
      </a>,
    ];
    const wrapper = shallow(<TabNav items={options} />);
    expect(wrapper.find('.tab-title').first().text()).toBe('Link 1');
    expect(wrapper.find(Tag).length).toBe(1);
    expect(wrapper.find(Tag).children().text()).toBe('5');
    expect(wrapper.find('.tab-title').last().text()).toBe('Link 2');
  });
});
