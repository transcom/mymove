import React from 'react';
import { shallow } from 'enzyme';
import { Tag } from '@trussworks/react-uswds';
import TabNav from '.';

describe('TabNav', () => {
  it('should render the tab navigation', () => {
    const options = [
      {
        title: 'Option 1',
        active: false,
        notice: null,
      },
      {
        title: 'Option 2',
        active: true,
        notice: '2',
      },
      {
        title: 'Option 3',
        active: false,
        notice: null,
      },
    ];
    const wrapper = shallow(<TabNav options={options} />);
    expect(
      wrapper
        .find('.tab-title')
        .first()
        .text(),
    ).toBe('Option 1');
    expect(wrapper.find('.tab-active').text()).toBe('Option 2');
    expect(wrapper.find(Tag).length).toBe(1);
    expect(
      wrapper
        .find(Tag)
        .children()
        .text(),
    ).toBe(1);
    expect(
      wrapper
        .find('.tab-title')
        .last()
        .text(),
    ).toBe('Option 3');
  });
});
