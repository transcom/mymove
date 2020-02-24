import React from 'react';
import { shallow } from 'enzyme';
import { TabPanel } from 'react-tabs';
import { Tag } from '@trussworks/react-uswds';
import TabNav from '.';

describe('TabNav', () => {
  it('should render the tab navigation', () => {
    const options = [
      {
        title: 'Option 1',
        notice: null,
      },
      {
        title: 'Option 2',
        notice: '2',
      },
      {
        title: 'Option 3',
        notice: null,
      },
    ];
    const wrapper = shallow(
      <TabNav options={options}>
        <TabPanel>Body Of Tab 1</TabPanel>
        <TabPanel>Body Of Tab 2</TabPanel>
        <TabPanel>Body Of Tab 3</TabPanel>
      </TabNav>,
    );
    expect(
      wrapper
        .find('.tab-title')
        .first()
        .text(),
    ).toBe('Option 1');
    expect(wrapper.find(Tag).length).toBe(1);
    expect(
      wrapper
        .find(Tag)
        .children()
        .text(),
    ).toBe('2');
    expect(
      wrapper
        .find('.tab-title')
        .last()
        .text(),
    ).toBe('Option 3');
  });
});
